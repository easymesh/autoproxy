package util

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

var logCfg = logconfig{
	Level: logs.LevelInformational,
	Daily: true,
	MaxSize: 10*1024*1024,
	MaxLines: 100*1024,
	MaxDays: 7,
	Color: false,
}

func NewLogger(logdir string, filename string) (*logs.BeeLogger, error) {
	os.MkdirAll(logdir, 0644)
	
	logCfg.Filename = fmt.Sprintf("%s%c%s", logdir, os.PathSeparator, filename)
	value, err := json.Marshal(&logCfg)
	if err != nil {
		return nil, err
	}
	logger := logs.NewLogger(100)
	err = logger.SetLogger(logs.AdapterFile, string(value))
	if err != nil {
		return nil, err
	}
	logger.Async(100)
	logger.EnableFuncCallDepth(true)
	logger.SetLogFuncCallDepth(3)
	return logger, nil
}
