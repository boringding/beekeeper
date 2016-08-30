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
	self.srv.waitGroup.Done()

	return err
}
