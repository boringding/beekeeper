package file

import (
	"errors"
	"os"
	"reflect"
)

var ErrNotSameFile = errors.New("not the same file")

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

	pathInode := reflect.ValueOf(pathFileStat).Elem().FieldByName("Ino").Field(0).Uint()
	fileInode := reflect.ValueOf(fileStat).Elem().FieldByName("Ino").Field(0).Uint()

	if pathInode != fileInode {
		return ErrNotSameFile
	}

	return nil
}
