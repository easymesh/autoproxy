package engin

import (
	"net"
	"net/http"
	"sync"
	"time"
)

type Forward interface {
	Close() error
	Http(r *http.Request) (*http.Response, error)
	Https(address string, r *http.Request) (net.Conn, error)
}

type defaultForward struct {
	sync.WaitGroup

	tmout int
	stop  chan struct{}
	trans *http.Transport
}

func (d *defaultForward) Close() error {
	d.stop <- struct{}{}
	d.trans.CloseIdleConnections()
	d.Wait()
	return nil
}

func (d *defaultForward) Http(r *http.Request) (*http.Response, error) {
	return d.trans.RoundTrip(r)
}

func (d *defaultForward) Https(address string, r *http.Request) (net.Conn, error) {
	return net.DialTimeout("tcp", address, time.Second*time.Duration(d.tmout))
}

func NewDefault(timeout int) (Forward, error) {
	forward := &defaultForward{
		trans: newTransport(timeout, nil),
		tmout: timeout,
		stop:  make(chan struct{}, 1),
	}
	forward.Add(1)
	go func() {
		ticker := time.NewTicker(5 * time.Second)
		defer ticker.Stop()
		defer forward.Done()
		for {
			select {
			case <-ticker.C:
				forward.trans.CloseIdleConnections()
			case <-forward.stop:
				return
			}
		}
	}()
	return forward, nil
}
