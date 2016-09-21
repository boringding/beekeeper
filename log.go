package beekeeper

import (
	"io"

	"github.com/boringding/beekeeper/conf"
	"github.com/boringding/beekeeper/log"
)

var defaultLog = log.NewRotateLog()

func InitLog(conf conf.LogConf) error {
	return defaultLog.Init(conf)
}

func LogLvl() int {
	return defaultLog.Lvl()
}

func SetLogLvl(lvl int) {
	defaultLog.SetLvl(lvl)
}

func SetLogWriter(writer io.Writer) {
	defaultLog.SetWriter(writer)
}

func Log(lvl int, format string, v ...interface{}) error {
	return defaultLog.Log(lvl, 3, format, v...)
}

func LogDebug(format string, v ...interface{}) error {
	return defaultLog.Log(log.LogDebug, 3, format, v...)
}

func LogInfo(format string, v ...interface{}) error {
	return defaultLog.Log(log.LogInfo, 3, format, v...)
}

func LogWarn(format string, v ...interface{}) error {
	return defaultLog.Log(log.LogWarn, 3, format, v...)
}

func LogErr(format string, v ...interface{}) error {
	return defaultLog.Log(log.LogErr, 3, format, v...)
}

func LogFatal(format string, v ...interface{}) error {
	return defaultLog.Log(log.LogFatal, 3, format, v...)
}
