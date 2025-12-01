package middleware

import (
	"context"
	"net/http"

	"github.com/eterline/fstmon/internal/infra/http/common/api"
	"github.com/eterline/fstmon/internal/infra/http/common/security"
)

func SecureHeaders(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		w.Header().Set("X-XSS-Protection", "1; mode=block")
		w.Header().Set("X-Frame-Options", "DENY")
		w.Header().Set("X-Content-Type-Options", "nosniff")
		w.Header().Set("Content-Security-Policy", "default-src 'none';")

		next.ServeHTTP(w, r)
	})
}

func SourceSubnetsAllow(ctx context.Context, allw SubnetAllower) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			root := CtxRoot(r.Context())

			if allw.InAllowedSubnets(root.ClientIP) {
				next.ServeHTTP(w, r)
				return
			}

			if err := api.NewResponse().
				SetCode(http.StatusForbidden).
				SetMessage("request forbidden").
				AddStringError("ip address not allowed").
				Write(w); err != nil {
				root.AccessLogger.Log("event", "response_error", "error", err.Error())
			}
		})
	}
}

func BearerAuth(btest BearerTester, headerName string) func(http.Handler) http.Handler {

	if headerName == "" {
		headerName = "Authorization"
	}

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			root := CtxRoot(r.Context())

			if btest.TestBearer(r.Header.Get(headerName)) {
				next.ServeHTTP(w, r)
				return
			}

			root.AccessLogger.Log(
				"event", "request_block",
				"reason", "invalid_token",
				"ip", root.ClientIP.String(),
				"path", root.Path,
			)

			if err := api.NewResponse().
				SetCode(http.StatusForbidden).
				SetMessage("invalid bearer").
				AddStringError("auth failed").
				Write(w); err != nil {
				root.AccessLogger.Log("event", "response_error", "error", err.Error())
			}
		})
	}
}

func AllowedHosts(hosts []string) func(http.Handler) http.Handler {

	filter := security.InitAllowedHostsFilter(hosts...)

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			root := CtxRoot(r.Context())

			if filter.InAllowedHosts(root.Host) {
				next.ServeHTTP(w, r)
				return
			}

			root.AccessLogger.Log(
				"event", "request_block",
				"reason", "invalid_host",
				"host", root.Host,
				"ip", root.ClientIP.String(),
				"path", root.Path,
			)

			if err := api.NewResponse().
				SetCode(http.StatusForbidden).
				SetMessage("request forbidden").
				AddStringError("invalid request host").
				Write(w); err != nil {
				root.AccessLogger.Log("event", "response_error", "error", err.Error())
			}

		})
	}
}
