package main

import (
	"strings"
	"sync"

	"github.com/astaxie/beego/logs"
)

type DomainCtrl struct {
	sync.RWMutex
	cache  map[string]string
	domain []string
}

var forwardCtrl DomainCtrl

func DomainInit(domain []string) {
	forwardCtrl.cache = make(map[string]string, 100)
	forwardCtrl.domain = domain
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
	if idx == -1 {
		return false
	}
	return true
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

	if flag == false {
		forwardCtrl.Lock()
		result = routeMatch(address)
		forwardCtrl.Unlock()
	}

	if result == "" {
		return false
	}

	return true
}
