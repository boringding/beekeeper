//Lock the file represented by the file descriptor.
//In linux it is implemented by system call.

package file

import (
	"syscall"
)

func LockFile(fd int) error {
	return syscall.Flock(fd, syscall.LOCK_EX)
}

func UnlockFile(fd int) error {
	return syscall.Flock(fd, syscall.LOCK_UN)
}
