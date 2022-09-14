package main

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"

	"yacoid_server/common"
	"yacoid_server/database"

	validator "github.com/go-playground/validator/v10"
	fiber "github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
)

const DATABASE_ADDRESS = "localhost"
const DATABASE_PORT = 27017
const REST_PORT = 3000

var ErrorEmailVerification = errors.New("EMAIL_VERIFICATION_ERROR")
var ErrorChangePassword = errors.New("CHANGE_PASSWORD_ERROR")

var ErrorCodeMap map[error]int = map[error]int{}

func GetErrorCode(err error) int {

	code, exists := ErrorCodeMap[err]

	if exists {
		return code
	}

	return fiber.StatusInternalServerError

}

// TODO: Pagination, Filter

func main() {

	ErrorCodeMap[database.InvalidID] = fiber.StatusBadRequest
	ErrorCodeMap[common.ValidationError] = fiber.StatusBadRequest

	ErrorCodeMap[database.ErrorUserNotFound] = fiber.StatusNotFound
	ErrorCodeMap[database.ErrorNotEnoughPermissions] = fiber.StatusUnauthorized
	ErrorCodeMap[database.ErrorInvalidCredentials] = fiber.StatusUnauthorized
	ErrorCodeMap[database.ErrorPasswordResetExpiryDateExceeded] = fiber.StatusBadRequest
	ErrorCodeMap[database.ErrorUserAlreadyExists] = fiber.StatusBadRequest
	ErrorCodeMap[database.ErrorUserAlreadyLoggedIn] = fiber.StatusBadRequest

	ErrorCodeMap[database.ErrorDefinitionNotFound] = fiber.StatusNotFound
	ErrorCodeMap[database.ErrorDefinitionAlreadyApproved] = fiber.StatusBadRequest
	ErrorCodeMap[database.ErrorDefinitionRejectionBelongsToAnotherUser] = fiber.StatusUnauthorized
	ErrorCodeMap[database.ErrorDefinitionRejectionNotAnsweredYet] = fiber.StatusBadRequest

	fmt.Println("Starting server...")

	dbContext, client := database.Connect(DATABASE_ADDRESS, DATABASE_PORT)
	defer client.Disconnect(dbContext)

	app := fiber.New()

	api := app.Group("/api")

	validate := validator.New()
	validate.RegisterValidation("ISO8601date", IsISO8601Date)

	definitionApi := api.Group("/definitions")
	AddDefinitionRequests(&definitionApi, validate)

	authApi := api.Group("/auth")
	AddAuthRequests(&authApi, validate)

	userApi := api.Group("/user")
	AddAuthRequests(&userApi, validate)

	fmt.Println("Started server on port " + strconv.Itoa(REST_PORT))

	app.Listen(":" + strconv.Itoa(REST_PORT))

}

func AddAuthRequests(authApi *fiber.Router, validate *validator.Validate) {

	(*authApi).Post("/register", func(ctx *fiber.Ctx) error {

		input := new(database.User)

		if err := ctx.BodyParser(input); err != nil {
			return ctx.Status(500).JSON(bson.M{"error": err})
		}

		validateErrors := input.Validate(validate)

		if validateErrors != nil {
			return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error":      "Error on fields: " + strings.Join(validateErrors, ", "),
				"definition": nil,
			})
		}

		user, err := database.Register(*input)
		if err != nil {
			return ctx.Status(500).JSON(bson.M{"error": err.Error()})
		}
		return ctx.JSON(bson.M{"error": nil, "user": user})

	})

	(*authApi).Get("/login/:email/:password", func(ctx *fiber.Ctx) error {

		email := ctx.Params("email")
		password := ctx.Params("password")
		fmt.Println("Input", email, password)
		user, err := database.Login(email, password)
		fmt.Println(user, err)
		if err != nil {
			return ctx.Status(500).JSON(bson.M{"error": err.Error()})
		}

		/* Hide some attributes */
		user.PasswordSalt = ""
		user.PasswordHash = ""

		return ctx.JSON(bson.M{"error": nil, "user": user})

	})

	(*authApi).Get("/password_salt/:email", func(ctx *fiber.Ctx) error {

		email := ctx.Params("email")
		fmt.Println("Input", email)
		salt, err := database.GetPasswordSalt(email)
		fmt.Println("SALT", salt, err)
		if err != nil {
			return ctx.Status(500).JSON(bson.M{"error": err.Error()})
		}
		return ctx.JSON(bson.M{"error": nil, "passwordSalt": salt})

	})

	(*authApi).Get("/logout", func(ctx *fiber.Ctx) error {

		err := database.Logout(ctx.GetReqHeaders()["Authtoken"])
		if err != nil {
			return ctx.Status(500).JSON(bson.M{"error": err.Error()})
		}
		return ctx.JSON(bson.M{"error": nil})

	})

	(*authApi).Get("/request_password_reset/:email", func(ctx *fiber.Ctx) error {

		email := ctx.Params("email")
		token, err := database.InitiatePasswordReset(email)
		if err != nil {
			return ctx.Status(500).JSON(bson.M{"error": err.Error()})
		}
		return ctx.JSON(bson.M{"error": nil, "token": token})

	})

	(*authApi).Get("/reset_password/:token/:password_hash", func(ctx *fiber.Ctx) error {

		token := ctx.Params("token")
		passwordHash := ctx.Params("password_hash")
		err := database.ResetPassword(token, passwordHash)
		if err != nil {
			return ctx.Status(500).JSON(bson.M{"error": err.Error()})
		}
		return ctx.JSON(bson.M{"error": nil})

	})

}

