package data

import (
	"fmt"

	"github.com/go-playground/validator/v10"
)

type Validation struct {
	validate *validator.Validate
}

func NewValidation() *Validation {
	validate := validator.New()
	return &Validation{validate}
}

type VError struct {
	validator.FieldError
}

func (v VError) Error() string {
	return fmt.Sprintf(
		"Key: '%s' Error: Field validation for '%s' failed on the '%s' tag",
		v.Namespace(),
		v.Field(),
		v.Tag(),
	)
}

type VErrors []VError

func (v VErrors) Errors() []string {
	errs := []string{}
	for _, err := range v {
		errs = append(errs, err.Error())
	}

	return errs
}

func (v *Validation) Validate(i interface{}) VErrors {
	errs := v.validate.Struct(i)
	if errs == nil {
		return nil
	}

	var returnErrs []VError
	for _, err := range errs.(validator.ValidationErrors) {
		ve := VError{err.(validator.FieldError)}
		returnErrs = append(returnErrs, ve)
	}

	return returnErrs
}
