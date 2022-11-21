package main

import (
	"fmt"
	"net/mail"

	"github.com/go-playground/validator/v10"
)

type Validation struct {
	validate *validator.Validate
}

func NewValidation() *Validation {
	validate := validator.New()
	validate.RegisterValidation("email", validateEmail)

	return &Validation{validate}
}

type ValidationError struct {
	validator.FieldError
}

func (v ValidationError) Error() string {
	return fmt.Sprintf(
		"Key: '%s' Error: Field validation for '%s' failed on the '%s' tag",
		v.Namespace(),
		v.Field(),
		v.Tag(),
	)
}

type ValidationErrors []ValidationError

func (v ValidationErrors) Errors() []string {
	errs := []string{}
	for _, err := range v {
		errs = append(errs, err.Error())
	}

	return errs
}

func (v *Validation) Validate(i interface{}) ValidationErrors {
	errs := v.validate.Struct(i)
	if errs == nil {
		return nil
	}

	var returnErrs []ValidationError
	for _, err := range errs.(validator.ValidationErrors) {
		ve := ValidationError{err.(validator.FieldError)}
		returnErrs = append(returnErrs, ve)
	}

	return returnErrs
}

func validateEmail(fl validator.FieldLevel) bool {
	_, err := mail.ParseAddress(fl.Field().String())
	return err == nil
}
