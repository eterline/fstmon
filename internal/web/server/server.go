// Copyright (c) 2025 EterLine (Andrew)
// This file is part of fstmon.
// Licensed under the MIT License. See the LICENSE file for details.

package server

import (
	"crypto/tls"
	"fmt"
	"net/http"
)

type SwitchProtoMap map[string]func(*http.Server, *tls.Conn, http.Handler)

type Server struct {
	srv *http.Server
}

var (
	tlsCiphers = []uint16{
		tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
		tls.TLS_ECDHE_RSA_WITH_AES_256_CBC_SHA,
		tls.TLS_RSA_WITH_AES_256_GCM_SHA384,
		tls.TLS_RSA_WITH_AES_256_CBC_SHA,
	}
	tlsCurve = []tls.CurveID{
		tls.CurveP521,
		tls.CurveP384,
		tls.CurveP256,
	}
)

func NewServer(mux http.Handler) *Server {

	tlsConf := &tls.Config{
		PreferServerCipherSuites: true,
		CurvePreferences:         tlsCurve,
		CipherSuites:             tlsCiphers,
		MinVersion:               tls.VersionTLS12,
	}

	tlsProtoMap := make(map[string]func(*http.Server, *tls.Conn, http.Handler), 0)

	srv := &http.Server{
		Handler:      mux,
		TLSConfig:    tlsConf,
		TLSNextProto: tlsProtoMap,
	}

	return &Server{
		srv: srv,
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
