// Copyright (c) 2025 EterLine (Andrew)
// This file is part of fstmon.
// Licensed under the MIT License. See the LICENSE file for details.

package middleware

import (
	"context"
	"crypto/subtle"
	"net"
	"net/http"
	"regexp"
	"time"

	"github.com/eterline/fstmon/internal/domain"
	"github.com/eterline/fstmon/internal/log"
	"github.com/eterline/fstmon/internal/services/secure"
	"github.com/eterline/fstmon/internal/web/controller"
)

func RequestWrapper(ipExt domain.IpExtractor) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			info := domain.InitRequestInfo(r, ipExt)

			w.Header().Set("X-Request-Time", time.Now().Format(time.RFC1123))

			r = r.WithContext(info.ToContext(r.Context()))
			next.ServeHTTP(w, r)
		})
	}
}

func RequestLogger(ctx context.Context) func(next http.Handler) http.Handler {
	log := log.MustLoggerFromContext(ctx)

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			log.DebugContext(
				r.Context(), "api request",
				"path", r.RequestURI,
				"method", r.Method,
			)

			next.ServeHTTP(w, r)
		})
	}
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

func SourceSubnetsAllow(ctx context.Context, ipExt domain.IpExtractor, cidr []string) func(http.Handler) http.Handler {
	log := log.MustLoggerFromContext(ctx)

	filter, err := secure.NewSubnetFilter(cidr)
	if err != nil {
		log.Warn("allowed subnets error", "error", err)
	}

	if s := filter.AllowedList(); len(s) != 0 {
		log.Warn("setup allowed subnets", "subnets", filter.AllowedList())
	}

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			ip, err := ipExt.ExtractIP(r)
			if err != nil {
				log.ErrorContext(r.Context(), "failed to parse request ip", "error", err.Error())
				controller.ResponseInternalError(w)
				return
			}

			if filter.InAllowedSubnets(ip) {
				log.DebugContext(r.Context(), "request ip allowed", "ip", ip.String())
				next.ServeHTTP(w, r)
				return
			}

			log.WarnContext(r.Context(), "request ip blocked", "ip", ip.String())
			controller.ResponseError(w, http.StatusForbidden, "forbidden: IP not allowed")
		})
	}
}

func AllowedHosts(ctx context.Context, hosts []string) func(http.Handler) http.Handler {
	log := log.MustLoggerFromContext(ctx)

	filter := secure.InitAllowedHostsFilter(hosts...)
	if filter != nil {
		log.Warn("setup allowed hosts", "hosts", filter.AllowedHosts())
	}

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			host, _, _ := net.SplitHostPort(r.Host)
			if !filter.InAllowedHosts(host) {
				log.WarnContext(r.Context(),
					"invalid request host",
					"request_host",
					host,
				)
				controller.ResponseError(w, http.StatusForbidden, "forbidden: invalid host")
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

func BearerCheck(ctx context.Context, bearer string, ipExt domain.IpExtractor) func(http.Handler) http.Handler {
	log := log.MustLoggerFromContext(ctx)
	enableAuth := !(bearer == "")
	log.Info("setup token auth policy", "auth_enabled", enableAuth)

	if !enableAuth {
		return func(next http.Handler) http.Handler {
			return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				next.ServeHTTP(w, r)
			})
		}
	}

	expected := []byte(bearer)
	bearerReg := regexp.MustCompile(`^Bearer:\s*(.+)`)

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			authHeader := r.Header.Get("Authorization")

			if matches := bearerReg.FindStringSubmatch(authHeader); matches != nil {
				token := []byte(matches[1])
				if subtle.ConstantTimeCompare(token, expected) == 1 {
					next.ServeHTTP(w, r)
					return
				}
			}

			log.WarnContext(
				r.Context(),
				"invalid request token",
				"auth_header", authHeader,
			)

			controller.ResponseError(w, http.StatusForbidden, "invalid token")
		})
	}
}
