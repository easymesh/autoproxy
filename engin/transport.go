package engin

import (
	"crypto/tls"
	"net"
	"net/http"
	"time"
)

func newTransport(timeout int, tlscfg *tls.Config) *http.Transport {
	tmout := time.Duration(timeout) * time.Second
	return &http.Transport{
		TLSClientConfig: tlscfg,
		DialContext: (&net.Dialer{
			Timeout:   tmout,
			KeepAlive: tmout,
			DualStack: true,
		}).DialContext,
		ForceAttemptHTTP2:     true,
		MaxIdleConns:          1000,
		IdleConnTimeout:       2 * tmout,
		TLSHandshakeTimeout:   tmout,
		ExpectContinueTimeout: 5 * time.Second}
}
