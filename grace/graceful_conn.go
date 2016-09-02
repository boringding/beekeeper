package grace

import (
	"net"
)

type gracefulConn struct {
	net.Conn
	srv *GracefulSrv
}

func (self *gracefulConn) Close() error {
	err := self.Conn.Close()
	if err == nil {
		self.srv.waitGroup.Done()
	}

	return err
}
