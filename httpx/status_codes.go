package httpx

import (
	"errors"
	"net/http"

	"github.com/boostgo/core/errorx"
)

func IsFailureCode(statusCode int) bool {
	return statusCode >= http.StatusBadRequest // >= 400
}

// StatusCodeByError - define which status code must be provided to response by error
func StatusCodeByError(err error) int {
	switch {
	case errors.Is(err, errorx.ErrBadRequest):
		return http.StatusBadRequest
	case errors.Is(err, errorx.ErrUnauthorized):
		return http.StatusUnauthorized
	case errors.Is(err, errorx.ErrPaymentRequired):
		return http.StatusPaymentRequired
	case errors.Is(err, errorx.ErrForbidden):
		return http.StatusForbidden
	case errors.Is(err, errorx.ErrNotFound):
		return http.StatusNotFound
	case errors.Is(err, errorx.ErrMethodNotAllowed):
		return http.StatusMethodNotAllowed
	case errors.Is(err, errorx.ErrNotAcceptable):
		return http.StatusNotAcceptable
	case errors.Is(err, errorx.ErrProxyAuthRequired):
		return http.StatusProxyAuthRequired
	case errors.Is(err, errorx.ErrTimeout):
		return http.StatusRequestTimeout
	case errors.Is(err, errorx.ErrConflict):
		return http.StatusConflict
	case errors.Is(err, errorx.ErrGone):
		return http.StatusGone
	case errors.Is(err, errorx.ErrLengthRequired):
		return http.StatusLengthRequired
	case errors.Is(err, errorx.ErrPreconditionFailed):
		return http.StatusPreconditionFailed
	case errors.Is(err, errorx.ErrEntityTooLarge):
		return http.StatusRequestEntityTooLarge
	case errors.Is(err, errorx.ErrURITooLong):
		return http.StatusRequestURITooLong
	case errors.Is(err, errorx.ErrUnsupportedMediaType):
		return http.StatusUnsupportedMediaType
	case errors.Is(err, errorx.ErrRangeNotSatisfiable):
		return http.StatusRequestedRangeNotSatisfiable
	case errors.Is(err, errorx.ErrExpectationFailed):
		return http.StatusExpectationFailed
	case errors.Is(err, errorx.ErrTeapot):
		return http.StatusTeapot
	case errors.Is(err, errorx.ErrMisdirectedRequest):
		return http.StatusMisdirectedRequest
	case errors.Is(err, errorx.ErrUnprocessableEntity):
		return http.StatusUnprocessableEntity
	case errors.Is(err, errorx.ErrLocked):
		return http.StatusLocked
	case errors.Is(err, errorx.ErrFailedDependency):
		return http.StatusFailedDependency
	case errors.Is(err, errorx.ErrTooEarly):
		return http.StatusTooEarly
	case errors.Is(err, errorx.ErrUpgradeRequired):
		return http.StatusUpgradeRequired
	case errors.Is(err, errorx.ErrPreconditionRequired):
		return http.StatusPreconditionRequired
	case errors.Is(err, errorx.ErrTooManyRequests):
		return http.StatusTooManyRequests
	case errors.Is(err, errorx.ErrRequestHeaderFieldsTooLarge):
		return http.StatusRequestHeaderFieldsTooLarge
	case errors.Is(err, errorx.ErrUnavailableForLegalReasons):
		return http.StatusUnavailableForLegalReasons
	case errors.Is(err, errorx.ErrInternal):
		return http.StatusInternalServerError
	case errors.Is(err, errorx.ErrNotImplemented):
		return http.StatusNotImplemented
	case errors.Is(err, errorx.ErrBadGateway):
		return http.StatusBadGateway
	case errors.Is(err, errorx.ErrServiceUnavailable):
		return http.StatusServiceUnavailable
	case errors.Is(err, errorx.ErrGatewayTimeout):
		return http.StatusGatewayTimeout
	case errors.Is(err, errorx.ErrHTTPVersionNotSupported):
		return http.StatusHTTPVersionNotSupported
	case errors.Is(err, errorx.ErrVariantAlsoNegotiates):
		return http.StatusVariantAlsoNegotiates
	case errors.Is(err, errorx.ErrInsufficientStorage):
		return http.StatusInsufficientStorage
	case errors.Is(err, errorx.ErrLoopDetected):
		return http.StatusLoopDetected
	case errors.Is(err, errorx.ErrNotExtended):
		return http.StatusNotExtended
	case errors.Is(err, errorx.ErrNetworkAuthenticationRequired):
		return http.StatusNetworkAuthenticationRequired
	default:
		return http.StatusInternalServerError
	}
}

