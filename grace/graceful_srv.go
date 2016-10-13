//Type GracefulSrv inherits http.Server.
//It implements a server which can restart smoothly
//without interrupting the service.

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
	listener  *gracefulListener
	waitGroup sync.WaitGroup
	srvType   int
	//The longest time for the old server to wait for all
	//established connections to close.
	shutdownTimeoutSeconds int
}

const (
	//The environment variable name which indicates
	//whether the process is a child process.
	BeekeeperChildEnv = "BEEKEEPER_CHILD"
	SrvTypeHttp       = iota
	SrvTypeFcgi
)

var (
	mu           sync.Mutex
	GracefulSrvs = map[string]*GracefulSrv{}
	sigChan      chan os.Signal
	//The flag which indicates whether the process has forked
	//a new process to accept connections(restart).
	isForked bool
	//The flag which indicates whether the process is about to close.
	isClosed bool
)

func init() {
	sigChan = make(chan os.Signal)
	isForked = false
	isClosed = false

	//If this is a new-forked child process
	//it should tell its parent process to exit.
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
	//When the listener is closed and the server jumps out of
	//the loop it should wait the established connections to
	//finish their process.
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
		//Before closing the listener of a http server
		//keep-alive option should be disabled.
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
		//When there is no go routines to wait calling
		//sync.WaitGroup.Done method will cause a panic.
		self.waitGroup.Done()
		//Yields the processor so that the server routine
		//may have a chance to return as soon as possible.
		runtime.Gosched()
	}

	return
}
