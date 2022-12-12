package main

import (
	"github.com/astaxie/beego/logs"
)

func LogInit(debug bool) error {
	if debug {
		logs.SetLogger(logs.AdapterConsole)
		logs.EnableFuncCallDepth(true)
		logs.SetLogFuncCallDepth(3)
	}
	return nil
}
