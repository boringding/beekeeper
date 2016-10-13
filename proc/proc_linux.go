//Terminate a process in linux using system call.

package proc

import (
	"syscall"
)

func TerminateProc(pid int) error {
	return syscall.Kill(pid, syscall.SIGTERM)
}
