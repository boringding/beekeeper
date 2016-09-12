package file

import (
	"errors"
	"os"
	"reflect"
)

func CheckFile(path string, file *os.File) error {
	pathFileInfo, err := os.Stat(path)
	if err != nil {
		return err
	}

	pathFileStat := pathFileInfo.Sys()

	fileInfo, err := file.Stat()
	if err != nil {
		return err
	}

	fileStat := fileInfo.Sys()

	pathInode := reflect.ValueOf(pathFileStat).Elem().FieldByName("Ino").Uint()
	fileInode := reflect.ValueOf(fileStat).Elem().FieldByName("Ino").Uint()

	if pathInode != fileInode {
		return errors.New("not the same file")
	}

	return nil
}
