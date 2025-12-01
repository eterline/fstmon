// Copyright (c) 2025 EterLine (Andrew)
// This file is part of fstmon.
// Licensed under the MIT License. See the LICENSE file for details.

package log

import (
	"context"
	"log/slog"
	"os"

	"github.com/eterline/fstmon/internal/domain"
)

type MiddlewareWithReqInfoWrapping struct {
	next slog.Handler
}

func NewLoggerWithReqInfoWrapping(next slog.Handler) *MiddlewareWithReqInfoWrapping {
	return &MiddlewareWithReqInfoWrapping{next: next}
}

func (mrw *MiddlewareWithReqInfoWrapping) Enabled(ctx context.Context, rec slog.Level) bool {
	return mrw.next.Enabled(ctx, rec)
}

func (mrw *MiddlewareWithReqInfoWrapping) Handle(ctx context.Context, rec slog.Record) error {
	if c, ok := domain.RequestInfoFromContext(ctx); ok {
		rec.Add("request_duration_ms", c.RequestDuration().Milliseconds())
		rec.Add("client", c.Client)
		rec.Add("source_ip", c.Source.Addr())
		rec.Add("source_port", c.Source.Port())
	}

	return mrw.next.Handle(ctx, rec)
}

func (mrw *MiddlewareWithReqInfoWrapping) WithAttrs(attrs []slog.Attr) slog.Handler {
	return &MiddlewareWithReqInfoWrapping{next: mrw.next.WithAttrs(attrs)}
}

func (mrw *MiddlewareWithReqInfoWrapping) WithGroup(name string) slog.Handler {
	return &MiddlewareWithReqInfoWrapping{next: mrw.next.WithGroup(name)}
}

// InitLogger – create singletone style structure logger
func InitLogger(level string, json bool) {
	logger := NewLogger(level, json)
	slog.SetDefault(logger)
}

// InitLogger – create new structure logger
func NewLogger(level string, json bool) *slog.Logger {
	opt := &slog.HandlerOptions{
		Level: selectLogLevel(level),
	}

	var h slog.Handler
	if json {
		h = slog.NewJSONHandler(os.Stdout, opt)
	} else {
		h = slog.NewTextHandler(os.Stdout, opt)
	}

	h = NewLoggerWithReqInfoWrapping(h)

	return slog.New(h)
}

// Uses only for context
type loggerContextKey int

const (
	loggerCtxKey loggerContextKey = iota
)

// WrapLoggerToContext – wrap logger to parent context
func WrapLoggerToContext(ctx context.Context, logger *slog.Logger) context.Context {
	return context.WithValue(ctx, loggerCtxKey, logger)
}

// LoggerFromContext – unwrap logger from context
func LoggerFromContext(ctx context.Context) (logger *slog.Logger, ok bool) {
	l, ok := ctx.Value(loggerCtxKey).(*slog.Logger)
	if ok {
		return l, true
	}
	return nil, false
}

// MustLoggerFromContext – unwrap logger from context
func MustLoggerFromContext(ctx context.Context) (logger *slog.Logger) {
	l, ok := ctx.Value(loggerCtxKey).(*slog.Logger)
	if ok {
		return l
	}
	panic("logger doesn't wrapped in this context.Context")
}

func selectLogLevel(level string) slog.Level {
	switch level {
	case "debug":
		return slog.LevelDebug
	case "info":
		return slog.LevelInfo
	case "warn":
		return slog.LevelWarn
	case "error":
		return slog.LevelError
	default:
		return slog.LevelInfo
	}
}
