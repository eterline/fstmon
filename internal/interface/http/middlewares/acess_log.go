// Copyright (c) 2025 EterLine (Andrew)
// This file is part of fstmon.
// Licensed under the MIT License. See the LICENSE file for details.
package middleware

import (
	"io"
	"log/slog"
)

type AccessLogWriter struct {
	logger *slog.Logger
}

func NewAccessLogWriter(w io.Writer) *AccessLogWriter {
	if w == nil {
		return nil
	}
	h := slog.NewTextHandler(w, &slog.HandlerOptions{Level: slog.LevelInfo})
	return &AccessLogWriter{slog.New(h)}
}

func (l *AccessLogWriter) Log(fields ...any) {
	if l != nil {
		l.logger.Info("access", fields...)
	}
}
