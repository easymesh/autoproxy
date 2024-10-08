package engin

import (
	"context"
	"crypto/tls"
	"encoding/base64"
	"fmt"
	"io"
	"net"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/astaxie/beego/logs"
)

type HttpAccess struct {
	Timeout    int
	Address    string
	httpserver *http.Server
	sync.WaitGroup

	authHandler          func(auth *AuthInfo) bool
	forwardHandler       func(address string, r *http.Request) Forward
	forwardUpdateHandler func(address string, forward Forward)
	defaultForward       Forward
}

type Access interface {
	Shutdown() error
	AuthHandlerSet(func(*AuthInfo) bool)
	ForwardHandlerSet(func(address string, r *http.Request) Forward)
	ForwardUpdateHandlerSet(func(address string, forward Forward))
}

func HttpError(w http.ResponseWriter, err string, code int) {
	time.Sleep(time.Second)
	http.Error(w, err, code)
}

func AuthFailHandler(w http.ResponseWriter, r *http.Request) {
	logs.Warn("Request authentication failed. RemoteAddr: ", r.RemoteAddr)
	w.Header().Add("Proxy-Authenticate", "Basic realm=\"Access to internal site\"")
	HttpError(w, "Request authentication failed.", http.StatusProxyAuthRequired)
}

func AuthInfoParse(r *http.Request) *AuthInfo {
	value := r.Header.Get("Proxy-Authorization")
	if value == "" {
		return nil
	}
	body, err := base64.StdEncoding.DecodeString(value[6:])
	if err != nil {
		return nil
	}
	ctx := strings.Split(string(body), ":")
	if len(ctx) != 2 {
		return nil
	}
	return &AuthInfo{User: ctx[0], Token: ctx[1]}
}

func (acc *HttpAccess) NoProxyHandler(w http.ResponseWriter, r *http.Request) {
	logs.Warn("request is illegal. RemoteAddr: ", r.RemoteAddr)
	HttpError(w,
		"This is a proxy server. Does not respond to non-proxy requests.",
		http.StatusInternalServerError)
}

func (acc *HttpAccess) AuthHandlerSet(handler func(auth *AuthInfo) bool) {
	acc.authHandler = handler
}

func (acc *HttpAccess) ForwardHandlerSet(handler func(address string, r *http.Request) Forward) {
	acc.forwardHandler = handler
}

func (acc *HttpAccess) ForwardUpdateHandlerSet(handler func(address string, forward Forward)) {
	acc.forwardUpdateHandler = handler
}

func (acc *HttpAccess) AuthHttp(r *http.Request) bool {
	if acc.authHandler == nil {
		return true
	}
	if AuthCache(r) {
		return true
	}
	auth := acc.authHandler(AuthInfoParse(r))
	if auth {
		AuthLogin(r)
	}
	return auth
}

func (acc *HttpAccess) Shutdown() error {
	context, cencel := context.WithTimeout(context.Background(), 30*time.Second)
	err := acc.httpserver.Shutdown(context)
	cencel()
	if err != nil {
		logs.Error("http access ready to shut down fail, %s", err.Error())
	}
	acc.Wait()
	return err
}

func DebugReqeust(r *http.Request) {
	var headers string
	for key, value := range r.Header {
		headers += fmt.Sprintf("[%s:%s]", key, value)
	}
	logs.Info("RemoteAddr:%s Host:%s Scheme:%s Method:%s URL:%s Header:%s", r.RemoteAddr, r.Host, r.URL.Scheme, r.Method, r.URL.String(), headers)
}

