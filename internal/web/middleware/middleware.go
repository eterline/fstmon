package middleware

import (
	"crypto/subtle"
	"log/slog"
	"net"
	"net/http"
	"regexp"
	"time"

	"github.com/eterline/fstmon/internal/domain"
	"github.com/eterline/fstmon/internal/services/ipfilter"
	"github.com/eterline/fstmon/internal/web/controller"
)

func RequestWrapper(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		info := domain.InitRequestInfo(r)

		w.Header().Set("X-Request-Time", info.RequestTime.UTC().Format(time.RFC1123))
		w.Header().Set("X-Request-ID", info.RequestID.String())

		r = r.WithContext(info.ToContext(r.Context()))
		next.ServeHTTP(w, r)
	})
}

func RequestLogger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		slog.DebugContext(
			r.Context(), "api request",
			"path", r.RequestURI,
			"method", r.Method,
		)

		next.ServeHTTP(w, r)
	})
}

func NoCacheControl(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		w.Header().Set("Cache-Control", "no-store, no-cache, must-revalidate, max-age=0")
		w.Header().Set("Pragma", "no-cache")
		w.Header().Set("Expires", "0")

		next.ServeHTTP(w, r)
	})
}

func SecureControl(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		// Protects against XSS in older browsers
		w.Header().Set("X-XSS-Protection", "1; mode=block")

		// Prevents embedding the page in an iframe (protects against clickjacking)
		w.Header().Set("X-Frame-Options", "DENY")

		// Prevents MIME sniffing so that the browser strictly adheres to Content-Type
		w.Header().Set("X-Content-Type-Options", "nosniff")

		// Content Security Policy — you can restrict loading resources
		// For example, allow only your own scripts, styles, etc.
		// For APIs with JSON, you can usually use the basic policy
		w.Header().Set("Content-Security-Policy", "default-src 'none';")

		next.ServeHTTP(w, r)
	})
}

func SourceSubnetsAllow(cidr string) func(http.Handler) http.Handler {
	filter := ipfilter.InitIpFilter(cidr)

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ip, _, err := net.SplitHostPort(r.RemoteAddr)
			if err != nil {
				controller.ResponseError(w, http.StatusForbidden, "forbidden: IP not allowed")
				return
			}

			if !filter.InAllowedCIDR(ip) {
				controller.ResponseError(w, http.StatusForbidden, "forbidden: IP not allowed")
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

func BearerCheck(bearer string) func(http.Handler) http.Handler {

	if bearer == "" {
		return emptyToken
	}

	expected := []byte(bearer)
	bearerReg := regexp.MustCompile(`^Bearer:\s*(.+)`)

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			authHeader := r.Header.Get("Authorization")

			if matches := bearerReg.FindStringSubmatch(authHeader); matches != nil {
				token := matches[1]
				if subtle.ConstantTimeCompare([]byte(token), expected) == 1 {
					next.ServeHTTP(w, r)
					return
				}
			}

			slog.WarnContext(r.Context(),
				"invalid request token",
				"auth_header",
				authHeader,
			)

			controller.ResponseError(w, http.StatusForbidden, "invalid token")
		})
	}
}

func emptyToken(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		slog.InfoContext(r.Context(), "auth disabled")
		next.ServeHTTP(w, r)
	})
}
