package main

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"sync"
	"time"

	"github.com/astaxie/beego/logs"
	"github.com/easymesh/autoproxy/engin"
)

const (
	CACHE_FILE = "autocache.json"
)

type LocalAccessInfo struct {
	Hostname string `json:"hostname"`
	Access   bool   `json:"access"`
}

type AutoCtrl struct {
	sync.RWMutex
	cache map[string]LocalAccessInfo
}

var autoCtrl AutoCtrl

func AutoInit() {
	autoCtrl.cache = make(map[string]LocalAccessInfo, 100)
	syncFromFile()
	go syncTask()
}

func cacheFileTimestamp() time.Time {
	info, err := os.Stat(CACHE_FILE)
	if err != nil {
		logs.Warning(err.Error())
		return time.Time{}
	}
	return info.ModTime()
}

func syncTask() {
	time1 := cacheFileTimestamp()
	for {
		time.Sleep(time.Second)
		time2 := cacheFileTimestamp()
		if time2 != time1 {
			autoCtrl.Lock()
			syncFromFile()
			autoCtrl.Unlock()
		}
		time1 = cacheFileTimestamp()
	}
}

func syncFromFile() {
	var list []LocalAccessInfo
	body, err := ioutil.ReadFile(CACHE_FILE)
	if err != nil {
		logs.Warning(err.Error())
		return
	}
	err = json.Unmarshal(body, &list)
	if err != nil {
		logs.Error(err.Error())
	}
	for _, access := range list {
		old_access, exist := autoCtrl.cache[access.Hostname]
		if exist {
			if old_access.Access != access.Access {
				logs.Info("update %s from %v to %v", access.Hostname, old_access.Access, access.Access)
			}
		}
		autoCtrl.cache[access.Hostname] = access
	}
	logs.Info("sync %d from cache file success", len(list))
}

func syncToFile() {
	var list []LocalAccessInfo
	for _, access := range autoCtrl.cache {
		list = append(list, access)
	}
	body, err := json.MarshalIndent(list, "", "\t")
	if err != nil {
		logs.Error(err.Error())
		return
	}
	err = ioutil.WriteFile(CACHE_FILE, body, 0644)
	if err != nil {
		logs.Error(err.Error())
	}
}

func AutoCheck(address string) bool {
	autoCtrl.RLock()
	result, ok := autoCtrl.cache[address]
	autoCtrl.RUnlock()
	if ok {
		return result.Access
	}

	result.Hostname = address
	result.Access = engin.IsConnect(address, 3)

	autoCtrl.Lock()
	autoCtrl.cache[address] = result
	syncToFile()
	autoCtrl.Unlock()

	return result.Access
}
