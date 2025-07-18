package errorx

import (
	"runtime/debug"

	"github.com/boostgo/core/convert"
)

var (
	ErrBadRequest                  = New("bad_request")
	ErrUnauthorized                = New("unauthorized")
	ErrPaymentRequired             = New("payment_required")
	ErrForbidden                   = New("forbidden")
	ErrNotFound                    = New("not_found")
	ErrMethodNotAllowed            = New("method_not_allowed")
	ErrNotAcceptable               = New("not_acceptable")
	ErrProxyAuthRequired           = New("proxy_auth_required")
	ErrTimeout                     = New("timeout")
	ErrConflict                    = New("conflict")
	ErrGone                        = New("gone")
	ErrLengthRequired              = New("length_required")
	ErrPreconditionFailed          = New("precondition_failed")
	ErrEntityTooLarge              = New("entity_too_large")
	ErrURITooLong                  = New("uri_too_long")
	ErrUnsupportedMediaType        = New("unsupported_media_type")
	ErrRangeNotSatisfiable         = New("range_not_satisfiable")
	ErrExpectationFailed           = New("expectation_failed")
	ErrTeapot                      = New("teapot")
	ErrMisdirectedRequest          = New("misdirected_request")
	ErrUnprocessableEntity         = New("unprocessable_entity")
	ErrLocked                      = New("locked")
	ErrFailedDependency            = New("failed_dependency")
	ErrTooEarly                    = New("too_early")
	ErrUpgradeRequired             = New("upgrade_required")
	ErrPreconditionRequired        = New("precondition_required")
	ErrTooManyRequests             = New("too_many_requests")
	ErrRequestHeaderFieldsTooLarge = New("request_header_fields_too_large")
	ErrUnavailableForLegalReasons  = New("unavailable_for_legal_reasons")

	ErrInternal                      = New("internal")
	ErrNotImplemented                = New("not_implemented")
	ErrBadGateway                    = New("bad_gateway")
	ErrServiceUnavailable            = New("service_unavailable")
	ErrGatewayTimeout                = New("gateway_timeout")
	ErrHTTPVersionNotSupported       = New("http_version_not_supported")
	ErrVariantAlsoNegotiates         = New("variant_also_negotiates")
	ErrInsufficientStorage           = New("insufficient_storage")
	ErrLoopDetected                  = New("loop_detected")
	ErrNotExtended                   = New("not_extended")
	ErrNetworkAuthenticationRequired = New("network_authentication_required")
)

var ErrPanicRecover = New("panic_recover")

type panicRecoverContext struct {
	Stack string `json:"stack"`
}

func NewPanicRecoverError() *Error {
	return ErrPanicRecover.SetData(panicRecoverContext{
		Stack: convert.StringFromBytes(debug.Stack()),
	})
}
