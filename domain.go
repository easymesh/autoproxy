package main

import (
	"encoding/json"
	"io/ioutil"
	"strings"
	"sync"
	"time"

	"github.com/astaxie/beego/logs"
	"github.com/easymesh/autoproxy/engin"
)

type DomainCtrl struct {
	sync.RWMutex
	cache  map[string]string
	domain []string
}

var forwardCtrl DomainCtrl

func DomainInit(file string) {
	forwardCtrl.cache = make(map[string]string, 100)
	forwardCtrl.domain = make([]string, 100)
	domainFromFile(file)
	go syncDomainTask(file)
}

func syncDomainTask(file string) {
	time1 := engin.GetFileTimestamp(file)
	for {
		time.Sleep(time.Second)
		time2 := engin.GetFileTimestamp(file)
		if time2 != time1 {
			domainFromFile(file)
		}
		time1 = time2
	}
}

func domainFromFile(file string) {
	forwardCtrl.Lock()
	defer forwardCtrl.Unlock()

	body, err := ioutil.ReadFile(file)
	if err != nil {
		logs.Warning(err.Error())
		return
	}
	var domain []string
	err = json.Unmarshal(body, &domain)
	if err != nil {
		logs.Warning(err.Error())
		return
	}

	forwardCtrl.cache = make(map[string]string, 100)
	forwardCtrl.domain = domain

	logs.Info("sync %d from domain file success", len(domain))
}

func domainGet(address string) string {
	domain := address
	idx := strings.Index(address, ":")
	if idx != -1 {
		domain = address[:idx]
	}
	return domain
}

func domainMatch(domain string, match string) bool {
	begin := strings.Index(match, "*")
	end := strings.Index(match[begin+1:], "*")
	if end != -1 {
		end += begin + 1
	}
	if begin != -1 && end == -1 {
		// suffix match
		return strings.HasSuffix(domain, match[begin+1:])
	}
	if begin == -1 && end != -1 {
		// prefix match
		return strings.HasPrefix(domain, match[:end])
	}
	if begin == -1 && end == -1 {
		// full match
		if domain == match {
			return true
		} else {
			return false
		}
	}
	idx := strings.Index(domain, match[begin+1:end])
	return idx != -1
}

// address: www.baidu.com:80 or www.baidu.com:443
func routeMatch(address string) string {
	domain := domainGet(address)
	for _, v := range forwardCtrl.domain {
		if domainMatch(domain, v) {
			forwardCtrl.cache[address] = v
			logs.Info("route address %s match to domain %s", address, v)
			return v
		}
	}
	logs.Info("route address %s no match", address)
	forwardCtrl.cache[address] = ""
	return ""
}

func DomainCheck(address string) bool {
	forwardCtrl.RLock()
	result, flag := forwardCtrl.cache[address]
	forwardCtrl.RUnlock()

	if !flag {
		forwardCtrl.Lock()
		result = routeMatch(address)
		forwardCtrl.Unlock()
	}

	if result == "" {
		return false
	}

	return true
}
