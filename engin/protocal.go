package engin

import (
	"fmt"
	"io/ioutil"
	"net"
	"net/http"

	"github.com/astaxie/beego/logs"
)

func (acc *HttpAccess) ForwardUpdate(address string, forward Forward) {
	if acc.forwardUpdateHandler != nil {
		acc.forwardUpdateHandler(address, forward)
	}
}

func (acc *HttpAccess) HttpsForward(address string, r *http.Request) (net.Conn, error) {
	if acc.forwardHandler == nil {
		return nil, fmt.Errorf("forward handler is null")
	}
	forward := acc.forwardHandler(address, r)
	conn, err := forward.Https(address, r)
	if err != nil {
		acc.ForwardUpdate(address, forward)
	}
	return conn, err
}

func (acc *HttpAccess) HttpForward(address string, r *http.Request) (*http.Response, error) {
	if acc.forwardHandler == nil {
		return nil, fmt.Errorf("forward handler is null")
	}
	forward := acc.forwardHandler(address, r)
	conn, err := forward.Http(r)
	if err != nil {
		acc.ForwardUpdate(address, forward)
	}
	return conn, err
}

func (acc *HttpAccess) HttpsRoundTripper(w http.ResponseWriter, r *http.Request) {
	hij, ok := w.(http.Hijacker)
	if !ok {
		logs.Error("httpserver does not support hijacking")
	}

	client, _, err := hij.Hijack()
	if err != nil {
		logs.Error("Cannot hijack connection", err.Error())
		panic("golang sdk is too old.")
	}

	address := Address(r.URL)

	server, err := acc.HttpsForward(address, r)
	if err != nil {
		errstr := fmt.Sprintf("can't forward hostname %s", address)
		logs.Error(errstr, err.Error())
		HttpError(w, errstr, http.StatusInternalServerError)

		client.Close()
		return
	}

	connection := fmt.Sprintf("HTTP/1.1 200 Connection Established\r\n\r\n")

	err = WriteFull(client, []byte(connection))
	if err != nil {
		errstr := fmt.Sprintf("client connect %s fail", client.RemoteAddr())
		logs.Error(errstr, err.Error())
		HttpError(w, errstr, http.StatusInternalServerError)

		client.Close()
		return
	}

	go Connect(acc, client, server)
}

func (acc *HttpAccess) HttpRoundTripper(r *http.Request) (*http.Response, error) {
	if r.Body != nil {
		r.Body = ioutil.NopCloser(r.Body)
	}
	return acc.HttpForward(Address(r.URL), r)
}
