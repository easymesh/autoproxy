package engin

import (
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/astaxie/beego/logs"
)

type AuthItem struct {
	Address string
	Login   time.Time
}

type AuthCtrl struct {
	sync.RWMutex
	Cache map[string]AuthItem
}

var authctrl AuthCtrl

func init() {
	authctrl.Cache = make(map[string]AuthItem, 100)
	go func() {
		for {
			now := time.Now()
			authctrl.Lock()
			for _, v := range authctrl.Cache {
				if now.Sub(v.Login) > time.Hour {
					logs.Info("auth ctrl address %s lease timeout", v.Address)
					delete(authctrl.Cache, v.Address)
				}
			}
			authctrl.Unlock()
			time.Sleep(time.Minute)
		}
	}()
}

func AuthLogin(r *http.Request) {
	address := strings.Split(r.RemoteAddr, ":")[0]

	authctrl.Lock()
	defer authctrl.Unlock()

	authctrl.Cache[address] = AuthItem{Address: address, Login: time.Now()}
}

func AuthCache(r *http.Request) bool {
	address := strings.Split(r.RemoteAddr, ":")[0]

	authctrl.RLock()
	defer authctrl.RUnlock()

	_, ok := authctrl.Cache[address]
	if ok {
		return true
	}
	return false
}
