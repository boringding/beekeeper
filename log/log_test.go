package log

import (
	"github.com/boringding/beekeeper/conf"
	"testing"
)

func Test_Log(t *testing.T) {
	var l Log
	var logConf conf.LogConf
	logConf.MaxFileCnt = 2
	logConf.MaxFileSize = 1024
	logConf.FileNamePrefix = "doubi"
	logConf.Dir = "D:\\"
	logConf.Lvl = LogErr
	l.Init(logConf)
	//l.SetWriter(os.Stdout)

	l.Log(LogErr, "%d-%s", 1, "测试")
	l.Log(LogFatal, "%d-%s", 2, "我们")
	l.Log(LogErr, "%d-%s", 3, "程序员")
}
