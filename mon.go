package beekeeper

import (
	"github.com/boringding/beekeeper/conf"
	"github.com/boringding/beekeeper/mon"
)

var defaultMonitor = mon.NewMonitor()

func InitMonitor(conf conf.MonConf) error {
	return defaultMonitor.Init(conf)
}