func (acc *HttpAccess) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	DebugReqeust(r)

	if r.Header.Get("AUTOPROXY") == GetUUID() {
		HttpError(w, "loop request", http.StatusBadRequest)
		return
	}

	if !acc.AuthHttp(r) {
		AuthFailHandler(w, r)
		return
	}

	r.Header.Add("AUTOPROXY", GetUUID())

	if r.Method == "CONNECT" {
		acc.HttpsRoundTripper(w, r)
		return
	}

	StatUpdate(1, 0)
	defer StatUpdate(-1, 0)

	var rsp *http.Response
	var err error

	if !r.URL.IsAbs() {
		logs.Warn("the request is not proxy request, transport to local network")
		r.URL.Host = r.Host
		r.URL.Scheme = "http"
		rsp, err = acc.defaultForward.Http(r)
	} else {
		removeProxyHeaders(r)
		rsp, err = acc.HttpRoundTripper(r)
	}

	if err != nil {
		errStr := fmt.Sprintf("transport %s %s failed! %s", r.Host, r.URL.String(), err.Error())
		logs.Warn(errStr)
		HttpError(w, errStr, http.StatusInternalServerError)
		return
	}

	if rsp == nil {
		errStr := fmt.Sprintf("transport %s read response failed!", r.URL.Host)
		logs.Warn(errStr)
		HttpError(w, errStr, http.StatusInternalServerError)
		return
	}

	copyHeaders(w.Header(), rsp.Header)
	w.WriteHeader(rsp.StatusCode)

	size, err := io.Copy(w, rsp.Body)
	if size == 0 && err != nil {
		if err != io.EOF {
			logs.Info("io copy fail", err.Error())
		}
	} else {
		StatUpdate(0, size)
	}
	rsp.Body.Close()
}

func copyHeaders(dst, src http.Header) {
	for k, vs := range src {
		for _, v := range vs {
			dst.Add(k, v)
		}
	}
}

func removeProxyHeaders(r *http.Request) {
	r.RequestURI = ""
	r.Header.Del("Proxy-Connection")
	r.Header.Del("Proxy-Authenticate")
	r.Header.Del("Proxy-Authorization")
}

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
	address := Address(r.URL)
	server, err := acc.HttpsForward(address, r)
	if err != nil {
		errstr := fmt.Sprintf("can't forward hostname %s", address)
		logs.Error(errstr, err.Error())
		HttpError(w, errstr, http.StatusInternalServerError)
		return
	}

	hij, ok := w.(http.Hijacker)
	if !ok {
		logs.Error("httpserver does not support hijacking")
		panic("golang sdk is too old.")
	}

	client, _, err := hij.Hijack()
	if err != nil {
		logs.Error("Cannot hijack connection", err.Error())
		panic("golang sdk is too old.")
	}

	err = WriteFull(client, []byte("HTTP/1.1 200 Connection Established\r\n\r\n"))
	if err != nil {
		errstr := fmt.Sprintf("client connect %s fail", client.RemoteAddr())
		logs.Error(errstr, err.Error())
		client.Close()
		return
	}

	go Connect(acc, client, server)
}

func (acc *HttpAccess) HttpRoundTripper(r *http.Request) (*http.Response, error) {
	if r.Body != nil {
		r.Body = io.NopCloser(r.Body)
	}
	return acc.HttpForward(Address(r.URL), r)
}

func NewHttpsAccess(addr string, timeout int, tlsEnable bool, certfile, keyfile string) (Access, error) {
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		logs.Error("listen address fail", addr)
		return nil, err
	}

	var config *tls.Config
	if tlsEnable {
		config, err = TlsConfigServer(certfile, keyfile)
		if err != nil {
			logs.Error("make tls config server fail, %s", err.Error())
			return nil, err
		}
		logs.Info("enable TLS for %s", addr)

		lis = tls.NewListener(lis, config)
	}

	acc := new(HttpAccess)
	acc.Address = addr
	acc.Timeout = timeout
	acc.defaultForward = NewLocalForward(timeout)

	tmout := time.Duration(timeout) * time.Second

	httpserver := &http.Server{
		Handler:      acc,
		ReadTimeout:  tmout,
		WriteTimeout: tmout,
		TLSConfig:    config,
	}

	acc.httpserver = httpserver

	acc.Add(1)

	go func() {
		defer acc.Done()
		err = httpserver.Serve(lis)
		if err != nil {
			logs.Error("http server ", err.Error())
		}
	}()

	schema := "http"
	if config != nil {
		schema = "https"
	}

	logs.Info("listen local %s://%s as service with timeout %d seconds success.", schema, lis.Addr().String(), timeout)
	return acc, nil
}
