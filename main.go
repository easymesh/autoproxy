package main

import (
	"flag"
	"fmt"
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
	Timeout int

	LocalAddr  string
	LocalAuth  string
	RemoteAddr string
	RemoteAuth string

	RunMode    string // local、proxy、domain、auto
	DomainFile string
	CertFile   string
	KeyFile    string
	LogFile    string
)

func init() {
	flag.StringVar(&KeyFile, "key-file", "", "tls key file pem format")
	flag.StringVar(&CertFile, "cert-file", "", "tls cert file pem format")

	flag.IntVar(&Timeout, "timeout", 30, "connect timeout (unit second)")

	flag.StringVar(&LocalAddr, "local-address", "http://0.0.0.0:8080", "Local proxy listening address")
	flag.StringVar(&LocalAuth, "local-auth", "", "Local proxy auth username and password")

	flag.StringVar(&RemoteAddr, "remote-address", "https://my.domain:8080", "Remote proxy listening address")
	flag.StringVar(&RemoteAuth, "remote-auth", "", "Remote proxy auth username and password")

	flag.StringVar(&RunMode, "mode", "proxy", "proxy mode(local/proxy/domain/auto)")
	flag.StringVar(&DomainFile, "domain", "domain.json", "match domain list file(domain mode requires)")

	flag.StringVar(&LogFile, "logfile", "autoproxy.log", "logger file")
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
	return strings.ToLower(ul.Scheme), engin.Address(ul), nil
}

func LocalAccessInit(scheme string, address string, auth *engin.AuthInfo) (engin.Access, error) {
	var tlsEnable bool
	if scheme == "https" {
		tlsEnable = true
	}
	access, err := engin.NewHttpsAccess(address, Timeout, tlsEnable, CertFile, KeyFile)
	if err != nil {
		logs.Error(err.Error())
		return nil, err
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
	return access, nil
}

func RemoteForwardInit(scheme string, address string, auth *engin.AuthInfo) (engin.Forward, error) {
	var tlsEnable bool
	if scheme == "https" {
		tlsEnable = true
	}
	forward, err := engin.NewHttpProxyForward(address, Timeout, auth, tlsEnable, CertFile, KeyFile)
	if err != nil {
		logs.Error(err.Error())
		return nil, err
	}
	return forward, nil
}

func main() {
	flag.Parse()
	if Help {
		flag.Usage()
		return
	}

	LogInit(LogFile)

	scheme, address, err := parseAddress(LocalAddr)
	if err != nil {
		panic(err.Error())
	}

	auth, err := parseAuth(LocalAuth)
	if err != nil {
		panic(err.Error())
	}

	var acc engin.Access
	acc, err = LocalAccessInit(scheme, address, auth)
	if err != nil {
		panic(err.Error())
	}

	var local, proxy engin.Forward

	local = engin.NewLocalForward(Timeout)

	if strings.ToLower(RunMode) != "local" {
		scheme, address, err = parseAddress(RemoteAddr)
		if err != nil {
			panic(err.Error())
		}
		auth, err = parseAuth(RemoteAuth)
		if err != nil {
			panic(err.Error())
		}
		proxy, err = RemoteForwardInit(scheme, address, auth)
		if err != nil {
			panic(err.Error())
		}
	}

	switch strings.ToLower(RunMode) {
	case "domain":
		{
			DomainInit(DomainFile)
			acc.ForwardHandlerSet(func(address string, r *http.Request) engin.Forward {
				if DomainCheck(address) {
					logs.Info("%s auto forward to remote proxy", address)
					return proxy
				}
				return local
			})
		}
	case "auto":
		{
			AutoInit()
			acc.ForwardHandlerSet(func(address string, r *http.Request) engin.Forward {
				if AutoCheck(address) {
					logs.Info("%s auto forward to local network", address)
					return local
				}
				logs.Info("%s auto forward to remote proxy", address)
				return proxy
			})
			acc.ForwardUpdateHandlerSet(func(address string, forward engin.Forward) {
				if forward == local {
					AutoCheckUpdate(address, false)
				}
				if forward == proxy {
					AutoCheckUpdate(address, true)
				}
			})
		}
	case "proxy":
		acc.ForwardHandlerSet(func(address string, r *http.Request) engin.Forward {
			return proxy
		})
	case "local":
		acc.ForwardHandlerSet(func(address string, r *http.Request) engin.Forward {
			return local
		})
	default:
		panic(fmt.Sprintf("running mode(%s) not support", RunMode))
	}

	logs.Info("autoproxy %s instance %s running mode %s success", VersionGet(), engin.GetUUID(), RunMode)

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGKILL, syscall.SIGINT, syscall.SIGTERM)

	sig := <-signalChan
	logs.Info("recv signal %s, ready to exit", sig.String())

	acc.Shutdown()
	local.Close()
	if proxy != nil {
		proxy.Close()
	}
	os.Exit(-1)
}
