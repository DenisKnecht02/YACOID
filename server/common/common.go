package common

import (
	"errors"
	"time"

	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var ValidationError = errors.New("INVALID_INPUT")
var ErrorInvalidType = errors.New("INVALID_TYPE")
var ErrorNotFound = errors.New("ENTITY_NOT_FOUND")

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
	ID            primitive.ObjectID `bson:"_id" json:"-"`
	SlugId        string             `bson:"slug_id" json:"slugId"`
	SubmittedBy   primitive.ObjectID `bson:"submitted_by" json:"submittedBy"`
	SubmittedDate time.Time          `bson:"submitted_date" json:"submittedDate"`
	FirstName     string             `bson:"first_name" json:"firstName"`
	LastName      string             `bson:"last_name" json:"lastName"`
}

type Source struct {
	ID            primitive.ObjectID   `bson:"_id" json:"-"`
	SubmittedBy   primitive.ObjectID   `bson:"submitted_by" json:"submittedBy"`
	SubmittedDate time.Time            `bson:"submitted_date" json:"submittedDate"`
	Authors       []primitive.ObjectID `bson:"authors" json:"authors" validate:"required,min=1"`
}

func (author *Source) Validate(validate *validator.Validate) []string {
	return ValidateStruct(author, validate)
}

type DefinitionFilter struct {
	Title           *string      `json:"title" bson:"title" validate:"omitempty"`
	Content         *string      `json:"content" bson:"content" validate:"omitempty"`
	PublishingDates *[]time.Time `json:"publishing_dates" bson:"publishing_dates" validate:"omitempty,min=1"`
	Authors         *[]*Author   `json:"authors" bson:"authors" validate:"omitempty,min=1,dive"`
	Sources         *[]*Source   `json:"sources" bson:"sources" validate:"omitempty,min=1,dive"`
	Tags            *[]string    `json:"tags" bson:"tags" validate:"omitempty,min=1"`
}
