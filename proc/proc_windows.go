//Terminate a process in windows using external command.

package proc

import (
	"os/exec"
	"strconv"
)

func TerminateProc(pid int) error {
	pidStr := strconv.Itoa(pid)
	cmd := exec.Command("taskkill.exe", "/f", "/pid", pidStr)
	return cmd.Start()
}
