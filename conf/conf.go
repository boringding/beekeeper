package conf

import (
	"encoding/xml"
	"os"
)

func Parse(path string, v interface{}) error {
	file, err := os.OpenFile(path, os.O_RDONLY, 0666)
	if err != nil {
		return err
	}

	content := []byte{}

	_, err = file.Read(content)
	if err != nil {
		return err
	}

	return xml.Unmarshal(content, v)
}
