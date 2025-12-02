package middleware

import "net"

// requestHostname returns the hostname without port from http.Request.Host.
// It handles both "host:port" and plain "host" formats safely.
func requestHostname(host string) string {
	// If host contains a port, split it.
	h, _, err := net.SplitHostPort(host)
	if err == nil {
		return h
	}

	// If no port exists, SplitHostPort returns error -> host already clean.
	return host
}
