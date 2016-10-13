//Type gracefulListener inherits net.Listener and
//overrides Accept and Close method.

package grace

import (
	"net"
	"os"
	"strconv"
	"time"

	"github.com/boringding/beekeeper/proc"
)

type gracefulListener struct {
	net.Listener
	name             string
	addr             string
	keepAlive        bool
	keepAliveSeconds int64
	srv              *GracefulSrv
}

func (self *gracefulListener) init() error {
	fdStr := proc.GetEnv(self.name)
	var err error

	//If the environment variable is not set
	//create the listener with the address directly,
	if len(fdStr) <= 0 {
		self.Listener, err = net.Listen("tcp", self.addr)
		if err != nil {
			return err
		}
	} else { //or create the listener with the file descriptor set in the environment variable.
		fd, err := strconv.Atoi(fdStr)
		if err != nil {
			return err
		}

		file := os.NewFile(uintptr(fd), "")
		self.Listener, err = net.FileListener(file)
		if err != nil {
			return err
		}
	}

	return nil
}

func (self *gracefulListener) Accept() (net.Conn, error) {
	tcpConn, err := self.Listener.(*net.TCPListener).AcceptTCP()
	if err != nil {
		return tcpConn, err
	}

	tcpConn.SetKeepAlive(self.keepAlive)
	tcpConn.SetKeepAlivePeriod(time.Duration(self.keepAliveSeconds) * time.Second)

	gracefulConn := gracefulConn{
		Conn: tcpConn,
		srv:  self.srv,
	}

	//After accept a connection successfully
	//increase the server's sync.WaitGroup.
	self.srv.waitGroup.Add(1)

	return &gracefulConn, nil
}

func (self *gracefulListener) Close() error {
	return self.Listener.Close()
}

func (self *gracefulListener) file() (*os.File, error) {
	return self.Listener.(*net.TCPListener).File()
}
