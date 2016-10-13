//Unified configure load and parse
//including command line parameters and xml configure files.

package conf

import (
	"encoding/xml"
	"errors"
	"flag"
	"fmt"
	"os"
	"reflect"
	"strings"
	"sync"
)

const (
	CmdConfName  = "cmd"
	UsageTagName = "usage"
	//Default max size of configure file in bytes.
	ConfFileMaxSize = 1024 * 1024
)

type Conf struct {
	mu    sync.RWMutex
	items map[string]interface{}
	//The environment name used in the name suffix of configure files.
	env string
	//The directory of configure files.
	dir string
}

func NewConf() *Conf {
	return &Conf{
		items: map[string]interface{}{},
		env:   "",
		dir:   "",
	}
}

func (self *Conf) Init(env string, dir string) {
	self.env = env
	self.dir = dir
}

func (self *Conf) SetEnv(env string) {
	self.mu.Lock()
	defer self.mu.Unlock()

	self.env = env
}

func (self *Conf) SetDir(dir string) {
	self.mu.Lock()
	defer self.mu.Unlock()

	self.dir = dir
}

func (self *Conf) AddItem(name string, v interface{}) error {
	self.mu.Lock()
	defer self.mu.Unlock()

	if _, ok := self.items[name]; ok {
		return errors.New("configure item already exists")
	}

	self.items[name] = v

	return nil
}

func (self *Conf) Parse() error {
	self.mu.Lock()
	defer self.mu.Unlock()

	//Do NOT use default flag.CommandLine because it will cause exit on error.
	//See flag.go.
	flagSet := flag.NewFlagSet(os.Args[0], flag.ContinueOnError)
	var err error

	for k, v := range self.items {
		//Handle command line parameters.
		if k == CmdConfName {
			reflectVal := reflect.ValueOf(v).Elem()
			reflectType := reflectVal.Type()
			var vals []interface{}

			for i := 0; i < reflectVal.NumField(); i++ {
				switch reflectVal.Field(i).Type().Name() {
				case "uint", "uint32":
					vals = append(vals, flagSet.Uint(
						strings.ToLower(reflectType.Field(i).Name),
						0,
						reflectType.Field(i).Tag.Get(UsageTagName)))
				case "uint64":
					vals = append(vals, flagSet.Uint64(
						strings.ToLower(reflectType.Field(i).Name),
						0,
						reflectType.Field(i).Tag.Get(UsageTagName)))
				case "int", "int32":
					vals = append(vals, flagSet.Int(
						strings.ToLower(reflectType.Field(i).Name),
						0,
						reflectType.Field(i).Tag.Get(UsageTagName)))
				case "int64":
					vals = append(vals, flagSet.Int64(
						strings.ToLower(reflectType.Field(i).Name),
						0,
						reflectType.Field(i).Tag.Get(UsageTagName)))
				case "float64":
					vals = append(vals, flagSet.Float64(
						strings.ToLower(reflectType.Field(i).Name),
						0.0,
						reflectType.Field(i).Tag.Get(UsageTagName)))
				case "string":
					vals = append(vals, flagSet.String(
						strings.ToLower(reflectType.Field(i).Name),
						"",
						reflectType.Field(i).Tag.Get(UsageTagName)))
				case "bool":
					vals = append(vals, flagSet.Bool(
						strings.ToLower(reflectType.Field(i).Name),
						false,
						reflectType.Field(i).Tag.Get(UsageTagName)))
				default:
					vals = append(vals, nil)
				}
			}

			err = flagSet.Parse(os.Args[1:])
			if err != nil {
				continue
			}

			for j := 0; j < reflectVal.NumField(); j++ {
				switch reflectVal.Field(j).Type().Name() {
				case "uint", "uint32":
					val1 := *vals[j].(*uint)
					reflectVal.Field(j).SetUint(uint64(val1))
				case "uint64":
					reflectVal.Field(j).SetUint(*vals[j].(*uint64))
				case "int", "int32":
					val2 := *vals[j].(*int)
					reflectVal.Field(j).SetInt(int64(val2))
				case "int64":
					reflectVal.Field(j).SetInt(*vals[j].(*int64))
				case "float64":
					reflectVal.Field(j).SetFloat(*vals[j].(*float64))
				case "string":
					reflectVal.Field(j).SetString(*vals[j].(*string))
				case "bool":
					reflectVal.Field(j).SetBool(*vals[j].(*bool))
				}
			}
		} else { //Handle xml configure files.
			path := fmt.Sprintf("%s%s.conf.%s.xml", self.dir, k, self.env)
			file, err := os.OpenFile(path, os.O_RDONLY, 0666)
			if err != nil {
				continue
			}

			content := make([]byte, ConfFileMaxSize)

			_, err = file.Read(content)
			if err != nil {
				continue
			}

			err = xml.Unmarshal(content, v)
			if err != nil {
				continue
			}
		}
	}

	return err
}
