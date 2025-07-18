package retry

import (
	"context"
	"errors"
	"fmt"
	"testing"
	"time"
)

func TestRetrySuccess(t *testing.T) {
	attempts := 0
	err := Retry(context.Background(), func(ctx context.Context) error {
		attempts++
		if attempts < 3 {
			return errors.New("temporary error")
		}
		return nil
	}, Options{
		Policy: NewFixedDelay(10*time.Millisecond, 5),
	})

	if err != nil {
		t.Fatalf("expected success, got error: %v", err)
	}

	if attempts != 3 {
		t.Fatalf("expected 3 attempts, got %d", attempts)
	}
}

func TestRetryMaxAttemptsExceeded(t *testing.T) {
	attempts := 0
	err := Retry(context.Background(), func(ctx context.Context) error {
		attempts++
		return errors.New("always fails")
	}, Options{
		Policy: NewFixedDelay(10*time.Millisecond, 3),
	})

	if err == nil {
		t.Fatal("expected error, got nil")
	}

	var retryErr *Error
	if !errors.As(err, &retryErr) {
		t.Fatalf("expected retry.Error, got %T", err)
	}

	if retryErr.Attempts != 3 {
		t.Fatalf("expected 3 attempts, got %d", retryErr.Attempts)
	}

	if !errors.Is(err, ErrMaxRetriesExceeded) {
		t.Fatal("expected ErrMaxRetriesExceeded")
	}
}

func TestRetryNonRetryableError(t *testing.T) {
	nonRetryableErr := errors.New("non-retryable")
	attempts := 0

	err := Retry(context.Background(), func(ctx context.Context) error {
		attempts++
		return nonRetryableErr
	}, Options{
		Policy: NewFixedDelay(10*time.Millisecond, 3),
		RetryIf: func(err error) bool {
			return !errors.Is(err, nonRetryableErr)
		},
	})

	if err == nil {
		t.Fatal("expected error, got nil")
	}

	if attempts != 1 {
		t.Fatalf("expected 1 attempt for non-retryable error, got %d", attempts)
	}

	if !errors.Is(err, ErrNonRetryable) {
		t.Fatal("expected ErrNonRetryable")
	}
}

func TestRetryContextCancellation(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	attempts := 0

	go func() {
		time.Sleep(50 * time.Millisecond)
		cancel()
	}()

	err := Retry(ctx, func(ctx context.Context) error {
		attempts++
		return errors.New("temporary error")
	}, Options{
		Policy: NewFixedDelay(100*time.Millisecond, 5),
	})

	if err == nil {
		t.Fatal("expected error, got nil")
	}

	if !errors.Is(err, context.Canceled) {
		t.Fatal("expected context.Canceled")
	}

	// Should have attempted at least once, but not all 5 times
	if attempts == 0 || attempts >= 5 {
		t.Fatalf("unexpected number of attempts: %d", attempts)
	}
}

func TestRetryOnRetryCallback(t *testing.T) {
	var callbackCalls []int
	attempts := 0

	err := Retry(context.Background(), func(ctx context.Context) error {
		attempts++
		if attempts < 3 {
			return fmt.Errorf("error on attempt %d", attempts)
		}
		return nil
	}, Options{
		Policy: NewFixedDelay(10*time.Millisecond, 5),
		OnRetry: func(attempt int, err error) {
			callbackCalls = append(callbackCalls, attempt)
		},
	})

	if err != nil {
		t.Fatalf("expected success, got error: %v", err)
	}

	expectedCalls := []int{1, 2}
	if len(callbackCalls) != len(expectedCalls) {
		t.Fatalf("expected %d callback calls, got %d", len(expectedCalls), len(callbackCalls))
	}

	for i, expected := range expectedCalls {
		if callbackCalls[i] != expected {
			t.Fatalf("callback call %d: expected attempt %d, got %d", i, expected, callbackCalls[i])
		}
	}
}

