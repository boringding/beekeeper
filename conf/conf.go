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
)

type Conf struct {
	mu    sync.RWMutex
	items map[string]interface{}
	env   string
	dir   string
}

func NewConf() *Conf {
	return &Conf{
		items: map[string]interface{}{},
		env:   "dev",
		dir:   "../conf/",
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

	flagSet := flag.NewFlagSet(os.Args[0], flag.ContinueOnError)
	var err error

	for k, v := range self.items {
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
				}
			}

			//parsing from the 2nd parameter, os.Args[1] is used for environment
			err = flagSet.Parse(os.Args[2:])
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
		} else {
			path := fmt.Sprintf("%s%s.conf.%s.xml", self.dir, k, self.env)
			file, err := os.OpenFile(path, os.O_RDONLY, 0666)
			if err != nil {
				continue
			}

			content := []byte{}

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
