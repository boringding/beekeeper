//Monitor interface.

package beekeeper

import (
	"github.com/boringding/beekeeper/conf"
	"github.com/boringding/beekeeper/mon"
)

var defaultMonitor = mon.NewMonitor()
var customMetrics = map[string]*mon.Metrics{}

func InitMonitor(conf conf.MonConf) error {
	return defaultMonitor.Init(conf)
}

func CloseMonitor() error {
	return defaultMonitor.Close()
}

func MetricsAdd(name string, delta int64) error {
	if v, ok := customMetrics[name]; ok {
		return v.Add(delta)
	}

	err, m := mon.NewMetrics(name, delta)
	if err != nil {
		return err
	}

	customMetrics[name] = m

	return nil
}

func MetricsSet(name string, val int64) error {
	if v, ok := customMetrics[name]; ok {
		return v.Set(val)
	}

	err, m := mon.NewMetrics(name, val)
	if err != nil {
		return err
	}

	customMetrics[name] = m

	return nil
}
