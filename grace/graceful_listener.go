package grace

import (
	"github.com/boringding/beekeeper/proc"
	"net"
	"os"
	"strconv"
	"time"
)

type gracefulListener struct {
	net.Listener
	name             string
	addr             string
	keepAlive        bool
	keepAliveSeconds int64
	srv              *GracefulSrv
}

func (self *gracefulListener) Init() error {
	fdStr := proc.GetEnv(self.name)
	var err error

	if len(fdStr) <= 0 {
		self.Listener, err = net.Listen("tcp", self.addr)
		if err != nil {
			return err
		}

		file, err := self.Listener.(*net.TCPListener).File()
		if err != nil {
			return err
		}

		fdStr = strconv.Itoa(int(file.Fd()))
		err = proc.SetEnv(self.name, fdStr)
		if err != nil {
			return err
		}
	} else {
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

	self.srv.waitGroup.Add(1)

	return &gracefulConn, nil
}

func (self *gracefulListener) Close() error {
	err := self.Listener.Close()
	err = proc.SetEnv(self.name, "")

	return err
}
