package errs

import (
	"errors"

	"github.com/boostgo/core/errorx"
	"google.golang.org/grpc/codes"
)

var (
	ErrCanceled           = errorx.New("canceled")
	ErrUnknown            = errorx.New("unknown")
	ErrInvalidArgument    = errorx.New("invalid_argument")
	ErrDeadlineExceeded   = errorx.New("deadline_exceeded")
	ErrNotFound           = errorx.New("not_found")
	ErrAlreadyExist       = errorx.New("already_exist")
	ErrPermissionDenied   = errorx.New("permission_denied")
	ErrResourceExhausted  = errorx.New("resource_exhausted")
	ErrFailedPrecondition = errorx.New("failed_precondition")
	ErrAborted            = errorx.New("aborted")
	ErrOutOfRange         = errorx.New("out_of_range")
	ErrUnimplemented      = errorx.New("unimplemented")
	ErrInternal           = errorx.New("internal")
	ErrUnavailable        = errorx.New("unavailable")
	ErrDataLoss           = errorx.New("data_loss")
	ErrUnauthenticated    = errorx.New("unauthenticated")
)

func Code(err error) codes.Code {
	switch {
	case errors.Is(err, ErrCanceled):
		return codes.Canceled
	case errors.Is(err, ErrUnknown):
		return codes.Unknown
	case errors.Is(err, ErrInvalidArgument),
		errors.Is(err, errorx.ErrBadRequest):
		return codes.InvalidArgument
	case errors.Is(err, ErrDeadlineExceeded),
		errors.Is(err, errorx.ErrTimeout):
		return codes.DeadlineExceeded
	case errors.Is(err, ErrNotFound),
		errors.Is(err, errorx.ErrNotFound):
		return codes.NotFound
	case errors.Is(err, ErrAlreadyExist),
		errors.Is(err, errorx.ErrConflict):
		return codes.AlreadyExists
	case errors.Is(err, ErrPermissionDenied),
		errors.Is(err, errorx.ErrForbidden):
		return codes.PermissionDenied
	case errors.Is(err, ErrResourceExhausted),
		errors.Is(err, errorx.ErrTooManyRequests):
		return codes.ResourceExhausted
	case errors.Is(err, ErrFailedPrecondition):
		return codes.FailedPrecondition
	case errors.Is(err, ErrAborted):
		return codes.Aborted
	case errors.Is(err, ErrOutOfRange):
		return codes.OutOfRange
	case errors.Is(err, ErrUnimplemented):
		return codes.Unimplemented
	case errors.Is(err, ErrInternal),
		errors.Is(err, errorx.ErrInternal):
		return codes.Internal
	case errors.Is(err, ErrUnavailable),
		errors.Is(err, errorx.ErrServiceUnavailable),
		errors.Is(err, errorx.ErrGatewayTimeout),
		errors.Is(err, errorx.ErrBadGateway):
		return codes.Unavailable
	case errors.Is(err, ErrDataLoss):
		return codes.DataLoss
	case errors.Is(err, ErrUnauthenticated),
		errors.Is(err, errorx.ErrUnauthorized):
		return codes.Unauthenticated
	default:
		return codes.Internal
	}
}
