package proc

import (
	"os"
)

func GetSelfPid() int {
	return os.Getpid()
}