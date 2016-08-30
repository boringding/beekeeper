package proc

import (
	"os"
)

var pid = -1
var ppid = -1

func GetSelfPid() int {
	if pid >= 0 {
		return pid
	} else {
		return os.Getpid()
	}
}

func GetParentPid() int {
	if ppid >= 0 {
		return ppid
	} else {
		return os.Getppid()
	}
}
