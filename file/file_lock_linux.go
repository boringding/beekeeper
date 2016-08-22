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
