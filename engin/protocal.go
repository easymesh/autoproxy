package engin

import (
	"bytes"
	"fmt"
	"github.com/astaxie/beego/logs"
	"io/ioutil"
	"net"
	"net/http"
	"sync/atomic"
)

var HTTPS_CLIENT_CONNECT_FLAG  = []byte("HTTP/1.1 200 Connection Established\r\n\r\n")

func (acc *HttpAccess)HttpsForward(address string, r *http.Request) (net.Conn, error) {
	if acc.forwardHandler != nil {
		forward := acc.forwardHandler(address, r)
		return forward.Https(address, r)
	}
	return nil, fmt.Errorf("forward handler is null")
}

func (acc *HttpAccess)HttpForward(address string, r *http.Request) (*http.Response, error) {
	if acc.forwardHandler != nil {
		forward := acc.forwardHandler(address, r)
		return forward.Http(r)
	}
	return nil, fmt.Errorf("forward handler is null")
}

func (acc *HttpAccess)HttpsRoundTripper(w http.ResponseWriter, r *http.Request) {
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

	err = WriteFull(client, HTTPS_CLIENT_CONNECT_FLAG)
	if err != nil {
		errstr := fmt.Sprintf("client connect %s fail", client.RemoteAddr())
		logs.Error(errstr, err.Error())
		http.Error(w, errstr, http.StatusInternalServerError)

		client.Close()
		return
	}

	server, err := acc.HttpsForward(address, r)
	if err != nil {
		errstr := fmt.Sprintf("can't forward hostname %s", address)
		logs.Error(errstr, err.Error())
		http.Error(w, errstr, http.StatusInternalServerError)

		client.Close()
		return
	}

	go func() {
		ConnectCopyWithTimeout(client, server, 60, func(cnt uint64) {
			atomic.AddUint64(&acc.flowsize, cnt)
		})
	}()
}

func (acc *HttpAccess)HttpRoundTripper(r *http.Request) (*http.Response, error) {
	var bodyBytes []byte
	if r.Body != nil {
		var err error
		bodyBytes, err = ioutil.ReadAll(r.Body)
		if err != nil {
			logs.Error("read all fail, %s", err.Error())
		}
		atomic.AddUint64(&acc.flowsize, uint64(len(bodyBytes)))
	}
	r.Body = ioutil.NopCloser(bytes.NewBuffer(bodyBytes))
	return acc.HttpForward(Address(r.URL), r)
}

