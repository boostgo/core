package requests

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/boostgo/core/httpx"
	"github.com/boostgo/core/reflectx"
)

type Response struct {
	request  *Request
	raw      *http.Response
	bodyBlob []byte
	isCore   bool
}

func newResponse(request *Request, resp *http.Response) *Response {
	return &Response{
		request: request,
		raw:     resp,
	}
}

func (response *Response) Raw() *http.Response {
	return response.raw
}

func (response *Response) Status() string {
	return response.raw.Status
}

func (response *Response) StatusCode() int {
	return response.raw.StatusCode
}

func (response *Response) BodyRaw() []byte {
	return response.bodyBlob
}

func (response *Response) Core(isCore bool) *Response {
	response.isCore = isCore
	return response
}

func (response *Response) Parse(export any) error {
	if response.bodyBlob == nil {
		return nil
	}

	if !reflectx.IsPointer(export) {
		return ErrExportResponseMustBePointer
	}

	if response.isCore {
		var r httpx.SuccessResponse
		if err := json.Unmarshal(response.bodyBlob, &r); err != nil {
			return err
		}

		bodyBlob, err := json.Marshal(r.Body)
		if err != nil {
			return err
		}

		if err = json.Unmarshal(bodyBlob, &export); err != nil {
			return newParseResponseBodyError(response.request.req.RequestURI, response.raw.StatusCode, response.bodyBlob).
				SetError(err)
		}

		return nil
	}

	if err := json.Unmarshal(response.bodyBlob, export); err != nil {
		return newParseResponseBodyError(response.request.req.RequestURI, response.raw.StatusCode, response.bodyBlob).
			SetError(err)
	}

	return nil
}

func (response *Response) Context(ctx context.Context) context.Context {
	return ctx
}

func (response *Response) ContentType() string {
	return response.raw.Header.Get("Content-Type")
}

func (response *Response) IsFailure() bool {
	return httpx.IsFailureCode(response.StatusCode())
}
