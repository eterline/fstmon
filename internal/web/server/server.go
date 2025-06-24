package server

import (
	"crypto/tls"
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
	if key == "" || crt == "" {
		return s.srv.ListenAndServe()
	}
	return s.srv.ListenAndServeTLS(crt, key)
}

func (s *Server) Close() error {
	return s.srv.Close()
}
