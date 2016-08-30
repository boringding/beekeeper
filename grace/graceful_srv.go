package grace

import (
	"fmt"
	"github.com/boringding/beekeeper/conf"
	"net/http"
	"net/http/fcgi"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

type GracefulSrv struct {
	http.Server
	listener  *gracefulListener
	waitGroup sync.WaitGroup
	sigChan   chan os.Signal
}

func (self *GracefulSrv) Init(conf conf.SrvConf) error {
	addr := fmt.Sprintf("%s:%d", conf.Host, conf.Port)

	self.listener = &gracefulListener{
		name:             conf.Name,
		addr:             addr,
		keepAlive:        conf.KeepAlive,
		keepAliveSeconds: conf.KeepAliveSeconds,
		srv:              self,
	}

	self.Server.Addr = addr
	self.Server.ReadTimeout = time.Duration(conf.ReadTimeoutSeconds) * time.Second
	self.Server.WriteTimeout = time.Duration(conf.WriteTimeoutSeconds) * time.Second
	self.Server.MaxHeaderBytes = conf.MaxHeaderBytes

	self.sigChan = make(chan os.Signal)

	return self.listener.Init()
}

func (self *GracefulSrv) ServeHttp(handler http.Handler) error {
	self.Server.Handler = handler
	err := self.Server.Serve(self.listener)
	self.waitGroup.Wait()

	return err
}

func (self *GracefulSrv) ServeFcgi(handler http.Handler) error {
	err := fcgi.Serve(self.listener, handler)
	self.waitGroup.Wait()

	return err
}

func (self *GracefulSrv) ShutdownHttp() error {
	self.SetKeepAlivesEnabled(false)
	return self.listener.Close()
}

func (self *GracefulSrv) handleSignal() {
	var sig os.Signal

	signal.Notify(self.sigChan, syscall.SIGHUP, syscall.SIGTERM)

	for {
		sig = <-self.sigChan
		switch sig {

		}
	}
}
