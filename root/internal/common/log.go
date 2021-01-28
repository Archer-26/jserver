package common

import (
	"root/pkg/iniconfig"
	"root/pkg/log"
	"time"
)

var appName string
var appId int

func InitLog(level int32, app string, appid int) {
	appName = app
	appId = appid
	log.Init(level, reportLog, iniconfig.String("logdir"), app, appid)
}

func reportLog(t time.Time, level string, file string, line int, msg string) {
	if log.Levels[level] < log.TAG_WARN_I {
		return
	}
}
