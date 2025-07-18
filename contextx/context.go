package contextx

import (
	"context"
	"errors"

	"github.com/boostgo/core/errorx"
)

func Timeout(ctx context.Context) error {
	if ctx == nil {
		return nil
	}

	err := ctx.Err()
	if err == nil {
		return nil
	}

	if errors.Is(err, context.DeadlineExceeded) || errors.Is(err, context.Canceled) {
		return errorx.ErrTimeout.SetError(err)
	}

	return nil
}

func Any(ctx context.Context) error {
	if ctx == nil {
		return nil
	}

	err := ctx.Err()
	if err != nil {
		return err
	}

	return nil
}

func Validate(ctx context.Context) error {
	err := Timeout(ctx)
	if err != nil {
		return err
	}

	err = Any(ctx)
	if err != nil {
		return err
	}

	return nil
}
