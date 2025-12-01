// Copyright (c) 2025 EterLine (Andrew)
// This file is part of fstmon.
// Licensed under the MIT License. See the LICENSE file for details.
package middleware

import (
	"io"
	"log/slog"
	"os"
)

type AccessLogWriter struct {
	logger *slog.Logger
}

func NewAccessLogWriter(w io.Writer) *AccessLogWriter {
	if w == nil {
		w = os.Stdout
	}
	h := slog.NewTextHandler(w, &slog.HandlerOptions{Level: slog.LevelInfo})
	return &AccessLogWriter{slog.New(h)}
}

func (l *AccessLogWriter) Log(fields ...any) {
	// slog.InfoContext здесь делает запись в один уровень, уровни не используются
	l.logger.Info("access", fields...)
}
