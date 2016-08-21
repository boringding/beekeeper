package log

import (
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
)

type writer struct {
	mu             sync.Mutex
	maxFileCnt     uint
	maxFileSize    uint64
	fileNamePrefix string
	dir            string
	fileName       string
	curFileNo      uint
	curFileSize    uint64
	file           *os.File
}

func (self *writer) closeFile() error {
	self.mu.Lock()
	defer func() {
		self.file = nil
		self.mu.Unlock()
	}()

	if self.file != nil {
		err := self.file.Close()
		if err != nil {
			return err
		}
	}

	return nil
}

func (self *writer) openFile() error {
	self.mu.Lock()
	defer self.mu.Unlock()

	self.fileName = self.fileNamePrefix + ".log"
	path := self.dir + self.fileName

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

func (self *writer) getMaxFileNo(fileName string) (maxFileNo uint, err error) {
	self.mu.Lock()
	defer self.mu.Unlock()

	maxFileNo = 0

	err = filepath.Walk(self.dir, func(path string, fileInfo os.FileInfo, err error) error {
		if fileInfo == nil || err != nil {
			return err
		}

		if fileInfo.IsDir() {
			return nil
		}

		_fileName := fileInfo.Name()
		length := len(fileName)
		_length := len(_fileName)
		pos := strings.Index(_fileName, fileName)
		var fileNo uint

		if pos < 0 {
			return nil
		}

		if _length > length {
			_fileNo := _fileName[pos+length+1 : _length-pos-length-1]
			fileNo, _err := strconv.Atoi(_fileNo)

			if _err != nil {
				return nil
			}
		}

		if fileNo > maxFileNo {
			maxFileNo = fileNo
		}

		return nil
	})

	if err != nil {
		return maxFileNo + 1, err
	}

	return maxFileNo + 1, nil
}
