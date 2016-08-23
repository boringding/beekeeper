package log

import (
	"sync"
	"fmt"
	"io"
	"log"
	"github.com/boringding/beekeeper/conf"
	"github.com/boringding/beekeeper/proc"
)

const (
	LogAll     = iota
	LogDebug
	LogInfo
	LogWarn
	LogErr
	LogFatal
	LogNone
)

var LogLvls=[]string{"ALL","DEBUG","INFO","WARN","ERR","FATAL","NONE"}

type Log struct {
	mu           sync.Mutex
	lvl          int
	rotateWriter RotateWriter
	logger       log.Logger
}

func (self *Log) Init(conf conf.LogConf) error {
	self.mu.Lock()
	self.lvl = conf.Lvl
	self.mu.Unlock()
	
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
	
	return nil
}

func (self *Log) SetWriter(writer io.Writer) {
	self.logger.SetOutput(writer)
}

func (self *Log) Log(lvl int, format string, v ...interface{}) error {
	if lvl < self.lvl {
		return nil
	}
	
	prefix := formatPrefix(lvl)
	content := fmt.Sprintf(format, v...)
	output := prefix+content
	
	return self.logger.Output(2, fmt.Sprintln(output))
}

func formatPrefix(lvl int) string {
	selfPid := proc.GetSelfPid()
	lvlStr := LogLvls[lvl]
	prefix := fmt.Sprintf("[%d][%s]", selfPid, lvlStr)
	
	return prefix
}
