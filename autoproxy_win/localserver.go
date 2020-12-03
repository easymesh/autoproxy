package main

import (
	"fmt"
	"github.com/astaxie/beego/logs"
	"github.com/easymesh/autoproxy/engin"
	"net/http"
	"strings"
	"sync"
)

var access engin.Access

func StatGet() (uint64, uint64) {
	acc := access
	if acc != nil {
		return acc.Stat()
	}
	return 0,0
}

/*
func AuthSwitch(auth *engin.AuthInfo) bool {
	if auth == nil {
		return false
	}
	return AuthCheck(auth.User, auth.Token)
}*/

var LocalForward engin.Forward
var RemoteForward engin.Forward

var mutex sync.Mutex

func LocalForwardFunc(address string, r *http.Request) engin.Forward {
	return LocalForward
}

func ProxyForwardFunc(address string, r *http.Request) engin.Forward {
	return RemoteForward
}

func AutoForwardFunc(address string, r *http.Request) engin.Forward {
	if RouteCheck(address) {
		logs.Info("%s auto forward to remote proxy", address)
		return RemoteForward
	}
	return LocalForward
}

func RemoteForwardUpdate() error {
	mutex.Lock()
	defer mutex.Unlock()

	return remoteUpdate()
}

func remoteUpdate() error {
	list := RemoteList()
	if len(list) == 0 {
		logs.Warn("no remote proxy server.")
		return nil
	}

	remote := RemoteList()[RemoteIndexGet()]
	logs.Info("remote swtich config : %v", remote)

	var tlsEnable bool
	if strings.ToLower(remote.Protocal) == "https" {
		tlsEnable = true
	}

	var auth *engin.AuthInfo
	if remote.Auth {
		auth = &engin.AuthInfo{User: remote.User, Token: remote.Password}
	}

	forward, err := engin.NewHttpsProtcal(remote.Address, 60, auth, tlsEnable )
	if err != nil {
		logs.Error(err.Error())
		return err
	}

	logs.Info("remote swtich to %s success", remote.Name )

	if RemoteForward != nil {
		RemoteForward.Close()
	}
	RemoteForward = forward
	return nil
}

func modeUpdate() error {
	acc := access
	if acc == nil {
		logs.Warn("server has been stop, mode update disable")
		return nil
	}

	mode := ModeOptionGet()
	if mode != "local" && len(RemoteList()) == 0 {
		return fmt.Errorf("Please add remote proxy config.")
	}

	logs.Info("mode switch to %s", mode)

	switch mode {
	case "auto":
		acc.ForwardHandlerSet(AutoForwardFunc)
	case "proxy" :
		acc.ForwardHandlerSet(ProxyForwardFunc)
	case "local":
		acc.ForwardHandlerSet(LocalForwardFunc)
	}

	logs.Info("server mode switch to %s success", mode)
	return nil
}

func ModeUpdate() error {
	mutex.Lock()
	defer mutex.Unlock()

	return modeUpdate()
}

func ServerStart() error {
	mutex.Lock()
	defer mutex.Unlock()

	var err error

	if access != nil {
		logs.Error("server has been start")
		return fmt.Errorf("server has been start")
	}

	address := fmt.Sprintf("%s:%d",
		IfaceOptions()[LocalIfaceOptionsIdx()],
		PortOptionGet())

	logs.Info("server start %s", address)

	access, err = engin.NewHttpsAccess(address, 60, false)
	if err != nil {
		logs.Error(err.Error())
		return err
	}

	LocalForward, _ = engin.NewDefault(60)

	err = remoteUpdate()
	if err != nil {
		return err
	}

	modeUpdate()

	logs.Info("server start %s success", address)
	return nil
}

func ServerRunning() bool {
	mutex.Lock()
	defer mutex.Unlock()

	if access == nil {
		return false
	}
	return true
}

func ServerShutdown() error {
	mutex.Lock()
	defer mutex.Unlock()

	if access == nil {
		return fmt.Errorf("server has been stop")
	}
	err := access.Shutdown()
	if err != nil {
		logs.Error("shutdown fail, %s", err.Error())
		return err
	}
	access = nil

	LocalForward.Close()
	LocalForward = nil
	return nil
}