// ErrorByStatusCode - define which error must be provided to response by status code
func ErrorByStatusCode(statusCode int) *errorx.Error {
	switch statusCode {
	case http.StatusBadRequest:
		return errorx.ErrBadRequest
	case http.StatusUnauthorized:
		return errorx.ErrUnauthorized
	case http.StatusPaymentRequired:
		return errorx.ErrPaymentRequired
	case http.StatusForbidden:
		return errorx.ErrForbidden
	case http.StatusNotFound:
		return errorx.ErrNotFound
	case http.StatusMethodNotAllowed:
		return errorx.ErrMethodNotAllowed
	case http.StatusNotAcceptable:
		return errorx.ErrNotAcceptable
	case http.StatusProxyAuthRequired:
		return errorx.ErrProxyAuthRequired
	case http.StatusRequestTimeout:
		return errorx.ErrTimeout
	case http.StatusConflict:
		return errorx.ErrConflict
	case http.StatusGone:
		return errorx.ErrGone
	case http.StatusLengthRequired:
		return errorx.ErrLengthRequired
	case http.StatusPreconditionFailed:
		return errorx.ErrPreconditionFailed
	case http.StatusRequestEntityTooLarge:
		return errorx.ErrEntityTooLarge
	case http.StatusRequestURITooLong:
		return errorx.ErrURITooLong
	case http.StatusUnsupportedMediaType:
		return errorx.ErrUnsupportedMediaType
	case http.StatusRequestedRangeNotSatisfiable:
		return errorx.ErrRangeNotSatisfiable
	case http.StatusExpectationFailed:
		return errorx.ErrExpectationFailed
	case http.StatusTeapot:
		return errorx.ErrTeapot
	case http.StatusMisdirectedRequest:
		return errorx.ErrMisdirectedRequest
	case http.StatusUnprocessableEntity:
		return errorx.ErrUnprocessableEntity
	case http.StatusLocked:
		return errorx.ErrLocked
	case http.StatusFailedDependency:
		return errorx.ErrFailedDependency
	case http.StatusTooEarly:
		return errorx.ErrTooEarly
	case http.StatusUpgradeRequired:
		return errorx.ErrUpgradeRequired
	case http.StatusPreconditionRequired:
		return errorx.ErrPreconditionRequired
	case http.StatusTooManyRequests:
		return errorx.ErrTooManyRequests
	case http.StatusRequestHeaderFieldsTooLarge:
		return errorx.ErrRequestHeaderFieldsTooLarge
	case http.StatusUnavailableForLegalReasons:
		return errorx.ErrUnavailableForLegalReasons
	case http.StatusInternalServerError:
		return errorx.ErrInternal
	case http.StatusNotImplemented:
		return errorx.ErrNotImplemented
	case http.StatusBadGateway:
		return errorx.ErrBadGateway
	case http.StatusServiceUnavailable:
		return errorx.ErrServiceUnavailable
	case http.StatusGatewayTimeout:
		return errorx.ErrGatewayTimeout
	case http.StatusHTTPVersionNotSupported:
		return errorx.ErrHTTPVersionNotSupported
	case http.StatusVariantAlsoNegotiates:
		return errorx.ErrVariantAlsoNegotiates
	case http.StatusInsufficientStorage:
		return errorx.ErrInsufficientStorage
	case http.StatusLoopDetected:
		return errorx.ErrLoopDetected
	case http.StatusNotExtended:
		return errorx.ErrNotExtended
	case http.StatusNetworkAuthenticationRequired:
		return errorx.ErrNetworkAuthenticationRequired
	default:
		return errorx.ErrInternal
	}
}
