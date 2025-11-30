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
