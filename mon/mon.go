package mon

import (
	"fmt"
	"net"
	"net/http"

	"github.com/boringding/beekeeper/conf"
)

type Monitor struct {
	host string
	port int
}

func NewMonitor() *Monitor {
	return &Monitor{
		host: "",
		port: 0,
	}
}

func (self *Monitor) Init(conf conf.MonConf) error {
	if conf.Enabled == false {
		return nil
	}

	self.host = conf.Host
	self.port = conf.Port

	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%d", self.host, self.port))
	if err != nil {
		return err
	}

	go func() {
		http.Serve(listener, nil)
	}()

	return nil
}
