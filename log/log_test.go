package log

import (
	"testing"
	"github.com/boringding/beekeeper/conf"
)

func Test_Log(t *testing.T){
	var l Log
	var logConf conf.LogConf
	logConf.MaxFileCnt = 2
	logConf.MaxFileSize = 1024
	logConf.FileNamePrefix = "doubi"
	logConf.Dir = "D:\\"
	logConf.Lvl = LogErr
	l.Init(logConf)
	
	l.Log(LogErr, "%d-%s", 1, "test")
	l.Log(LogFatal, "%d-%s", 2, "test")
	l.Log(LogErr, "%d-%s", 3, "test")
}