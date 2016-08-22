package log

import (
	"fmt"
	"github.com/boringding/beekeeper/file"
	"io"
	"os"
	"strconv"
	"strings"
	"sync"
)

type RotateWriter struct {
	mu             sync.Mutex
	maxFileCnt     int
	maxFileSize    uint64
	fileNamePrefix string
	dir            string
	fileName       string
	curFileNo      int
	curFileSize    uint64
	file           *os.File
	Writer         io.Writer
}

func (self *RotateWriter) getPathFileName() (string, string) {
	fileName := self.fileNamePrefix + ".log"
	path := self.dir + fileName

	return path, fileName
}

func (self *RotateWriter) closeFile() error {
	defer func() {
		self.file = nil
	}()

	if self.file != nil {
		err := self.file.Close()
		if err != nil {
			return err
		}
	}

	return nil
}

func (self *RotateWriter) openFile() error {
	var path string
	path, self.fileName = self.getPathFileName()

	var err error
	self.file, err = os.OpenFile(path, os.O_CREATE|os.O_RDWR|os.O_APPEND, 0666)
	if err != nil {
		self.file = nil
		return err
	}

	fileSize, err := self.file.Seek(0, os.SEEK_END)
	if err != nil {
		self.file.Close()
		self.file = nil
		return err
	}

	if fileSize < 0 {
		self.curFileSize = 0
	} else {
		self.curFileSize = uint64(fileSize)
	}

	return nil
}

func (self *RotateWriter) getMaxFileNo(fileName string) (int, error) {
	maxFileNo := 0
	fileNameLen := len(fileName)

	file, err := os.Open(self.dir)
	if err != nil {
		return maxFileNo + 1, err
	}

	names, err := file.Readdirnames(-1)
	file.Close()
	if err != nil {
		return maxFileNo + 1, err
	}

	for _, name := range names {
		path := self.dir + name
		fileInfo, err := os.Lstat(path)

		if err != nil {
			continue
		}

		if fileInfo.IsDir() {
			continue
		}

		nameLen := len(name)
		pos := strings.Index(name, fileName)
		fileNo := 0

		if pos < 0 {
			continue
		}

		if nameLen > fileNameLen {
			fileNoStr := name[pos+fileNameLen+1 : nameLen]
			fileNo, err = strconv.Atoi(fileNoStr)

			if err != nil {
				continue
			}
		}

		if fileNo > maxFileNo {
			maxFileNo = fileNo
		}
	}

	return maxFileNo + 1, nil
}

func (self *RotateWriter) initCurFileNo() error {
	_, fileName := self.getPathFileName()

	var err error
	self.curFileNo, err = self.getMaxFileNo(fileName)
	if err != nil {
		return err
	}

	return nil
}

func (self *RotateWriter) updateCurFileNo() error {
	var err error
	self.curFileNo, err = self.getMaxFileNo(self.fileName)
	if err != nil {
		return err
	}

	return nil
}

func (self *RotateWriter) renameFile(src string, dst string) error {
	err := os.Remove(dst)
	if err != nil {
		return err
	}

	err = os.Rename(src, dst)
	if err != nil {
		return err
	}

	return nil
}

func (self *RotateWriter) removeFile(path string) error {
	return os.Remove(path)
}

func (self *RotateWriter) shiftFile() error {
	fd := int(self.file.Fd())

	file.LockFile(fd)
	defer file.UnlockFile(fd)

	path := self.dir + self.fileName

	err := file.CheckFile(path, self.file)
	if err != nil {
		self.closeFile()
		return self.openFile()
	}

	fileInfo, err := os.Stat(path)
	if err != nil {
		self.closeFile()
		return self.openFile()
	}

	self.curFileSize = uint64(fileInfo.Size())

	if self.curFileSize < self.maxFileSize {
		return nil
	}

	self.updateCurFileNo()

	src := path
	dst := fmt.Sprintf("%s.%d", src, self.curFileNo)

	self.renameFile(src, dst)

	rm := fmt.Sprintf("%s.%d", src, self.curFileNo-self.maxFileCnt)

	self.removeFile(rm)

	self.closeFile()
	return self.openFile()
}

func (self *RotateWriter) MaxFileCnt() int {
	self.mu.Lock()
	defer self.mu.Unlock()

	return self.maxFileCnt
}

func (self *RotateWriter) SetMaxFileCnt(maxFileCnt int) {
	self.mu.Lock()
	defer self.mu.Unlock()

	self.maxFileCnt = maxFileCnt
}

func (self *RotateWriter) MaxFileSize() uint64 {
	self.mu.Lock()
	defer self.mu.Unlock()

	return self.maxFileSize
}

func (self *RotateWriter) SetMaxFileSize(maxFileSize uint64) {
	self.mu.Lock()
	defer self.mu.Unlock()

	self.maxFileSize = maxFileSize
}

func (self *RotateWriter) FileNamePrefix() string {
	self.mu.Lock()
	defer self.mu.Unlock()

	return self.fileNamePrefix
}

func (self *RotateWriter) SetFileNamePrefix(fileNamePrefix string) {
	self.mu.Lock()
	defer self.mu.Unlock()

	self.fileNamePrefix = fileNamePrefix
}

func (self *RotateWriter) Dir() string {
	self.mu.Lock()
	defer self.mu.Unlock()

	return self.dir
}

func (self *RotateWriter) SetDir(dir string) {
	self.mu.Lock()
	defer self.mu.Unlock()

	self.dir = dir
}

func (self *RotateWriter) Write(p []byte) (n int, err error) {
	err = self.shiftFile()
	if err != nil {
		return 0, err
	}

	n, err = self.file.Write(p)
	if err != nil {
		return 0, err
	}

	self.curFileSize += uint64(n)
	return n, err
}
