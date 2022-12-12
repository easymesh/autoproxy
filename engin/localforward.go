package engin

import (
	"net"
	"net/http"
	"sync"
	"time"
)

type LocalForward struct {
	sync.WaitGroup

	tmout int
	stop  chan struct{}
	trans *http.Transport
}

func (d *LocalForward) Close() error {
	d.stop <- struct{}{}
	d.trans.CloseIdleConnections()
	d.Wait()
	return nil
}

func (d *LocalForward) Http(r *http.Request) (*http.Response, error) {
	return d.trans.RoundTrip(r)
}

func (d *LocalForward) Https(address string, r *http.Request) (net.Conn, error) {
	return net.DialTimeout("tcp", address, time.Second*time.Duration(d.tmout))
}

func NewLocalForward(timeout int) Forward {
	forward := &LocalForward{
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
	return forward
}
