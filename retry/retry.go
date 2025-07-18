package retry

import (
	"context"
	"errors"
	"time"

	"github.com/boostgo/core/errorx"
)

type RetryableFunc func(ctx context.Context) error

// Policy defines the interface for retry strategies
type Policy interface {
	// NextDelay returns the delay before the next retry
	// attempt starts from 1
	NextDelay(attempt int) time.Duration

	// MaxAttempts returns the maximum number of retry attempts
	MaxAttempts() int
}

// Options contains configuration for retry behavior
type Options struct {
	// Policy defines the retry strategy
	Policy Policy

	// RetryIf is a function that determines if an error should trigger a retry
	// If nil, all errors are retried
	RetryIf func(error) bool

	// OnRetry is called before each retry
	// attempt starts from 1
	OnRetry func(attempt int, err error)
}

// Retry executes the given function with retry logic based on the provided options
func Retry(ctx context.Context, fn RetryableFunc, options ...Options) error {
	var opts Options
	if len(options) > 0 {
		opts = options[0]
	}

	if opts.Policy == nil {
		// Default policy: 3 attempts with 100-millisecond fixed delay
		opts.Policy = NewFixedDelay(time.Millisecond*100, 3)
	}

	if opts.RetryIf == nil {
		opts.RetryIf = IsRetryable
	}

	if ctx == nil {
		ctx = context.Background()
	}

	var lastErr error
	var totalDelay time.Duration
	maxAttempts := opts.Policy.MaxAttempts()

	for attempt := 1; attempt <= maxAttempts; attempt++ {
		// Check if context is already cancelled
		select {
		case <-ctx.Done():
			return &Error{
				LastError:  ctx.Err(),
				Attempts:   attempt - 1,
				TotalDelay: totalDelay,
			}
		default:
		}

		// Execute the function
		err := fn(ctx)
		if err == nil {
			return nil
		}

		lastErr = err

		// Check if error is retryable
		if !opts.RetryIf(err) {
			return &Error{
				LastError:  errors.Join(ErrNonRetryable, err),
				Attempts:   attempt,
				TotalDelay: totalDelay,
			}
		}

		// Don't retry if this was the last attempt
		if attempt == maxAttempts {
			break
		}

		// Call OnRetry callback if provided
		if opts.OnRetry != nil {
			opts.OnRetry(attempt, err)
		}

		// Calculate delay for next attempt
		delay := opts.Policy.NextDelay(attempt)
		totalDelay += delay

		// Wait with context cancellation support
		timer := time.NewTimer(delay)
		select {
		case <-timer.C:
			// Continue to next attempt
		case <-ctx.Done():
			timer.Stop()
			return &Error{
				LastError:  errors.Join(lastErr, ctx.Err()),
				Attempts:   attempt,
				TotalDelay: totalDelay,
			}
		}
	}

	return &Error{
		LastError:  errors.Join(ErrMaxRetriesExceeded, lastErr),
		Attempts:   maxAttempts,
		TotalDelay: totalDelay,
	}
}

// Try is a convenience function that calls Retry but with try (panic tolerance)
func Try(ctx context.Context, fn RetryableFunc, opts ...Options) error {
	return errorx.Try(func() error {
		return Retry(ctx, fn, opts...)
	})
}

// DoWithData is a generic convenience function for operations that return data
func DoWithData[T any](ctx context.Context, fn func(ctx context.Context) (T, error), opts ...Options) (T, error) {
	var result T
	err := Retry(ctx, func(ctx context.Context) error {
		var fnErr error
		result, fnErr = fn(ctx)
		return fnErr
	}, opts...)
	return result, err
}
