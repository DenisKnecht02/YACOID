package database

import (
	"errors"
	"fmt"
	"time"

	"yacoid_server/common"

	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var ErrorDefinitionAlreadyApproved = errors.New("DEFINITION_ALREADY_APPROVED")
var ErrorDefinitionRejectionNotAnsweredYet = errors.New("DEFINITION_REJECTION_NOT_ANSWERED_YET")
var ErrorDefinitionRejectionBelongsToAnotherUser = errors.New("DEFINITION_REJECTION_BELONGS_TO_ANOTHER_USER")

type Author struct {
	ID        primitive.ObjectID `bson:"_id" json:"-"`
	FirstName string             `bson:"first_name" json:"firstName" validate:"required"`
	LastName  string             `bson:"last_name" json:"lastName" validate:"required"`
}

type Source struct {
	ID      primitive.ObjectID `bson:"_id" json:"-"`
	Authors []*Author          `bson:"authors" json:"authors" validate:"required,min=1,dive"`
}

type Rejection struct {
	ID           primitive.ObjectID `bson:"_id" json:"-"`
	RejectedBy   primitive.ObjectID `bson:"rejected_by" json:"rejectedBy" validate:"required"`
	RejectedDate time.Time          `bson:"rejected_date" json:"rejectedDate" validate:"required"`
	Content      string             `bson:"content" json:"content" validate:"required"`
}

type Definition struct {
	ID                   primitive.ObjectID  `bson:"_id" json:"id"`
	SubmittedBy          primitive.ObjectID  `bson:"submitted_by" json:"submittedBy"`
	SubmittedDate        time.Time           `bson:"submitted_date" json:"submittedDate"`
	LastSubmitChangeDate time.Time           `bson:"last_submit_change_date" json:"lastSubmitChangeDate"`
	ApprovedBy           *primitive.ObjectID `bson:"approved_by" json:"approvedBy"`
	ApprovedDate         *time.Time          `bson:"approved_date" json:"approvedDate"`
	Approved             bool                `bson:"approved" json:"approved"`
	RejectionLog         *[]*Rejection       `bson:"rejection_log" json:"-"`
	Title                string              `bson:"title" json:"title" validate:"required"`
	Content              string              `bson:"content" json:"content" validate:"required"`
	Source               *Source             `bson:"source" json:"source" validate:"required,dive"`
	PublishingDate       time.Time           `bson:"publishing_date" json:"publishingDate" validate:"ISO8601date"`
	Tags                 *[]string           `bson:"tags" json:"tags"`
}

func (definition *Definition) Validate(validate *validator.Validate) []string {
	return common.ValidateStruct(definition, validate)
}

func (definition *Definition) IsApproved() bool {
	return definition.ApprovedBy != nil && definition.ApprovedDate != nil
}

func SubmitDefinition(definition *Definition, authToken string) (*Definition, error) {

	user, userError := GetUserByAuthToken(authToken)

	if userError != nil {
		return nil, userError
	}

	now := time.Now()
	definition.ID = primitive.NewObjectID()
	definition.SubmittedBy = user.ID
	definition.SubmittedDate = now
	definition.LastSubmitChangeDate = now
	definition.ApprovedBy = nil
	definition.ApprovedDate = nil
	definition.Approved = false

	rejectionLog := []*Rejection{}
	definition.RejectionLog = &rejectionLog

	if definition.Tags == nil {
		definition.Tags = &[]string{}
	}

	_, err := definitionsCollection.InsertOne(dbContext, definition)
	// TODO: send email to user??

	if err != nil {
		return nil, err
	}

	return definition, nil

}

func ApproveDefinition(definitionId string, authToken string) error {

	definitionObjectId, definitionObjectIdError := primitive.ObjectIDFromHex(definitionId)

	if definitionObjectIdError != nil {
		return InvalidID
	}

	user, userError := GetUserByAuthToken(authToken)

	if userError != nil {
		return userError
	}

	if user.Admin == false {
		return ErrorNotEnoughPermissions
	}

	filter := bson.M{"_id": definitionObjectId}
	update := bson.M{
		"$set": bson.M{
			"approved_by":   user.ID,
			"approved_date": time.Now(),
			"approved":      true,
		},
	}

	var result bson.M
	updateError := definitionsCollection.FindOneAndUpdate(dbContext, filter, update, nil).Decode(&result)
	// TODO: send email to user

	if updateError != nil {
		if updateError == mongo.ErrNoDocuments {
			return ErrorDefinitionNotFound
		}
		return updateError
	}

	return nil

}

func RejectDefinition(definitionId string, authToken string, content string) error {

	definitionObjectId, definitionObjectIdError := primitive.ObjectIDFromHex(definitionId)

	if definitionObjectIdError != nil {
		return InvalidID
	}

	user, userError := GetUserByAuthToken(authToken)

	if userError != nil {
		return userError
	}

	if user.Admin == false {
		return ErrorNotEnoughPermissions
	}

	definition, findError := GetDefinitionByObjectId(definitionObjectId)

	if findError != nil {
		return ErrorDefinitionNotFound
	}

	if definition.Approved == true {
		return ErrorDefinitionAlreadyApproved
	}

	rejection := Rejection{
		ID:           primitive.NewObjectID(),
		RejectedBy:   user.ID,
		RejectedDate: time.Now(),
		Content:      content,
	}

	var latestRejectionDate time.Time
	for _, d := range *definition.RejectionLog {
		if d.RejectedDate.After(latestRejectionDate) {
			latestRejectionDate = d.RejectedDate
		}
	}

	if !latestRejectionDate.IsZero() && latestRejectionDate.After(definition.LastSubmitChangeDate) {
		return ErrorDefinitionRejectionNotAnsweredYet
	}

	filter := bson.M{"_id": definitionObjectId}

	update := bson.M{
		"$push": bson.M{
			"rejection_log": rejection,
		},
	}

	result := definitionsCollection.FindOneAndUpdate(dbContext, filter, update, nil)
	// TODO: send email to user

	if result.Err() != nil {
		if result.Err() == mongo.ErrNoDocuments {
			return ErrorDefinitionNotFound
		}
		return result.Err()
	}

	return nil

}

func ChangeDefinition(id string, title *string, content *string, source *Source, tags *[]string, authToken string) error {

	definitionObjectId, definitionObjectIdError := primitive.ObjectIDFromHex(id)

	if definitionObjectIdError != nil {
		return InvalidID
	}

	user, userError := GetUserByAuthToken(authToken)

	if userError != nil {
		return userError
	}

	definition, findError := GetDefinitionByObjectId(definitionObjectId)

	if findError != nil {
		return ErrorDefinitionNotFound
	}

	if definition.Approved == true {
		return ErrorDefinitionAlreadyApproved
	}

	if definition.SubmittedBy != user.ID {
		return ErrorDefinitionRejectionBelongsToAnotherUser
	}

	fmt.Println(definitionObjectId)
	filter := bson.M{"_id": definitionObjectId}

	var updateEntries bson.D
	if title != nil {
		updateEntries = append(updateEntries, bson.E{"title", title})
	}
	if content != nil {
		updateEntries = append(updateEntries, bson.E{"content", content})
	}
	if source != nil {
		updateEntries = append(updateEntries, bson.E{"source", source})
	}
	if tags != nil {
		updateEntries = append(updateEntries, bson.E{"tags", tags})
	}

	if len(updateEntries) > 0 {

		updateEntries = append(updateEntries, bson.E{"last_submit_change_date", time.Now()})
		update := bson.M{"$set": updateEntries}
		fmt.Println("UPDATE", update)

		result := definitionsCollection.FindOneAndUpdate(dbContext, filter, update, nil)

		if result.Err() != nil {
			if result.Err() == mongo.ErrNoDocuments {
				return ErrorDefinitionNotFound
			}
			return result.Err()
		}
	}

	return nil

}

func GetDefinitionById(id string) (*Definition, error) {

	objectId, idError := primitive.ObjectIDFromHex(id)

	if idError != nil {
		return nil, InvalidID
	}

	filter := bson.M{"_id": objectId}
	return getDefinition(filter, nil)

}

func GetDefinitionByObjectId(id primitive.ObjectID) (*Definition, error) {

	filter := bson.M{"_id": id}
	return getDefinition(filter, nil)

}

func GetNewestDefinitions(limit int) (*[]*Definition, error) {

	options := options.Find().SetSort(bson.M{"creation_date": -1}).SetLimit(int64(limit))
	return getDefinitions(bson.M{}, options)

}

func getDefinition(filter interface{}, options *options.FindOneOptions) (*Definition, error) {

	var definition Definition
	err := definitionsCollection.FindOne(dbContext, filter, options).Decode(&definition)

	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, ErrorDefinitionNotFound
		}
		return nil, err
	}

	return &definition, nil
}

func getDefinitions(filter interface{}, options *options.FindOptions) (*[]*Definition, error) {

	cursor, err := definitionsCollection.Find(dbContext, filter, options)

	if err != nil {
		defer cursor.Close(dbContext)
		return nil, err
	}

	definitions := []*Definition{}

	for cursor.Next(dbContext) {

		definition := Definition{}
		err := cursor.Decode(&definition)

		if err != nil {
			return nil, err
		}

		definitions = append(definitions, &definition)
	}

	return &definitions, nil

}
