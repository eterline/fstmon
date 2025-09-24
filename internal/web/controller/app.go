// Copyright (c) 2025 EterLine (Andrew)
// This file is part of fstmon.
// Licensed under the MIT License. See the LICENSE file for details.

package controller

import "net/http"

func NotFound(w http.ResponseWriter, r *http.Request) {
	ResponseError(w, http.StatusNotFound, "page did not exists")
}

func BadMethod(w http.ResponseWriter, r *http.Request) {
	ResponseError(w, http.StatusMethodNotAllowed, "unsupported method")
}
