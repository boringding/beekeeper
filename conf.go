package beekeeper

import (
	"github.com/boringding/beekeeper/conf"
)

var defaultConf = conf.NewConf()

func InitConf(env string, dir string) {
	defaultConf.Init(env, dir)
}

func SetConfEnv(env string) {
	defaultConf.SetEnv(env)
}

func SetConfDir(dir string) {
	defaultConf.SetDir(dir)
}

func AddConfItem(name string, v interface{}) error {
	return defaultConf.AddItem(name, v)
}

func AddCmdConfItem(v interface{}) error {
	return defaultConf.AddItem(conf.CmdConfName, v)
}

func ParseConf() error {
	return defaultConf.Parse()
}