type DeleteUserRequest struct {
	PasswordHash string `bson:"password_hash,omitempty" json:"passwordHash,omitempty"`
	Reason       string `bson:"reason,omitempty" json:"reason,omitempty"`
}

type ChangeAccountDataRequest struct {
	FirstName       *string `bson:"first_name,omitempty" json:"firstName,omitempty"`
	LastName        *string `bson:"last_name,omitempty" json:"lastName,omitempty"`
	Email           *string `bson:"email,omitempty" json:"email,omitempty"`
	City            *string `bson:"city,omitempty" json:"city,omitempty"`
	CurrentPassword *string `bson:"current_password,omitempty" json:"currentPassword,omitempty"`
	NewPassword     *string `bson:"new_password,omitempty" json:"newPassword,omitempty"`
}

func AddUserRequests(userApi *fiber.Router, validate *validator.Validate) {

	(*userApi).Post("/delete_user", func(ctx *fiber.Ctx) error {

		request := new(DeleteUserRequest)

		if err := ctx.BodyParser(request); err != nil {
			return ctx.Status(500).JSON(bson.M{"error": err})
		}

		authToken := ctx.GetReqHeaders()["Authtoken"]
		err := database.DeleteUser(authToken, request.PasswordHash, request.Reason)
		if err != nil {
			return ctx.Status(500).JSON(bson.M{"error": err.Error()})
		}
		return ctx.JSON(bson.M{"error": nil})

	})

	(*userApi).Post("/change_account_data", func(ctx *fiber.Ctx) error {

		request := new(ChangeAccountDataRequest)

		fmt.Println("1")
		if err := ctx.BodyParser(request); err != nil {
			return ctx.Status(500).JSON(bson.M{"error": err.Error()})
		}
		fmt.Println(request)

		authToken := ctx.GetReqHeaders()["Authtoken"]
		response, err := database.ChangeAccountData(authToken, request.FirstName, request.LastName, request.Email, request.City, request.CurrentPassword, request.NewPassword)

		if response.EmailVerification != nil && response.EmailVerification.Error != nil {
			errorText := ErrorEmailVerification.Error()
			response.EmailVerification.Error = &errorText
		}

		if response.ChangePassword != nil && response.ChangePassword.Error != nil {
			if *response.ChangePassword.Error != "INVALID_CREDENTIALS" {
				errorText := ErrorChangePassword.Error()
				response.ChangePassword.Error = &errorText
			}
		}

		fmt.Println(response, err)

		if err != nil {
			return ctx.Status(500).JSON(bson.M{"error": err.Error()})
		}

		return ctx.JSON(bson.M{"error": nil, "response": response})

	})
}

