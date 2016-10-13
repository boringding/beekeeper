//Type RotateLog implements a log facility that limits
//log file size and count.
//It is based on log package in standard library.

package log

import (
	"fmt"
	"io"
	"log"
	"sync"

	"github.com/boringding/beekeeper/conf"
	"github.com/boringding/beekeeper/proc"
)

const (
	LogAll = iota
	LogDebug
	LogInfo
	LogWarn
	LogErr
	LogFatal
	LogNone
)

var LogLvls = []string{"ALL", "DEBUG", "INFO", "WARN", "ERR", "FATAL", "NONE"}

type RotateLog struct {
	mu           sync.RWMutex
	lvl          int
	rotateWriter RotateWriter
	logger       log.Logger
}

func NewRotateLog() *RotateLog {
	return &RotateLog{
		lvl: LogAll,
	}
}

func (self *RotateLog) Init(conf conf.LogConf) error {
	for i, v := range LogLvls {
		if v == conf.Lvl {
			self.lvl = i
			break
		}
	}

	self.rotateWriter.SetMaxFileCnt(conf.MaxFileCnt)
	self.rotateWriter.SetMaxFileSize(conf.MaxFileSize)
	self.rotateWriter.SetFileNamePrefix(conf.FileNamePrefix)
	self.rotateWriter.SetDir(conf.Dir)

	err := self.rotateWriter.Init()
	if err != nil {
		return err
	}

	self.logger.SetFlags(log.Ldate | log.Ltime | log.Lmicroseconds | log.Lshortfile)
	self.logger.SetOutput(&self.rotateWriter)
	self.logger.SetPrefix(fmt.Sprintf("[%d]", proc.GetSelfPid()))

	return nil
}

func (self *RotateLog) Lvl() int {
	self.mu.RLock()
	defer self.mu.RUnlock()

	return self.lvl
}

func (self *RotateLog) SetLvl(lvl int) {
	self.mu.Lock()
	defer self.mu.Unlock()

	if lvl < LogAll || lvl > LogNone {
		return
	}

	self.lvl = lvl
}

func (self *RotateLog) SetWriter(writer io.Writer) {
	self.logger.SetOutput(writer)
}

//Call this function with a callDepth=2.
func (self *RotateLog) Log(lvl int, callDepth int, format string, v ...interface{}) error {
	if lvl < self.lvl {
		return nil
	}

	lvlStr := fmt.Sprintf("[%s]", LogLvls[lvl])
	content := fmt.Sprintf(format, v...)
	output := lvlStr + content

	return self.logger.Output(callDepth, fmt.Sprintln(output))
}

func (self *RotateLog) Debug(format string, v ...interface{}) error {
	return self.Log(LogDebug, 3, format, v...)
}

func (self *RotateLog) Info(format string, v ...interface{}) error {
	return self.Log(LogInfo, 3, format, v...)
}

func (self *RotateLog) Warn(format string, v ...interface{}) error {
	return self.Log(LogWarn, 3, format, v...)
}

func (self *RotateLog) Err(format string, v ...interface{}) error {
	return self.Log(LogErr, 3, format, v...)
}

func (self *RotateLog) Fatal(format string, v ...interface{}) error {
	return self.Log(LogFatal, 3, format, v...)
}
