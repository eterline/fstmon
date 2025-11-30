package api

import (
	"fmt"
	"net/http"
)

func HandleNotFound(w http.ResponseWriter, r *http.Request) {
	NewResponse().
		SetCode(http.StatusNotFound).
		SetMessage("handler not found").
		Write(w)
}

func HandleNotAllowedMethod(w http.ResponseWriter, r *http.Request) {
	NewResponse().
		SetCode(http.StatusMethodNotAllowed).
		SetMessage(fmt.Sprintf("%s: method not allowed", r.Method)).
		Write(w)
}
