package file

import (
	"errors"
	"os"
)

var ErrNotSameFile = errors.New("not the same file")

func CheckFile(path string, file *os.File) error {
	return nil
}