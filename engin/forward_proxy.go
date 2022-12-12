package engin

import (
	"bytes"
	"crypto/tls"
	"encoding/base64"
	"fmt"
	"net"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/astaxie/beego/logs"
	"golang.org/x/net/http/httpproxy"
)

type httpsProtocal struct {
	auth      *AuthInfo
	config    *tls.Config
	address   string
	timeout   time.Duration
	proxycfg  *httpproxy.Config
	proxyfunc func(reqURL *url.URL) (*url.URL, error)
	trans     *http.Transport
}

func (h *httpsProtocal) ProxyFunc(r *http.Request) (*url.URL, error) {
	return h.proxyfunc(r.URL)
}

func headerCoder(values []string) string {
	body := bytes.NewBuffer(make([]byte, 0))
	for i, v := range values {
		if i == len(values)-1 {
			body.WriteString(v)
		} else {
			body.WriteString(v + " ")
		}
	}
	return body.String()
}

func httpsProxyRequest(r *http.Request) []byte {
	body := bytes.NewBuffer(make([]byte, 0))
	body.WriteString(fmt.Sprintf("CONNECT %s %s\r\n", r.URL.Host, r.Proto))
	body.WriteString(fmt.Sprintf("Host: %s\r\n", r.Host))
	for key, value := range r.Header {
		body.WriteString(fmt.Sprintf("%s: %s\r\n", key, headerCoder(value)))
	}
	body.WriteString("\r\n")
	return body.Bytes()
}

func httpsProxyAuthAdd(r *http.Request, auth *AuthInfo) {
	if auth == nil {
		return
	}
	authBody := auth.User + ":" + auth.Token
	basic := "Basic " + base64.StdEncoding.EncodeToString([]byte(authBody))
	r.Header.Add("Proxy-Authorization", basic)
}

func (h *httpsProtocal) Http(r *http.Request) (*http.Response, error) {
	rsp, err := h.trans.RoundTrip(r)
	if err != nil {
		errStr := fmt.Sprintf("http roundtrip %s %s fail!", r.Host, r.RemoteAddr)
		logs.Warn(errStr, err.Error())
	}
	return rsp, err
}

func (h *httpsProtocal) Https(address string, r *http.Request) (net.Conn, error) {
	server, err := net.DialTimeout("tcp", h.address, h.timeout)
	if err != nil {
		return nil, fmt.Errorf("connect to proxy %s failed, err=%s",
			h.address, err.Error())
	}

	if h.config != nil {
		server = tls.Client(server, h.config)
	}

	r.Header.Del("Proxy-Authenticate")
	httpsProxyAuthAdd(r, h.auth)

	err = WriteFull(server, httpsProxyRequest(r))
	if err != nil {
		server.Close()
		return nil, fmt.Errorf("write to proxy failed! %s", err.Error())
	}

	var readbuf [1024]byte
	cnt, err := server.Read(readbuf[:])
	if err != nil {
		server.Close()
		return nil, fmt.Errorf("read from remote proxy failed! %s", err.Error())
	}

	if -1 == strings.Index(string(readbuf[:cnt]), "200") {
		server.Close()
		return nil, fmt.Errorf("read from remote proxy fail, rsponse[%s]", string(readbuf[:cnt]))
	}

	return server, nil
}

func (d *httpsProtocal) Close() error {
	d.trans.CloseIdleConnections()
	return nil
}

func NewHttpsProtcal(address string, timeout int, auth *AuthInfo, tlsEnable bool, certfile, keyfile string) (Forward, error) {
	var config *tls.Config

	if tlsEnable {
		var err error
		config, err = TlsConfigClient(address, certfile, keyfile)
		if err != nil {
			logs.Error("make tls config client fail, %s", err.Error())
			return nil, err
		}
	}

	h := new(httpsProtocal)
	h.address = address
	h.config = config
	h.timeout = time.Second * time.Duration(timeout)

	scheme := "http"
	if config != nil {
		scheme = "https"
	}
	h.auth = auth

	var proxy string
	if auth == nil {
		proxy = fmt.Sprintf("%s://%s", scheme, address)
	} else {
		proxy = fmt.Sprintf("%s://%s:%s@%s", scheme,
			url.QueryEscape(auth.User), url.QueryEscape(auth.Token), address)
	}

	logs.Info("proxy %s init success", proxy)

	h.proxycfg = &httpproxy.Config{HTTPProxy: proxy, HTTPSProxy: proxy}
	h.proxyfunc = h.proxycfg.ProxyFunc()

	h.trans = newTransport(timeout, config)
	h.trans.Proxy = h.ProxyFunc

	return h, nil
}
