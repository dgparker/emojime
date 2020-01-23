package httpsrv

import (
	"crypto/tls"
	"net/http"
	"time"
)

// READTIMEOUT is the default read timeout set for the http server
const READTIMEOUT = 5 * time.Second

// WRITETIMEOUT is the default write timeout set for the http server
const WRITETIMEOUT = 10 * time.Second

// IDLETIMEOUT is the default idle timeout set for the http server
const IDLETIMEOUT = 120 * time.Second

// tlsConfig sets the default TLS config for the http server
var tlsConfig = &tls.Config{
	PreferServerCipherSuites: true,
	CurvePreferences: []tls.CurveID{
		tls.CurveP256,
		tls.X25519,
	},
	MinVersion: tls.VersionTLS12,
	CipherSuites: []uint16{
		tls.TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384,
		tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
		tls.TLS_ECDHE_ECDSA_WITH_CHACHA20_POLY1305,
		tls.TLS_ECDHE_RSA_WITH_CHACHA20_POLY1305,
		tls.TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256,
		tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,
	},
}

// New returns a new http server using the provided configuration
func New(readTimeout, writeTimeout, idleTimeout time.Duration, tlsCfg *tls.Config, handler http.Handler) *http.Server {
	return &http.Server{
		ReadTimeout:  readTimeout,
		WriteTimeout: writeTimeout,
		IdleTimeout:  idleTimeout,
		TLSConfig:    tlsCfg,
		Handler:      handler,
	}
}

// NewWithDefault creates a new http server using the default configuration defined in the package
func NewWithDefault(handler http.Handler) *http.Server {
	return New(READTIMEOUT, WRITETIMEOUT, IDLETIMEOUT, tlsConfig, handler)
}
