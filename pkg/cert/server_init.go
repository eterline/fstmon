package cert

import (
	"crypto/tls"
	"net/http"
)

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

func NewServerTlsConfig(pref bool) *tls.Config {
	return &tls.Config{
		PreferServerCipherSuites: true,
		CurvePreferences:         tlsCurve,
		CipherSuites:             tlsCiphers,
		MinVersion:               tls.VersionTLS12,
		MaxVersion:               tls.VersionTLS13,
	}
}

func NewProtoMap() map[string]func(*http.Server, *tls.Conn, http.Handler) {
	return make(map[string]func(*http.Server, *tls.Conn, http.Handler), 0)
}
