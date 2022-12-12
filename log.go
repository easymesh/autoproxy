package main

import (
	"fmt"

	"github.com/astaxie/beego/logs"
)

func LogInit(file string) {
	if file == "" {
		logs.SetLogger(logs.AdapterConsole)
	} else {
		logs.SetLogger(logs.AdapterFile, fmt.Sprintf("{\"filename\":\"%s\"}", file))
	}
}