func AddDefinitionRequests(definitionApi *fiber.Router, validate *validator.Validate) {

	(*definitionApi).Get("/definition/:id", func(ctx *fiber.Ctx) error {

		id := ctx.Params("id")

		definition, err := database.GetDefinitionById(id)

		if err != nil {
			return ctx.Status(GetErrorCode(err)).JSON(fiber.Map{
				"error":      err.Error(),
				"definition": nil,
			})
		}

		return ctx.JSON(fiber.Map{
			"error":      nil,
			"definition": definition,
		})

	})

	/*(*definitionApi).Post("/definition", func(ctx *fiber.Ctx) error {

		input := new(database.Definition)

		if err := ctx.BodyParser(input); err != nil {
			return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": err,
				"definition": nil,
			})
		}

		definition, err := database.CreateDefinition(input)
		if err != nil {
			return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": err.Error(),
				"definition": nil,
			})
		}

		return ctx.JSON(fiber.Map{
			"error": nil,
			"definition": definition,
		})
	})*/

	(*definitionApi).Post("/submit", func(ctx *fiber.Ctx) error {

		input := new(database.Definition)

		if err := ctx.BodyParser(input); err != nil {
			return ctx.Status(GetErrorCode(err)).JSON(fiber.Map{
				"error":      err.Error(),
				"definition": nil,
			})
		}

		validateErrors := input.Validate(validate)

		if validateErrors != nil {
			return ctx.Status(GetErrorCode(common.ValidationError)).JSON(fiber.Map{
				"error":      "Error on fields: " + strings.Join(validateErrors, ", "),
				"definition": nil,
			})
		}

		authToken := ctx.GetReqHeaders()["Authtoken"]
		definition, err := database.SubmitDefinition(input, authToken)
		if err != nil {
			return ctx.Status(GetErrorCode(err)).JSON(fiber.Map{
				"error":      err.Error(),
				"definition": nil,
			})
		}

		return ctx.JSON(fiber.Map{
			"error":      nil,
			"definition": definition,
		})
	})

	(*definitionApi).Get("/approve/:id", func(ctx *fiber.Ctx) error {

		definitionId := ctx.Params("id")

		authToken := ctx.GetReqHeaders()["Authtoken"]
		err := database.ApproveDefinition(definitionId, authToken)

		if err != nil {
			return ctx.Status(GetErrorCode(err)).JSON(fiber.Map{
				"error": err.Error(),
			})
		}

		return ctx.JSON(fiber.Map{
			"error": nil,
		})

	})

	(*definitionApi).Post("/reject", func(ctx *fiber.Ctx) error {

		input := new(RejectRequest)

		if err := ctx.BodyParser(input); err != nil {
			return ctx.Status(GetErrorCode(err)).JSON(fiber.Map{
				"error":      err.Error(),
				"definition": nil,
			})
		}

		validateErrors := input.Validate(validate)

		if validateErrors != nil {
			return ctx.Status(GetErrorCode(common.ValidationError)).JSON(fiber.Map{
				"error":      "Error on fields: " + strings.Join(validateErrors, ", "),
				"definition": nil,
			})
		}

		authToken := ctx.GetReqHeaders()["Authtoken"]
		err := database.RejectDefinition(input.ID, authToken, input.Content)
		if err != nil {
			return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": err.Error(),
			})
		}

		return ctx.JSON(fiber.Map{
			"error": nil,
		})
	})

	(*definitionApi).Post("/change", func(ctx *fiber.Ctx) error {

		input := new(ChangeDefinitionRequest)

		if err := ctx.BodyParser(input); err != nil {
			return ctx.Status(GetErrorCode(err)).JSON(fiber.Map{
				"error":      err.Error(),
				"definition": nil,
			})
		}

		validateErrors := input.Validate(validate)

		if validateErrors != nil {
			return ctx.Status(GetErrorCode(common.ValidationError)).JSON(fiber.Map{
				"error":      "Error on fields: " + strings.Join(validateErrors, ", "),
				"definition": nil,
			})
		}

		authToken := ctx.GetReqHeaders()["Authtoken"]
		err := database.ChangeDefinition(input.ID, input.Title, input.Content, input.Source, input.Tags, authToken)
		if err != nil {
			return ctx.Status(GetErrorCode(err)).JSON(fiber.Map{
				"error": err.Error(),
			})
		}

		return ctx.JSON(fiber.Map{
			"error": nil,
		})
	})

	(*definitionApi).Get("/newest_definitions/:limit?", func(ctx *fiber.Ctx) error {

		limitParam := ctx.Params("limit")
		var limit int

		if len(limitParam) == 0 {
			limit = 4
		} else {
			tempLimit, err := strconv.Atoi(limitParam)

			if err != nil {
				return ctx.Status(GetErrorCode(err)).JSON(fiber.Map{
					"error":       "INVALID_LIMIT",
					"definitions": nil,
				})
			}

			limit = tempLimit

		}

		definitions, err := database.GetNewestDefinitions(limit)

		if err != nil {
			return ctx.Status(GetErrorCode(err)).JSON(fiber.Map{
				"error":       "DEFINITIONS_COULD_NOT_BE_RETRIEVED",
				"definitions": nil,
			})
		}

		return ctx.JSON(fiber.Map{
			"error":       nil,
			"definitions": definitions,
		})

	})
}

type RejectRequest struct {
	ID      string `json:"id" validate:"required"`
	Content string `json:"content" validate:"required,min=1"`
}

func (rejection *RejectRequest) Validate(validate *validator.Validate) []string {
	return common.ValidateStruct(rejection, validate)
}

type ChangeDefinitionRequest struct {
	ID             string           `json:"id" validate:"required"`
	Title          *string          `json:"title"`
	Content        *string          `json:"content"`
	Source         *database.Source `json:"source" validate:"omitempty,dive"`
	PublishingDate *time.Time       `json:"publishingDate" validate:"omitempty,ISO8601date"`
	Tags           *[]string        `json:"tags"`
}

func (rejection *ChangeDefinitionRequest) Validate(validate *validator.Validate) []string {
	return common.ValidateStruct(rejection, validate)
}

func IsISO8601Date(field validator.FieldLevel) bool {
	timeValue := field.Field().Interface().(time.Time)
	timeString := timeValue.Format(time.RFC3339)
	ISO8601DateRegexString := "^(?:[1-9]\\d{3}-(?:(?:0[1-9]|1[0-2])-(?:0[1-9]|1\\d|2[0-8])|(?:0[13-9]|1[0-2])-(?:29|30)|(?:0[13578]|1[02])-31)|(?:[1-9]\\d(?:0[48]|[2468][048]|[13579][26])|(?:[2468][048]|[13579][26])00)-02-29)T(?:[01]\\d|2[0-3]):[0-5]\\d:[0-5]\\d(?:\\.\\d{1,9})?(?:Z|[+-][01]\\d:[0-5]\\d)$"
	ISO8601DateRegex := regexp.MustCompile(ISO8601DateRegexString)
	return ISO8601DateRegex.MatchString(timeString)
}
