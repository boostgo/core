// Package errorx provides custom error implementation.
// Features:
// - List of popular errors list (copied from HTTP codes).
// - Nested errors.
// - Joined errors.
// - Context data. Any type of keys & values.
// - Implement errors package Is and Unwrap functions.
// - Wrap error and collect messages & types to one list.
// - Copy [Error].
package errorx

import (
	"errors"
	"strings"
)

// Error is custom error which implements built-in error interface.
//
// Struct contains hierarchy of error messages and their types; context (map) and inner error.
//
//	For example, error types could be like "User Handler - User Usecase - User Repository - SQL"
//	It means that first error created on "SQL" level (sql, sqlx or any other module), then error wrapped
//	by "User Repository" level, then "User Usecase" level and so on.
type Error struct {
	message       string
	localeMessage string
	inner         error
	data          any
	params        []Parameter
	noCopy        bool
}

// New creates new Error object with provided message
func New(message string) *Error {
	return &Error{
		message: message,
	}
}

// Extend copies provided err to the new one.
//
// Inner errors sets inside new error as one inner error.
//
// If inner errors contains only 1 error it will be 1 error, if errors more than 1, it will be "Join error"
func Extend(err error) *Error {
	var extended *Error
	var message string
	var inner error

	if errors.As(err, &extended) {
		if extended.noCopy {
			return extended
		}

		message = extended.message
		inner = extended.inner
	} else {
		message = err.Error()
	}

	return &Error{
		message: message,
		inner:   inner,
	}
}

// String returns string representation of current error.
//
// Method uses string builder and it's grow method.
//
// Method prints: types, messages and context
func (e *Error) String() string {
	builder := strings.Builder{}
	builder.WriteString(e.message)

	inner := e.inner
	for inner != nil {
		var convertedInner *Error
		if errors.As(inner, &convertedInner) {
			builder.WriteString(": ")
			builder.WriteString(convertedInner.message)

			inner = convertedInner.inner
			continue
		} else if inner != nil {
			builder.WriteString(": ")
			builder.WriteString(inner.Error())
		}

		break
	}

	return builder.String()
}

// Error returns result of String() method
func (e *Error) Error() string {
	return e.String()
}

// Unwrap takes inner error and try to take inside wrapped errors.
//
// Method works only for custom errors, otherwise to result error slice will be added just inner error by itself
func (e *Error) Unwrap() error {
	return e.inner
}

// Is compares current error with provided target error.
//
// By comparing errors method check if provided error is custom or not:
//
//	if custom - use equals method.
//	If not custom - unwrap current error and compare unwrapped inner errors with provided target
func (e *Error) Is(target error) bool {
	if target == nil {
		return false
	}

	// Проверяем, совпадает ли текущая ошибка с target
	var convertedTarget *Error
	if errors.As(target, &convertedTarget) {
		if e.message == convertedTarget.message {
			return true
		}
	} else {
		// Если target не является *Error, сравниваем по сообщению
		if e.message == target.Error() {
			return true
		}
	}

	// Проверяем, является ли target joinErrors
	var joinedTarget *joinErrors
	if errors.As(target, &joinedTarget) {
		// Проверяем, содержит ли joinErrors текущую ошибку
		for _, err := range joinedTarget.errors {
			if errors.Is(e, err) {
				return true
			}
		}
	}

	// Рекурсивно проверяем внутренние ошибки
	if e.inner != nil {
		// Если внутренняя ошибка является joinErrors
		var innerJoined *joinErrors
		if errors.As(e.inner, &innerJoined) {
			for _, err := range innerJoined.errors {
				if errors.Is(err, target) {
					return true
				}
			}
		} else {
			// Обычная проверка внутренней ошибки
			if errors.Is(e.inner, target) {
				return true
			}
		}
	}

	return false
}

// Message returns message.
func (e *Error) Message() string {
	return e.message
}

// SetLocaleMessage sets locale message.
func (e *Error) SetLocaleMessage(message string) *Error {
	target := Extend(e)
	target.localeMessage = message
	return target
}

// LocaleMessage returns locale message.
func (e *Error) LocaleMessage() string {
	return e.localeMessage
}

// SetError sets inner error.
//
// If inner errors more than 1 it will be "join error", if error is 1 it will be provided by itself
func (e *Error) SetError(err ...error) *Error {
	target := Extend(e)

	if len(err) == 0 {
		return target
	}

	if len(err) == 1 && err[0] != nil {
		target.inner = err[0]
		return target
	}

	target.inner = Join(err...)
	return target
}

// Inner returns inner error.
func (e *Error) Inner() error {
	return e.inner
}

// SetData sets context data (any type).
func (e *Error) SetData(data any) *Error {
	if data == nil {
		return e
	}

	target := Extend(e)

	target.data = data
	return target
}

// Data returns current error context (map)
func (e *Error) Data() any {
	return e.data
}

// AddParam append new key-value param
func (e *Error) AddParam(key string, value any) *Error {
	target := Extend(e)
	target.params = append(target.params, Parameter{
		Key:   key,
		Value: value,
	})
	return target
}

// SetParams append slice of key-value params
func (e *Error) SetParams(params []Parameter) *Error {
	target := Extend(e)
	target.params = append(target.params, params...)
	return target
}

func (e *Error) NoCopy() *Error {
	e.noCopy = true
	return e
}

// Params return params.
func (e *Error) Params() []Parameter {
	return e.params
}

// Wrap convert provided error to custom with the provided error type and message.
//
// If provided error is built-in (default), then it will be converted to custom.
//
// If it is already custom, just take custom and set to it one more type & message
func Wrap(err, wrapper error, data ...any) error {
	if err == nil || wrapper == nil {
		return nil
	}

	var setData any
	if len(data) > 0 && data[0] != nil {
		setData = data[0]
	}

	var convertedWrapper *Error
	if errors.As(wrapper, &convertedWrapper) {
		return Extend(convertedWrapper).
			SetError(err).
			SetData(setData)
	}

	return New(wrapper.Error()).
		SetError(err).
		SetData(setData)
}
