package log

import (
	"fmt"
	"testing"
)

func Test_getMaxFileNo(t *testing.T) {
	var w RotateWriter
	w.dir = "D:\\"
	maxFileNo, err := w.getMaxFileNo("test.log")
	fmt.Println(maxFileNo, err)
}

func Test_openFile(t *testing.T) {
	var w RotateWriter
	w.dir = "D:\\"
	w.fileNamePrefix = "test"
	err := w.openFile()
	w.closeFile()
	fmt.Println(w.curFileSize, err)
}