func TestExponentialBackoffPolicy(t *testing.T) {
	policy := NewExponentialBackoff(100*time.Millisecond, 2*time.Second, 5)

	tests := []struct {
		attempt     int
		expectedMin time.Duration
		expectedMax time.Duration
	}{
		{1, 100 * time.Millisecond, 100 * time.Millisecond},
		{2, 200 * time.Millisecond, 200 * time.Millisecond},
		{3, 400 * time.Millisecond, 400 * time.Millisecond},
		{4, 800 * time.Millisecond, 800 * time.Millisecond},
		{5, 1600 * time.Millisecond, 1600 * time.Millisecond},
		{6, 2 * time.Second, 2 * time.Second}, // capped at max
	}

	for _, tt := range tests {
		delay := policy.NextDelay(tt.attempt)
		if delay < tt.expectedMin || delay > tt.expectedMax {
			t.Errorf("attempt %d: expected delay between %v and %v, got %v",
				tt.attempt, tt.expectedMin, tt.expectedMax, delay)
		}
	}
}

func TestExponentialBackoffWithJitterPolicy(t *testing.T) {
	policy := NewExponentialBackoffWithJitter(100*time.Millisecond, 2*time.Second, 5, 0.5)

	// Test multiple times to account for randomness
	for i := 0; i < 10; i++ {
		delay := policy.NextDelay(3) // 3rd attempt should be ~400ms base

		// With 0.5 jitter factor, delay should be between 200ms and 600ms
		minExpected := 200 * time.Millisecond
		maxExpected := 600 * time.Millisecond

		if delay < minExpected || delay > maxExpected {
			t.Errorf("iteration %d: expected delay between %v and %v, got %v",
				i, minExpected, maxExpected, delay)
		}
	}
}

func TestDoConvenienceFunction(t *testing.T) {
	attempts := 0
	err := Try(context.Background(), func(ctx context.Context) error {
		attempts++
		if attempts < 2 {
			return errors.New("temporary")
		}
		return nil
	}, Options{
		Policy: NewFixedDelay(10*time.Millisecond, 3),
	})

	if err != nil {
		t.Fatalf("expected success, got error: %v", err)
	}

	if attempts != 2 {
		t.Fatalf("expected 2 attempts, got %d", attempts)
	}
}

func TestDoWithDataConvenienceFunction(t *testing.T) {
	attempts := 0
	result, err := DoWithData(context.Background(), func(ctx context.Context) (string, error) {
		attempts++
		if attempts < 2 {
			return "", errors.New("temporary")
		}
		return "success", nil
	}, Options{
		Policy: NewFixedDelay(10*time.Millisecond, 3),
	})

	if err != nil {
		t.Fatalf("expected success, got error: %v", err)
	}

	if result != "success" {
		t.Fatalf("expected 'success', got %q", result)
	}

	if attempts != 2 {
		t.Fatalf("expected 2 attempts, got %d", attempts)
	}
}

// Example of custom retryable error
type temporaryError struct {
	message string
}

func (e *temporaryError) Error() string {
	return e.message
}

func (e *temporaryError) Retryable() bool {
	return true
}

func TestCustomRetryableError(t *testing.T) {
	attempts := 0
	permanentErr := errors.New("permanent error")

	err := Retry(context.Background(), func(ctx context.Context) error {
		attempts++
		if attempts < 2 {
			return &temporaryError{message: "temporary error"}
		}
		return permanentErr
	}, Options{
		Policy: NewFixedDelay(10*time.Millisecond, 3),
		RetryIf: func(err error) bool {
			// Only retry if it implements Retryable and returns true
			type retryable interface {
				Retryable() bool
			}
			if r, ok := err.(retryable); ok {
				return r.Retryable()
			}
			// Don't retry plain errors
			return false
		},
	})

	if err == nil {
		t.Fatal("expected error, got nil")
	}

	// Should have tried twice: once for temporary error, once for permanent
	if attempts != 2 {
		t.Fatalf("expected 2 attempts, got %d", attempts)
	}

	if !errors.Is(err, permanentErr) {
		t.Fatal("expected permanent error")
	}
}
