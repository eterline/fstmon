package log

import (
	"context"
	"log/slog"
	"os"
	"reflect"

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
		t := reflect.TypeOf(c)
		v := reflect.ValueOf(c)

		for i := 0; i < t.NumField(); i++ {
			field := t.Field(i)
			tag := field.Tag.Get("json")
			if tag == "-" {
				continue
			}
			if tag == "" {
				tag = field.Name
			}
			value := v.Field(i).Interface()

			rec.Add(tag, value)
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
