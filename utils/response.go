package utils

import "net/http"

type Response struct {
	Code    int           `json:"code,omitempty"`
	Message string        `json:"message,omitempty"`
	Data    interface{}   `json:"data,omitempty"`
	Error   string        `json:"error,omitempty"`
	Errors  []interface{} `json:"errors,omitempty"`
}

func NewResponse(code int, message string, data interface{}) *Response {
	return &Response{
		Code:    code,
		Message: message,
		Data:    data,
	}
}

func NewErrorResponse(code int, message string, err string, errs []interface{}) *Response {
	return &Response{
		Code:    code,
		Message: message,
		Error:   err,
		Errors:  errs,
		Data:    nil,
	}
}

func NewErrorUnableToReadRequestBody() *Response {
	return NewErrorResponse(http.StatusBadRequest, "Unable to read request body", "Bad Request", nil)
}

func NewErrorUnableToParseRequestBody() *Response {
	return NewErrorResponse(http.StatusBadRequest, "Unable to parse body", "Bad Request", nil)
}

func NewErrorInvalidRequestBody(errs []interface{}) *Response {
	return NewErrorResponse(http.StatusBadRequest, "Invalid request body", "Bad Request", errs)
}

func NewErrorEmptyRequestBody() *Response {
	return NewErrorResponse(http.StatusBadRequest, "Empty request body", "Bad Request", nil)
}
