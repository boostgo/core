package httpx

import "github.com/boostgo/core/errorx"

var (
	ErrParseRequestBody = errorx.New("request_parse_body").SetError(errorx.ErrBadRequest)
	ErrReadFormFile     = errorx.New("read_form_file").SetError(errorx.ErrBadRequest)
	ErrOpenFormFile     = errorx.New("open_form_file").SetError(errorx.ErrBadRequest)

	ErrParseIntParam   = errorx.New("param.parse_int").SetError(errorx.ErrBadRequest)
	ErrParseFloatParam = errorx.New("param.parse_float").SetError(errorx.ErrBadRequest)
	ErrParseUUIDParam  = errorx.New("param.parse_uuid").SetError(errorx.ErrBadRequest)

	ErrPathParamIsEmpty = errorx.New("path_param_empty").SetError(errorx.ErrBadRequest)

	ErrRouteNotFound = errorx.New("route_not_found").SetError(errorx.ErrNotFound)
	ErrStartServer   = errorx.New("server_start")
)

type parseContext struct {
	Value string `json:"value"`
}

type emptyParamContext struct {
	ParamName string `json:"param_name"`
}

func NewEmptyPathParamError(paramName string) error {
	return ErrPathParamIsEmpty.SetData(emptyParamContext{
		ParamName: paramName,
	})
}
