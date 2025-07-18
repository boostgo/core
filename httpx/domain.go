package httpx

import (
	"github.com/boostgo/core/errorx"
)

const (
	statusSuccess = "Success"
	statusFailure = "Failure"
)

type FailureResponse struct {
	Status     string             `json:"status"`
	StatusCode int                `json:"status_code"`
	Message    string             `json:"message"`
	Code       string             `json:"code"`
	Inner      string             `json:"inner,omitempty"`
	Context    any                `json:"context,omitempty"`
	Params     []errorx.Parameter `json:"params,omitempty"`
	RequestID  string             `json:"request_id,omitempty"`
}

func NewFailureResponse(err *errorx.Error, statusCode int, requestID string) FailureResponse {
	message := err.Message()
	if err.LocaleMessage() != "" {
		message = err.LocaleMessage()
	}

	return FailureResponse{
		Status:     statusFailure,
		Message:    message,
		Code:       err.Message(),
		Context:    err.Data(),
		Params:     err.Params(),
		StatusCode: statusCode,
		RequestID:  requestID,
	}
}

type CreatedResponse struct {
	ID any `json:"id"`
}

func NewCreatedResponse(id any) CreatedResponse {
	return CreatedResponse{
		ID: id,
	}
}

type SuccessResponse struct {
	Status     string `json:"status"`
	StatusCode int    `json:"status_code"`
	Body       any    `json:"body"`
	RequestID  string `json:"request_id,omitempty"`
}

func NewSuccessResponse(body any, statusCode int, requestID string) SuccessResponse {
	return SuccessResponse{
		Status:     statusSuccess,
		StatusCode: statusCode,
		Body:       body,
		RequestID:  requestID,
	}
}
