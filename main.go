package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/astaxie/beego/logs"
	"github.com/easymesh/autoproxy/engin"
)

var (
	Help    bool
	Debug   bool
	Timeout int

	LocalAddr  string
	LocalAuth  string
	RemoteAddr string
	RemoteAuth string

	RunMode    string // local、proxy、domain、auto
	DomainFile string
	CertFile   string
	KeyFile    string
)

func init() {
	flag.StringVar(&KeyFile, "key-file", "", "tls key file pem format")
	flag.StringVar(&CertFile, "cert-file", "", "tls cert file pem format")

	flag.IntVar(&Timeout, "timeout", 30, "connect timeout (unit second)")

	flag.StringVar(&LocalAddr, "local-address", "http://0.0.0.0:8080", "Local proxy listening address")
	flag.StringVar(&LocalAuth, "local-auth", "", "Local proxy auth username and password")

	flag.StringVar(&RemoteAddr, "remote-address", "https://you.domain.com:8080", "Remote proxy listening address")
	flag.StringVar(&RemoteAuth, "remote-auth", "", "Remote proxy auth username and password")

	flag.StringVar(&RunMode, "mode", "proxy", "proxy mode(local/proxy/domain/auto)")
	flag.StringVar(&DomainFile, "domain", "domain.json", "match domain list file(domain mode requires)")

	flag.BoolVar(&Debug, "debug", false, "enable debug")
	flag.BoolVar(&Help, "help", false, "usage help")
}

func parseAuth(auth string) (*engin.AuthInfo, error) {
	if auth == "" {
		return nil, nil
	}
	list := strings.Split(auth, ":")
	if len(list) != 2 {
		return nil, fmt.Errorf("Authentication information '%s' is incorrect", auth)
	}
	return &engin.AuthInfo{User: list[0], Token: list[1]}, nil
}

func parseAddress(addr string) (string, string, error) {
	ul, err := url.Parse(addr)
	if err != nil {
		return "", "", err
	}
	scheme := strings.ToLower(ul.Scheme)
	host := ul.Host
	if -1 == strings.Index(host, ":") {
		if scheme == "https" {
			host += "443"
		} else {
			host += "80"
		}
	}
	return scheme, host, nil
}

func parseDomain(domain string) ([]string, error) {
	body, err := ioutil.ReadFile(domain)
	if err != nil {
		return nil, err
	}
	var output []string
	err = json.Unmarshal(body, &output)
	if err != nil {
		return nil, err
	}
	return output, nil
}

var LocalAccess engin.Access
var LocalForward engin.Forward
var RemoteForward engin.Forward

func LocalAccessInit(scheme string, address string, auth *engin.AuthInfo) error {
	var tlsEnable bool
	if scheme == "https" {
		tlsEnable = true
	}
	access, err := engin.NewHttpsAccess(address, Timeout, tlsEnable, CertFile, KeyFile)
	if err != nil {
		logs.Error(err.Error())
		return err
	}
	if auth != nil {
		access.AuthHandlerSet(func(info *engin.AuthInfo) bool {
			logs.Info("auth request %v", info)
			if info == nil {
				return false
			}
			if info.User == auth.User && info.Token == auth.Token {
				logs.Info("auth success")
				return true
			}
			logs.Info("auth fail")
			return false
		})
	}
	LocalAccess = access
	return nil
}

func RemoteForwardInit(scheme string, address string, auth *engin.AuthInfo) error {
	if net.ParseIP(address) == nil {
		addr, err := net.ResolveTCPAddr("tcp", address)
		if err != nil {
			return err
		}
		logs.Info("resolve %s to %s", address, addr.String())
		address = addr.String()
	}
	if false == engin.IsConnect(address, Timeout) {
		return fmt.Errorf("connect %s fail", address)
	}
	var tlsEnable bool
	if scheme == "https" {
		tlsEnable = true
	}
	forward, err := engin.NewHttpsProtcal(address, Timeout, auth, tlsEnable, CertFile, KeyFile)
	if err != nil {
		logs.Error(err.Error())
		return err
	}
	RemoteForward = forward
	return nil
}

func LocalForwardFunc(address string, r *http.Request) engin.Forward {
	return LocalForward
}

func ProxyForwardFunc(address string, r *http.Request) engin.Forward {
	return RemoteForward
}

func DomainForwardFunc(address string, r *http.Request) engin.Forward {
	if DomainCheck(address) {
		logs.Info("%s auto forward to remote proxy", address)
		return RemoteForward
	}
	return LocalForward
}

func AutoForwardFunc(address string, r *http.Request) engin.Forward {
	if AutoCheck(address) {
		return LocalForward
	}
	logs.Info("%s auto forward to remote proxy", address)
	return RemoteForward
}

func main() {
	flag.Parse()
	if Help {
		flag.Usage()
		return
	}

	err := LogInit(Debug)
	if err != nil {
		panic(err.Error())
	}

	scheme, address, err := parseAddress(LocalAddr)
	if err != nil {
		panic(err.Error())
	}

	auth, err := parseAuth(LocalAuth)
	if err != nil {
		panic(err.Error())
	}

	err = LocalAccessInit(scheme, address, auth)
	if err != nil {
		panic(err.Error())
	}

	LocalForward, _ = engin.NewDefault(Timeout)

	if strings.ToLower(RunMode) != "local" {
		scheme, address, err = parseAddress(RemoteAddr)
		if err != nil {
			panic(err.Error())
		}
		auth, err = parseAuth(RemoteAuth)
		if err != nil {
			panic(err.Error())
		}
		err = RemoteForwardInit(scheme, address, auth)
		if err != nil {
			panic(err.Error())
		}
	}

	switch strings.ToLower(RunMode) {
	case "domain":
		{
			domainList, err := parseDomain(DomainFile)
			if err != nil {
				panic(err)
			}
			DomainInit(domainList)
			LocalAccess.ForwardHandlerSet(DomainForwardFunc)
		}
	case "auto":
		{
			AutoInit()
			LocalAccess.ForwardHandlerSet(AutoForwardFunc)
		}
	case "proxy":
		LocalAccess.ForwardHandlerSet(ProxyForwardFunc)
	case "local":
		LocalAccess.ForwardHandlerSet(LocalForwardFunc)
	default:
		panic(fmt.Sprintf("running mode(%s) not support", RunMode))
	}

	logs.Info("autoproxy %s instance %s running mode %s success", VersionGet(), engin.GetUUID(), RunMode)

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGKILL, syscall.SIGINT, syscall.SIGTERM)

	sig := <-signalChan
	logs.Info("recv signal %s, ready to exit", sig.String())

	LocalAccess.Shutdown()
	LocalForward.Close()
	if RemoteForward != nil {
		RemoteForward.Close()
	}
	os.Exit(-1)
}
