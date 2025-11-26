// Copyright (c) 2025 EterLine (Andrew)
// This file is part of fstmon.
// Licensed under the MIT License. See the LICENSE file for details.

package server

import (
	"context"
	"errors"
	"net/http"
	"time"

	"github.com/eterline/fstmon/internal/log"
	"github.com/eterline/fstmon/pkg/cert"
)

type Server struct {
	srv *http.Server
}

func NewServer(mux http.Handler) *Server {
	return &Server{
		srv: &http.Server{
			Handler:      mux,
			TLSConfig:    cert.NewServerTlsConfig(true),
			TLSNextProto: cert.NewProtoMap(),
		},
	}
}

func (s *Server) Run(ctx context.Context, addr, key, crt string) error {
	logger := log.MustLoggerFromContext(ctx)
	logger = logger.With("listen", addr)

	s.srv.Addr = addr

	tlsEnabled := !(key == "" || crt == "")
	logger = logger.With("tls", tlsEnabled)

	errCh := make(chan error, 1)

	go func() {
		var err error

		if tlsEnabled {
			logger.Info("starting server", "crt", crt, "key", key)
			err = s.srv.ListenAndServeTLS(crt, key)
		} else {
			logger.Info("starting server")
			err = s.srv.ListenAndServe()
		}

		errCh <- err
	}()

	// --- WAIT FOR STOP ---
	select {
	case <-ctx.Done():
		logger.Info("shutdown signal received")

		shCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		if err := s.srv.Shutdown(shCtx); err != nil {
			logger.Error("shutdown error", "error", err)
			return err
		}

		logger.Info("server stopped gracefully")
		return nil

	case err := <-errCh:
		if errors.Is(err, http.ErrServerClosed) {
			logger.Info("server closed normally")
			return nil
		}

		logger.Error("server error", "error", err)
		return err
	}
}

func (s *Server) Close() error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	return s.srv.Shutdown(ctx)
}
