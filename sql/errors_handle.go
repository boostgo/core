package sql

import (
	"database/sql"
	"errors"
	"github.com/boostgo/core/sql/duplicate"
	"github.com/lib/pq"

	"github.com/boostgo/core/errorx"
)

// NotFound check if provided error is not found error
func NotFound(err error) bool {
	return errors.Is(err, sql.ErrNoRows) || errors.Is(err, errorx.ErrNotFound)
}

func Duplicate(err error) (bool, string) {
	var pqErr *pq.Error
	if !errors.As(err, &pqErr) {
		return false, ""
	}

	if duplicateErr := duplicate.Handle(pqErr); duplicateErr != nil {
		return true, duplicateErr.Field
	}

	return false, ""
}

func HandleError(err error) error {
	var pqErr *pq.Error
	if !errors.As(err, &pqErr) {
		return nil
	}

	switch pqErr.Code {
	case "23505":
		if duplicateErr := duplicate.Handle(pqErr); duplicateErr != nil {
			return NewDuplicateError(duplicateErr)
		}

		return nil
	case "23503":
		return NewForeignKeyViolationError(pqErr.Detail)
	case "23502":
		return NewNotNullError(pqErr.Column)
	default:
		return err
	}
}
