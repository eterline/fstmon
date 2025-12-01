// Copyright (c) 2025 EterLine (Andrew)
// This file is part of fstmon.
// Licensed under the MIT License. See the LICENSE file for details.
package middleware

import (
	"context"
	"io"
	"net/http"
	"net/netip"
	"time"

	"github.com/eterline/fstmon/internal/infra/http/common/api"
	"github.com/eterline/fstmon/internal/log"
)

// ================================================================================

type CtxKeyRoot int

const (
	RootKey CtxKeyRoot = iota
)

type RootRequestContext struct {
	StartedAt time.Time

	// parsed request parameters
	Method       string
	Path         string
	Host         string
	UserAgent    string
	EncodingList []string

	ClientIP netip.Addr
	SourceIP netip.AddrPort

	// ResponseWriter wrapper
	RW *ResponseRecorder

	// Access logger
	AccessLogger *AccessLogWriter
}

func NewRootRequestContext() *RootRequestContext {
	return &RootRequestContext{
		StartedAt: time.Now(),
	}
}

func (c *RootRequestContext) Duration() time.Duration {
	return time.Since(c.StartedAt)
}

func CtxRoot(ctx context.Context) *RootRequestContext {
	v, ok := ctx.Value(RootKey).(*RootRequestContext)
	if !ok {
		panic("request root context is missing")
	}
	return v
}

func WithRootContext(ctx context.Context, rctx *RootRequestContext) context.Context {
	return context.WithValue(ctx, RootKey, rctx)
}

// ================================================================================

type ResponseRecorder struct {
	http.ResponseWriter
	Status int
}

func NewRecorder(w http.ResponseWriter) *ResponseRecorder {
	return &ResponseRecorder{
		ResponseWriter: w,
		Status:         http.StatusOK,
	}
}

func (rec *ResponseRecorder) WriteHeader(status int) {
	rec.Status = status
	rec.ResponseWriter.WriteHeader(status)
}

// ================================================================================

func RootMiddleware(ctx context.Context, ipExt IpExtractor, accessWriter io.Writer, logAccess bool) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			rootLog := log.MustLoggerFromContext(ctx)

			root := NewRootRequestContext()

			// Request parsing
			root.Method = r.Method
			root.Path = r.URL.Path
			root.Host = r.Host
			root.UserAgent = r.UserAgent()
			root.EncodingList = r.TransferEncoding

			// IP extraction
			cleintAddr, srcAddr, err := ipExt.ExtractIP(r)
			if err != nil {
				rootLog.Error("extarting request IP failed", "error", err)
				api.InternalErrorResponse().
					SetMessage("internal servser error").
					Write(w)

				return
			}

			root.ClientIP = cleintAddr
			root.SourceIP = srcAddr

			// Wrap Response
			rec := NewRecorder(w)
			root.RW = rec

			// Prepare access logger
			root.AccessLogger = NewAccessLogWriter(accessWriter)

			// Inject context
			ctx := WithRootContext(r.Context(), root)

			next.ServeHTTP(rec, r.WithContext(ctx))

			if logAccess {
				root.AccessLogger.Log(
					"method", root.Method,
					"path", root.Path,
					"status_code", rec.Status,
					"duration_ms", root.Duration().Milliseconds(),
					"client_ip", root.ClientIP.String(),
					"user_agent", root.UserAgent,
					"source_ip", root.SourceIP.Addr().String(),
					"source_port", root.SourceIP.Port(),
					"encodings", root.EncodingList,
				)
			}
		})
	}
}
