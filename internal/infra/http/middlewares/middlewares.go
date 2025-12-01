// Copyright (c) 2025 EterLine (Andrew)
// This file is part of fstmon.
// Licensed under the MIT License. See the LICENSE file for details.

package middleware

import (
	"context"
	"log/slog"
	"net"
	"net/http"
	"time"

	"github.com/eterline/fstmon/internal/domain"
	"github.com/eterline/fstmon/internal/infra/http/common/api"
	"github.com/eterline/fstmon/internal/infra/http/common/security"
	"github.com/eterline/fstmon/internal/log"
)

func RequestLoggerWrap(logger *slog.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := log.WrapLoggerToContext(r.Context(), logger)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// RequestWrapper – wraps incoming HTTP requests with additional metadata.
// Initializes RequestInfo, extracts client IP, attaches it to the request context,
// and sets the X-Request-Time header.
func RequestWrapper(ipExt IpExtractor) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			info := domain.InitRequestInfo(r, ipExt)

			w.Header().Set("X-Request-Time", info.RequestCreated().Format(time.RFC1123))

			r = r.WithContext(info.ToContext(r.Context()))
			next.ServeHTTP(w, r)
		})
	}
}

// TODO: Make access log stream with another walues
// RequestLogger – logs basic information about each HTTP request, such as path and method,
// using the logger stored in the provided context.
func RequestLogger(ctx context.Context) func(next http.Handler) http.Handler {
	log := log.MustLoggerFromContext(ctx)

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			next.ServeHTTP(w, r)

			log.DebugContext(
				r.Context(), "api request",
				"path", r.RequestURI,
				"method", r.Method,
			)
		})
	}
}

// NoCacheControl – disables browser caching by setting strict no-cache headers.
func NoCacheControl(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		w.Header().Set("Cache-Control", "no-store, no-cache, must-revalidate, max-age=0")
		w.Header().Set("Pragma", "no-cache")
		w.Header().Set("Expires", "0")

		next.ServeHTTP(w, r)
	})
}

// SecureControl – applies security-related HTTP headers to harden responses.
// Adds protection against XSS, clickjacking, MIME sniffing, and enforces a strict CSP.
func SecureControl(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		// Protects against XSS in older browsers
		w.Header().Set("X-XSS-Protection", "1; mode=block")

		// Prevents embedding the page in an iframe (protects against clickjacking)
		w.Header().Set("X-Frame-Options", "DENY")

		// Prevents MIME sniffing so that the browser strictly adheres to Content-Type
		w.Header().Set("X-Content-Type-Options", "nosniff")

		// Content Security Policy – restricts allowed sources
		w.Header().Set("Content-Security-Policy", "default-src 'none';")

		next.ServeHTTP(w, r)
	})
}

// SourceSubnetsAllow – allows access only from specific CIDR subnets.
// Extracts request IP, checks it against the whitelist, and blocks requests from disallowed networks.
func SourceSubnetsAllow(ctx context.Context, ipExt IpExtractor, cidr []string) func(http.Handler) http.Handler {
	log := log.MustLoggerFromContext(ctx)

	filter, err := security.NewSubnetFilter(cidr)
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
				api.InternalErrorResponse().Write(w)
				return
			}

			if filter.InAllowedSubnets(ip) {
				log.DebugContext(r.Context(), "request ip allowed", "ip", ip.String())
				next.ServeHTTP(w, r)
				return
			}

			log.WarnContext(r.Context(), "request ip blocked", "ip", ip.String())

			if err := api.NewResponse().
				SetCode(http.StatusForbidden).
				SetMessage("request forbidden").
				AddStringError("ip address whitelist mismatch").
				Write(w); err != nil {
				log.Error("response error", "error", err)
			}
		})
	}
}

// AllowedHosts – checks the Host header against a list of allowed hosts and blocks unauthorized origins.
// Protects against host header attacks.
func AllowedHosts(ctx context.Context, hosts []string) func(http.Handler) http.Handler {
	log := log.MustLoggerFromContext(ctx)

	filter := security.InitAllowedHostsFilter(hosts...)
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

				if err := api.NewResponse().
					SetCode(http.StatusForbidden).
					SetMessage("request forbidden").
					AddStringError("invalid request host").
					Write(w); err != nil {
					log.Error("response error", "error", err)
				}

				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

// BearerAuth – validates Authorization header.
func BearerAuth(ctx context.Context, btest BearerTester) func(http.Handler) http.Handler {
	log := log.MustLoggerFromContext(ctx)

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			if b := r.Header.Get("Authorization"); btest.TestBearer(b) {
				next.ServeHTTP(w, r)
				return
			}

			log.WarnContext(r.Context(), "invalid auth bearer")

			if err := api.NewResponse().
				SetCode(http.StatusForbidden).
				SetMessage("request forbidden").
				AddStringError("invalid token bearer").
				Write(w); err != nil {
				log.Error("response error", "error", err)
			}
		})
	}
}
