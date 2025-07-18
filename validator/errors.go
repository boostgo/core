package validator

import "github.com/boostgo/core/errorx"

var (
	ErrModelValidation    = errorx.New("validator.model")
	ErrVariableValidation = errorx.New("validator.variable")
)

type validationContext struct {
	Validation  string   `json:"validation,omitempty"`
	Validations []string `json:"validations,omitempty"`
}
