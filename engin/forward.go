package engin

import (
	"net"
	"net/http"
)

type Forward interface {
	Close() error
	Http(r *http.Request) (*http.Response, error)
	Https(address string, r *http.Request) (net.Conn, error)
}
