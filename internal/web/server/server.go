// Copyright (c) 2025 EterLine (Andrew)
// This file is part of fstmon.
// Licensed under the MIT License. See the LICENSE file for details.

package server

import (
	"fmt"
	"net/http"

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

func (s *Server) Run(addr, key, crt string) error {
	s.srv.Addr = addr
	var err error = nil

	if key == "" || crt == "" {
		err = s.srv.ListenAndServe()
	} else {
		err = s.srv.ListenAndServeTLS(crt, key)
	}

	switch {
	case err == http.ErrServerClosed:
		return nil
	case err == nil:
		return nil
	default:
		return fmt.Errorf("server error: %w", err)
	}
}

func (s *Server) Close() error {
	return s.srv.Close()
}
