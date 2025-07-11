package log

import (
	"context"
	"log/slog"
	"os"
	"time"

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

		rec.Add("request_time", c.RequestTime.Format(time.RFC3339))

		if c.SourceIP != "" {
			rec.Add("source_ip", c.SourceIP)
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

func InitLogger(debug bool) {
	var h slog.Handler

	if debug {
		h = slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug})
	} else {
		h = slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo})
	}

	h = NewLoggerWithReqInfoWrapping(h)

	logger := slog.New(h)
	slog.SetDefault(logger)
}
