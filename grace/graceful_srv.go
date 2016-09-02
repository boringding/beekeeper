package grace

import (
	"errors"
	"fmt"
	"net/http"
	"net/http/fcgi"
	"os"
	"os/exec"
	"os/signal"
	"runtime"
	"strconv"
	"sync"
	"syscall"
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
	SrvTypeHttp = iota
	SrvTypeFcgi
)

const BeekeeperChildEnv = "BEEKEEPER_CHILD"

var (
	gracefulSrvs = map[string]*GracefulSrv{}
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
	if v, ok := gracefulSrvs[conf.Name]; ok {
		return v, errors.New("server already exists")
	}

	var gracefulSrv GracefulSrv
	err := gracefulSrv.init(conf)
	if err != nil {
		return nil, err
	}

	gracefulSrvs[conf.Name] = &gracefulSrv

	return &gracefulSrv, err
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

func handleSignal() {
	var sig os.Signal

	signal.Notify(sigChan, syscall.SIGHUP, syscall.SIGTERM)

	for {
		sig = <-sigChan
		switch sig {
		case syscall.SIGHUP:
			sighupHandler()
		case syscall.SIGTERM:
			sigtermHandler()
		}
	}

	return
}

func sighupHandler() error {
	if isForked == true {
		return nil
	}

	isForked = true

	var files []*os.File

	i := 0
	for k, v := range gracefulSrvs {
		file, err := v.listener.file()
		if err != nil {
			continue
		}

		files = append(files, file)
		envVal := 3 + i
		envValStr := strconv.Itoa(envVal)
		i++

		err = proc.SetEnv(k, envValStr)
		if err != nil {
			continue
		}
	}

	err := proc.SetEnv(BeekeeperChildEnv, "1")
	if err != nil {
		return err
	}

	exeFilePath := os.Args[0]
	var args []string
	if len(os.Args) > 1 {
		args = os.Args[1:]
	}

	cmd := exec.Command(exeFilePath, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.ExtraFiles = files
	cmd.Env = os.Environ()

	return cmd.Start()
}

func sigtermHandler() error {
	if isClosed == true {
		return nil
	}

	isClosed = true

	var err error
	for _, v := range gracefulSrvs {
		e := v.shutdown()
		if e != nil {
			err = e
		}
	}

	return err
}
