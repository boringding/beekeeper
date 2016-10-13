//Type Monitor provides a way to create, start and stop
//http.DefaultServeMux which gives access to process
//runtime statistics (see expvar.go)
//and metrics.

package mon

import (
	"fmt"
	"net"
	"net/http"

	"github.com/boringding/beekeeper/conf"
)

type Monitor struct {
	host     string
	port     int
	listener net.Listener
}

func NewMonitor() *Monitor {
	return &Monitor{
		host: "",
		port: 0,
	}
}

func (self *Monitor) Init(conf conf.MonConf) (err error) {
	if conf.Enabled == false {
		return
	}

	self.host = conf.Host
	self.port = conf.Port

	self.listener, err = net.Listen("tcp", fmt.Sprintf("%s:%d", self.host, self.port))
	if err != nil {
		return
	}

	go func() {
		http.Serve(self.listener, nil)
	}()

	return
}

func (self *Monitor) Close() error {
	return self.listener.Close()
}
