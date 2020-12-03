package main

import (
	"encoding/json"
	"fmt"
	"github.com/astaxie/beego/logs"
	"os"
)

type logconfig struct {
	Filename string  `json:"filename"`
	Level    int     `json:"level"`
	MaxLines int     `json:"maxlines"`
	MaxSize  int     `json:"maxsize"`
	Daily    bool    `json:"daily"`
	MaxDays  int     `json:"maxdays"`
	Color    bool    `json:"color"`
}

var logCfg = logconfig{Filename: os.Args[0], Level: logs.LevelInformational, Daily: true, MaxDays: 30, Color: false}

func LogInit(logpath string, debug bool) error {
	os.MkdirAll(logpath, 0644)
	logCfg.Filename = fmt.Sprintf("%s%c%s", logpath, os.PathSeparator, "autoproxy.log")
	value, err := json.Marshal(&logCfg)
	if err != nil {
		return err
	}
	if debug {
		err = logs.SetLogger(logs.AdapterConsole)
	} else {
		err = logs.SetLogger(logs.AdapterFile, string(value))
	}
	if err != nil {
		return err
	}
	logs.EnableFuncCallDepth(true)
	logs.SetLogFuncCallDepth(3)
	return nil
}
