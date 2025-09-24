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
		if c.SourceIP != "" {
			rec.Add("source_ip", c.SourceIP)
		}

		if c.ClientIP != "" {
			rec.Add("client_ip", c.ClientIP)
		}
	}

	return mrw.next.Handle(ctx, rec)
}

func (mrw *MiddlewareWithReqInfoWrapping) WithAttrs(attrs []slog.Attr) slog.Handler {
	return &MiddlewareWithReqInfoWrapping{next: mrw.next.WithAttrs(attrs)}
}

func (mrw *MiddlewareWithReqInfoWrapping) WithGroup(name string) slog.Handler {
	return &MiddlewareWithReqInfoWrapping{next: mrw.next.WithGroup(name)}
}

func InitLogger(debug, json bool) {
	var h slog.Handler
	opt := &slog.HandlerOptions{}

	if debug {
		opt.Level = slog.LevelDebug

		opt = &slog.HandlerOptions{Level: slog.LevelDebug}
		h = slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug})
	} else {
		opt.Level = slog.LevelInfo
		h = slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo})
	}

	if json {
		h = slog.NewJSONHandler(os.Stdout, opt)
	} else {
		h = slog.NewTextHandler(os.Stdout, opt)
	}

	h = NewLoggerWithReqInfoWrapping(h)

	logger := slog.New(h)
	slog.SetDefault(logger)
}
