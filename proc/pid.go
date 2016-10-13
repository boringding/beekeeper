//Process id operations.

package proc

import (
	"os"
	"strconv"
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

func DumpSelfPid(path string) error {
	file, err := os.OpenFile(path, os.O_CREATE|os.O_RDWR, 0666)
	if err != nil {
		return err
	}

	pidStr := strconv.Itoa(pid)

	_, err = file.WriteString(pidStr)
	return err
}
