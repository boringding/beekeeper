package grace

import (
	"errors"
	"fmt"
	"net/http"
	"net/http/fcgi"
	"os"
	"runtime"
	"sync"
	"time"

	"github.com/boringding/beekeeper/conf"
	"github.com/boringding/beekeeper/proc"
)

type GracefulSrv struct {
	http.Server
	listener               *gracefulListener
	waitGroup              sync.WaitGroup
	srvType                int
	shutdownTimeoutSeconds int
}

const (
	BeekeeperChildEnv = "BEEKEEPER_CHILD"
	SrvTypeHttp       = iota
	SrvTypeFcgi
)

var (
	mu           sync.Mutex
	GracefulSrvs = map[string]*GracefulSrv{}
	sigChan      chan os.Signal
	isForked     bool
	isClosed     bool
)

func init() {
	sigChan = make(chan os.Signal)
	isForked = false
	isClosed = false

	if len(proc.GetEnv(BeekeeperChildEnv)) > 0 {
		proc.TerminateProc(proc.GetParentPid())
	}

	go handleSignal()
}

func NewGracefulSrv(conf conf.SrvConf) (*GracefulSrv, error) {
	mu.Lock()
	defer mu.Unlock()

	if v, ok := GracefulSrvs[conf.Name]; ok {
		return v, errors.New("server already exists")
	}

	var srv GracefulSrv
	err := srv.init(conf)
	if err != nil {
		return nil, err
	}

	GracefulSrvs[conf.Name] = &srv

	return &srv, err
}

func (self *GracefulSrv) init(conf conf.SrvConf) error {
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

	self.shutdownTimeoutSeconds = conf.ShutdownTimeoutSeconds

	return self.listener.init()
}

func (self *GracefulSrv) serveHttp(handler http.Handler) error {
	self.Server.Handler = handler
	err := self.Server.Serve(self.listener)
	self.waitGroup.Wait()

	return err
}

func (self *GracefulSrv) serveFcgi(handler http.Handler) error {
	err := fcgi.Serve(self.listener, handler)
	self.waitGroup.Wait()

	return err
}

func (self *GracefulSrv) Serve(srvType int, handler http.Handler) error {
	self.srvType = srvType

	switch srvType {
	case SrvTypeHttp:
		return self.serveHttp(handler)
	case SrvTypeFcgi:
		return self.serveFcgi(handler)
	}

	return nil
}

func (self *GracefulSrv) shutdown() error {
	go self.stopConns()

	switch self.srvType {
	case SrvTypeHttp:
		self.SetKeepAlivesEnabled(false)
		return self.listener.Close()
	case SrvTypeFcgi:
		return self.listener.Close()
	}

	return nil
}

func (self *GracefulSrv) stopConns() {
	defer func() {
		if r := recover(); r != nil {
			return
		}
	}()

	time.Sleep(time.Duration(self.shutdownTimeoutSeconds) * time.Second)

	for {
		self.waitGroup.Done()
		runtime.Gosched()
	}

	return
}
