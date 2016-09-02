package proc

import (
	"os"
)

var pid = 0
var ppid = 0

func init() {
	pid = os.Getpid()
	ppid = os.Getppid()
}

func GetSelfPid() int {
	return pid
}

func GetParentPid() int {
	return ppid
}
