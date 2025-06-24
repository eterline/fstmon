package controller

import (
	"encoding/json"
	"net/http"

	"github.com/eterline/fstmon/internal/domain"
)

func writeJSON(w http.ResponseWriter, code int, resp any) error {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(code)
	return json.NewEncoder(w).Encode(resp)
}

func ResponseOK[T any](w http.ResponseWriter, data T) error {
	resp := domain.ResponseAPI[T]{
		Code: http.StatusOK,
		Data: data,
	}
	return writeJSON(w, http.StatusOK, resp)
}

func ResponseMessage(w http.ResponseWriter, message string) error {
	resp := domain.ResponseAPI[any]{
		Code:    http.StatusOK,
		Message: message,
	}
	return writeJSON(w, http.StatusOK, resp)
}

func ResponseError(w http.ResponseWriter, code int, message string) error {
	resp := domain.ResponseAPI[any]{
		Code:    code,
		Message: message,
	}
	return writeJSON(w, code, resp)
}

func ResponseInternalError(w http.ResponseWriter, err error) error {
	return ResponseError(w, http.StatusInternalServerError, err.Error())
}
