package common

import (
	"errors"

	"github.com/go-playground/validator/v10"
)

var ValidationError = errors.New("INVALID_INPUT")

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
