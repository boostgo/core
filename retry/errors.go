package retry

import (
	"context"
	"errors"
	"time"

	"github.com/boostgo/core/errorx"
)

var (
	ErrMaxRetriesExceeded = errorx.New("retry.maximum_retry_attempts")
	ErrNonRetryable       = errorx.New("retry.non_retryable_error")
)

// Error wraps the original error with retry metadata
type Error struct {
	LastError  error
	Attempts   int
	TotalDelay time.Duration
}

func (e *Error) Error() string {
	return e.LastError.Error()
}

func (e *Error) Unwrap() error {
	return e.LastError
}

// IsRetryable is a helper function to check if an error should be retried
// This is a default implementation that can be overridden in RetryOptions
func IsRetryable(err error) bool {
	if err == nil {
		return false
	}

	// Don't retry context errors
	if errors.Is(err, context.Canceled) || errors.Is(err, context.DeadlineExceeded) {
		return false
	}

	// Check if error implements a Retryable interface
	type retryable interface {
		Retryable() bool
	}

	if r, ok := err.(retryable); ok {
		return r.Retryable()
	}

	// By default, retry all other errors
	return true
}
