// Copyright (c) 2025 EterLine (Andrew)
// This file is part of fstmon.
// Licensed under the MIT License. See the LICENSE file for details.
package middleware

import (
	"net/http"
	"net/netip"
)

type IpExtractor interface {
	ExtractIP(*http.Request) (netip.Addr, error)
}

type BearerTester interface {
	TestBearer(string) bool
}
