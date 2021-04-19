package pages

import (
	"fmt"
	"github.com/astaxie/beego/logs"
	"github.com/easymesh/autoproxy/console/models"
	util "github.com/easymesh/autoproxy/console/uitl"
	"github.com/easymesh/autoproxy/engin"
	"github.com/GoAdminGroup/go-admin/modules/logger"
	"net"
	"net/http"
	"net/url"
	"strings"
)

type ProxyEngin struct {
	proxy  *models.Proxy
	remote *models.Remote
	access        engin.Access
	localForward  engin.Forward
	remoteForward engin.Forward
}

func parseAddress(protocal string, address string) (string, int, error) {
	idx := strings.Index(address, ":")
	if -1 == idx {
		switch strings.ToLower(protocal) {
		case "https":
			{
				return address, 443, nil
			}
		case "http":
			{
				return address, 80, nil
			}
		case "sock5":
			{
				return address, 1080, nil
			}
		}
		return "", 0, fmt.Errorf("protocal not support")
	}
	return address[:idx], util.Atoi(address[idx + 1:]), nil
}

func (p *ProxyEngin)authHandler(info *engin.AuthInfo) bool {
	if info == nil {
		logger.Warnf("not any auth info")
		return false
	}
	if user := models.UserFind(info.User); user != nil && user.Password == info.Token {
		if user.Enable == 0 {
			logger.Warnf("%s auth disable", info)
			return false
		}
		logger.Infof("%s auth success", info)
		return true
	}
	logger.Warnf("%s auth fail", info)
	return false
}

func (p *ProxyEngin)AuthSwitch(auth int)  {
	if auth > 0 {
		p.access.AuthHandlerSet(p.authHandler)
	} else {
		p.access.AuthHandlerSet(nil)
	}
	p.proxy.Auth = auth
}

func (p *ProxyEngin)Stop()  {
	p.access.Shutdown()
	if p.remoteForward != nil {
		p.remoteForward.Close()
	}
}

func remoteForwardInit(remote *models.Remote) (engin.Forward, error) {
	address, port , err := parseAddress(remote.Protocal, remote.Address)
	if err != nil {
		return nil, err
	}
	address = fmt.Sprintf("%s:%d", address, port)

	if net.ParseIP(address) == nil {
		addr, err := net.ResolveTCPAddr("tcp", address)
		if err != nil {
			return nil, err
		}
		address = addr.String() // resolve domain to ipv4
	}

	// try to connect first
	if false == engin.IsConnect(address, 5) {
		return nil, fmt.Errorf("connect %s fail", address)
	}

	logger.Infof("remote address %s", address)

	var tlsEnable bool
	if remote.Protocal == "https" {
		tlsEnable = true
	}

	var auth *engin.AuthInfo
	if remote.Auth > 0 {
		auth = &engin.AuthInfo{User: remote.User, Token: remote.Password}
	}

	forward, err := engin.NewHttpsProtcal(address, 30, auth, tlsEnable )
	if err != nil {
		return nil, err
	}

	return forward, nil
}

func remoteForwardTest(testurl string, forward engin.Forward) error {
	urls, err := url.Parse(testurl)
	if err != nil {
		return err
	}

	request, err := http.NewRequest("GET", testurl, nil)
	if err != nil {
		return err
	}

	if strings.ToLower(urls.Scheme) == "https" {
		conn, err := forward.Https(engin.Address(urls), request)
		if err != nil {
			return err
		}
		conn.Close()
	} else {
		rsp, err := forward.Http(request)
		if err != nil {
			return err
		}
		rsp.Body.Close()
	}

	return nil
}

func (p *ProxyEngin)DomainForwardFunc(address string, r *http.Request) engin.Forward {
	if DomainCheck(address) {
		logs.Info("%s auto forward to remote proxy", address)
		return p.remoteForward
	}
	return PublicLocalForward
}

func (p *ProxyEngin)LocalForwardFunc(address string, r *http.Request) engin.Forward {
	return PublicLocalForward
}

func (p *ProxyEngin)ProxyForwardFunc(address string, r *http.Request) engin.Forward {
	return p.remoteForward
}

func NewProxyEngin(proxy *models.Proxy, remote *models.Remote) (*ProxyEngin, error) {
	if proxy.Mode != models.MODE_LOCAL && remote == nil {
		return nil, fmt.Errorf("remote config not exist")
	}

	var err error
	var forword engin.Forward
	if remote != nil {
		forword, err = remoteForwardInit(remote)
		if err != nil {
			return nil, err
		}
		err = remoteForwardTest("https://www.google.com/", forword)
		if err != nil {
			models.RemoteUpdate(remote.Tag, func(u *models.Remote) {
				u.Status = err.Error()
			})
			return nil, err
		}
		models.RemoteUpdate(remote.Tag, func(u *models.Remote) {
			u.Status = "connected"
		})
	}

	var tlsEnable bool
	if proxy.Protocal == "https" {
		tlsEnable = true
	}

	var access engin.Access

	address := fmt.Sprintf("%s:%d", proxy.Iface, proxy.Port)
	access, err = engin.NewHttpsAccess(address, 30, tlsEnable)
	if err != nil {
		return nil, err
	}

	proxyEngin := new(ProxyEngin)
	proxyEngin.proxy = proxy
	proxyEngin.remote = remote
	proxyEngin.access = access
	proxyEngin.AuthSwitch(proxy.Auth)

	if proxy.Mode != models.MODE_LOCAL && remote == nil {
		return nil, fmt.Errorf("remote config not exist")
	}

	proxyEngin.remoteForward = forword

	switch strings.ToLower(proxy.Mode) {
		case models.MODE_DOMAIN:
			access.ForwardHandlerSet(proxyEngin.DomainForwardFunc)
		case models.MODE_REMOTE:
			access.ForwardHandlerSet(proxyEngin.ProxyForwardFunc)
		case models.MODE_LOCAL:
			access.ForwardHandlerSet(proxyEngin.LocalForwardFunc)
	default:
		panic(fmt.Sprintf("proxy mode(%s) not support", proxy.Mode))
	}

	return proxyEngin, nil
}

var PublicLocalForward engin.Forward

func init()  {
	PublicLocalForward, _ = engin.NewDefault(30)
}

