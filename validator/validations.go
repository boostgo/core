package validator

import (
	baseValidator "github.com/go-playground/validator/v10"
	"github.com/google/uuid"
)

type Validation func() (baseValidator.Func, string)

func Register(validator *Validator, tag string, p baseValidator.Func) error {
	return validator.RegisterValidation(tag, p)
}

func RegisterBasic(validator *Validator) error {
	validations := []Validation{
		validationUUID,
		validationUndefined,
	}

	for _, validation := range validations {
		fn, tag := validation()
		if err := Register(validator, tag, fn); err != nil {
			return err
		}
	}

	return nil
}

func validationUUID() (baseValidator.Func, string) {
	const tag = "uuid"
	return func(fl baseValidator.FieldLevel) (isValid bool) {
		switch val := fl.Field().Interface().(type) {
		case string:
			_, err := uuid.Parse(val)
			if err != nil {
				return false
			}

			return true
		case *string:
			_, err := uuid.Parse(*val)
			if err != nil {
				return false
			}

			return true
		case uuid.UUID:
			isValid = val != uuid.Nil
		case *uuid.UUID:
			if fl.Field().IsNil() {
				return
			}

			isValid = *val != uuid.Nil
		}

		return
	}, tag
}

func validationUndefined() (baseValidator.Func, string) {
	const tag = "undefined"
	return func(fl baseValidator.FieldLevel) (isValid bool) {
		const undefined = "undefined"
		switch val := fl.Field().Interface().(type) {
		case string:
			return val != undefined
		case *string:
			if val == nil {
				return true
			}

			value := *val
			return value != undefined
		}

		return
	}, tag
}
