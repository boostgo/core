package validator

import (
	"errors"
	"sync"

	"github.com/boostgo/core/errorx"
	
	baseValidator "github.com/go-playground/validator/v10"
)

var (
	_validator *Validator
	_once      sync.Once
)

func Get() *Validator {
	_once.Do(func() {
		_validator = New()

		if err := RegisterBasic(_validator); err != nil {
			panic(err)
		}
	})

	return _validator
}

type Validator struct {
	*baseValidator.Validate
	turnOff bool
}

func New() *Validator {
	return &Validator{
		Validate: baseValidator.New(),
	}
}

func (validator *Validator) TurnOff() *Validator {
	validator.turnOff = true
	return validator
}

func (validator *Validator) Struct(object any) error {
	if validator.turnOff {
		return nil
	}

	validateError := validator.Validate.Struct(object)
	if validateError == nil {
		return nil
	}

	err := ErrModelValidation.SetError(errorx.ErrUnprocessableEntity)

	var validationErrors baseValidator.ValidationErrors
	ok := errors.As(validateError, &validationErrors)
	if !ok {
		return err.SetError(validateError)
	}

	if len(validationErrors) == 0 {
		return nil
	}

	validations := make([]string, 0, len(validationErrors))
	for _, validationError := range validationErrors {
		validations = append(validations, validationError.Error())
	}

	return err.SetData(validationContext{
		Validations: validations,
	})
}

func (validator *Validator) Var(variable any, tag string) error {
	if validator.turnOff {
		return nil
	}

	err := ErrVariableValidation.SetError(errorx.ErrUnprocessableEntity)

	validateError := validator.Validate.Var(variable, tag)
	if err == nil {
		return err
	}

	var validationErrors baseValidator.ValidationErrors
	ok := errors.As(validateError, &validationErrors)
	if !ok {
		return err.SetData(validationContext{
			Validation: validateError.Error(),
		})
	}

	if len(validationErrors) == 0 {
		return nil
	}

	validations := make([]string, 0, len(validationErrors))
	for _, validationError := range validationErrors {
		validations = append(validations, validationError.Error())
	}

	return err.SetData(validationContext{
		Validations: validations,
	})
}
