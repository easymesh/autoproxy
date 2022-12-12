package main

import (
	"github.com/astaxie/beego/logs"
)

func LogInit(debug bool) error {
	logs.SetLogger(logs.AdapterConsole)
	if debug {
		logs.EnableFuncCallDepth(true)
		logs.SetLogFuncCallDepth(3)
	}
	return nil
}
