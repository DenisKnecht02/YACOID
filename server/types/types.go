package types

import (
	"time"
	"yacoid_server/common"

	"github.com/go-playground/validator/v10"
)

type SubmitDefinitionRequest struct {
	Title          string    `json:"title" validate:"required,min=1"`
	Content        string    `json:"content" validate:"required,min=1"`
	Source         string    `json:"source" validate:"required"`
	PublishingDate time.Time `json:"publishingDate" validate:"required,ISO8601date"`
	Tags           *[]string `json:"tags" validate:"required,min=1"`
}

func (author *SubmitDefinitionRequest) Validate(validate *validator.Validate) []string {
	return common.ValidateStruct(author, validate)
}

type DefinitionPageRequest struct {
	PageSize int                      `json:"pageSize" validate:"required"`
	Page     int                      `json:"page" validate:"required,min=1"`
	Filter   *common.DefinitionFilter `json:"filter" validate:"omitempty,dive"`
	Sort     *interface{}             `json:"sort"`
}

func (DefinitionPageRequest *DefinitionPageRequest) Validate(validate *validator.Validate) []string {
	return common.ValidateStruct(DefinitionPageRequest, validate)
}

type RejectRequest struct {
	ID      string `json:"id" validate:"required"`
	Content string `json:"content" validate:"required,min=1"`
}

func (rejection *RejectRequest) Validate(validate *validator.Validate) []string {
	return common.ValidateStruct(rejection, validate)
}

type ChangeDefinitionRequest struct {
	ID             string         `json:"id" validate:"required"`
	Title          *string        `json:"title"`
	Content        *string        `json:"content"`
	Source         *common.Source `json:"source" validate:"omitempty,dive"`
	PublishingDate *time.Time     `json:"publishingDate" validate:"omitempty,ISO8601date"`
	Tags           *[]string      `json:"tags"`
}

func (rejection *ChangeDefinitionRequest) Validate(validate *validator.Validate) []string {
	return common.ValidateStruct(rejection, validate)
}

type CreateAuthorRequest struct {
	FirstName string `json:"firstName" validate:"required,min=1"`
	LastName  string `json:"lastName" validate:"required,min=1"`
}

func (rejection *CreateAuthorRequest) Validate(validate *validator.Validate) []string {
	return common.ValidateStruct(rejection, validate)
}

type CreateSourceRequest struct {
	Authors []string `bson:"authors" json:"authors" validate:"required,min=1"`
}

func (rejection *CreateSourceRequest) Validate(validate *validator.Validate) []string {
	return common.ValidateStruct(rejection, validate)
}
