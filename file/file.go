package file

import (
	"os"
)

func RemoveFile(path string) {
	os.Remove(path)
}

func RenameFile(src string, dst string) error {
	os.Remove(dst)

	return os.Rename(src, dst)
}
