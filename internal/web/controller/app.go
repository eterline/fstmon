package controller

import "net/http"

func NotFound(w http.ResponseWriter, r *http.Request) {
	ResponseError(w, http.StatusNotFound, "page did not exists")
}

func BadMethod(w http.ResponseWriter, r *http.Request) {
	ResponseError(w, http.StatusMethodNotAllowed, "unsupported method")
}
