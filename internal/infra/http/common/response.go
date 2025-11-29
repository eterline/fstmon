package common

import (
	"encoding/json"
	"net/http"
)

/*
ResponseHttpWrapper - unified HTTP JSON response wrapper.

Holds status code, message, array of errors, and optional payload of type T.
Must only be used on transport level and never referenced by domain logic.
*/
type ResponseHttpWrapper[T any] struct {
	Code    int      `json:"code"`              // HTTP status or internal code
	Message string   `json:"message,omitempty"` // Optional descriptive message
	Errors  []string `json:"errors,omitempty"`  // Array of error messages
	Data    *T       `json:"data,omitempty"`    // Optional payload of type T
}

// initErrs - ensures Errors slice is initialized with at least startLen capacity.
func (r *ResponseHttpWrapper[T]) initErrs(startLen int) {
	if r.Errors == nil {
		r.Errors = make([]string, 0, startLen)
	}
}

// SetCode - sets the HTTP status code and returns the wrapper for chaining.
func (r *ResponseHttpWrapper[T]) SetCode(code int) *ResponseHttpWrapper[T] {
	r.Code = code
	return r
}

// SetMessage - sets the message and returns the wrapper for chaining.
func (r *ResponseHttpWrapper[T]) SetMessage(msg string) *ResponseHttpWrapper[T] {
	r.Message = msg
	return r
}

// WrapData - sets the payload and returns the wrapper for chaining.
func (r *ResponseHttpWrapper[T]) WrapData(data T) *ResponseHttpWrapper[T] {
	r.Data = &data
	return r
}

// AddError - adds one or more error values to the Errors slice.
func (r *ResponseHttpWrapper[T]) AddError(err ...error) *ResponseHttpWrapper[T] {
	r.initErrs(len(err))
	for _, e := range err {
		r.Errors = append(r.Errors, e.Error())
	}
	return r
}

// AddStringError - adds one or more string errors to the Errors slice.
func (r *ResponseHttpWrapper[T]) AddStringError(err ...string) *ResponseHttpWrapper[T] {
	r.initErrs(len(err))
	r.Errors = append(r.Errors, err...)
	return r
}

// ================= Builders ==================

// NewResponse - creates a new empty response wrapper.
func NewResponse[T any]() *ResponseHttpWrapper[T] {
	return &ResponseHttpWrapper[T]{}
}

// OkResponse - creates a 200 response with a message and no payload.
func OkResponse[T any](message string) *ResponseHttpWrapper[T] {
	return NewResponse[T]().
		SetCode(http.StatusOK).
		SetMessage(message)
}

// OkDataResponse - creates a 200 response with payload.
func OkDataResponse[T any](data T) *ResponseHttpWrapper[T] {
	return NewResponse[T]().
		SetCode(http.StatusOK).
		WrapData(data)
}

// OkMsgResponse - creates a 200 response with both message and payload.
func OkMsgResponse[T any](message string, data T) *ResponseHttpWrapper[T] {
	return NewResponse[T]().
		SetCode(http.StatusOK).
		SetMessage(message).
		WrapData(data)
}

// ErrorSimpleResponse - creates an error response with code and string errors.
func ErrorSimpleResponse[T any](code int, errsDetails ...string) *ResponseHttpWrapper[T] {
	return NewResponse[T]().
		SetCode(code).
		SetMessage("error").
		AddStringError(errsDetails...)
}

// ErrorDetailedResponse - creates an error response with code, message, and Go errors.
func ErrorDetailedResponse[T any](code int, message string, errs ...error) *ResponseHttpWrapper[T] {
	return NewResponse[T]().
		SetCode(code).
		SetMessage(message).
		AddError(errs...)
}

// NewCustomResponse - creates a fully configurable response with string error.
func NewCustomResponse[T any](code int, message, errMsg string, data T) *ResponseHttpWrapper[T] {
	return NewResponse[T]().
		SetCode(code).
		SetMessage(message).
		AddStringError(errMsg).
		WrapData(data)
}

// NewCustomDetailedResponse - creates a fully configurable response with Go error.
func NewCustomDetailedResponse[T any](code int, message string, err error, data T) *ResponseHttpWrapper[T] {
	return NewResponse[T]().
		SetCode(code).
		SetMessage(message).
		AddError(err).
		WrapData(data)
}

// Write - writes the response as JSON into http.ResponseWriter.
// Handles status code, headers, JSON encoding and fallback on encoding error.
func (r *ResponseHttpWrapper[T]) Write(w http.ResponseWriter) error {
	w.Header().Set("Content-Type", StringContentType(""))

	if r.Code == 0 {
		w.WriteHeader(http.StatusOK)
	} else {
		w.WriteHeader(r.Code)
	}

	err := json.NewEncoder(w).Encode(r)
	if err != nil {
		http.Error(
			w,
			`{"code":500,"error":"failed to encode response"}`,
			http.StatusInternalServerError,
		)
	}

	return err
}
