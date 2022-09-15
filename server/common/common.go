package common

import (
	"errors"
	"time"

	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var ValidationError = errors.New("INVALID_INPUT")
var ErrorInvalidType = errors.New("INVALID_TYPE")

func ValidateStruct(s interface{}, validate *validator.Validate) []string {

	errorFields := []string{}
	err := validate.Struct(s)
	if err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			errorMessage := err.StructNamespace() + " (tag:" + err.Tag() + ", should_be:" + err.Param() + ")"
			errorFields = append(errorFields, errorMessage)
		}
	}

	if len(errorFields) == 0 {
		return nil
	}
	return errorFields

}

type Author struct {
	ID        primitive.ObjectID `bson:"_id" json:"-"`
	FirstName string             `bson:"first_name" json:"firstName" validate:"required"`
	LastName  string             `bson:"last_name" json:"lastName" validate:"required"`
}

type Source struct {
	ID      primitive.ObjectID `bson:"_id" json:"-"`
	Authors []*Author          `bson:"authors" json:"authors" validate:"required,min=1,dive"`
}
type DefinitionFilter struct {
	Title           *string      `json:"title" bson:"title" validate:"omitempty"`
	Content         *string      `json:"content" bson:"content" validate:"omitempty"`
	PublishingDates *[]time.Time `json:"publishing_dates" bson:"publishing_dates" validate:"omitempty,min=1"`
	Authors         *[]*Author   `json:"authors" bson:"authors" validate:"omitempty,min=1,dive"`
	Sources         *[]*Source   `json:"sources" bson:"sources" validate:"omitempty,min=1,dive"`
	Tags            *[]string    `json:"tags" bson:"tags" validate:"omitempty,min=1"`
}
