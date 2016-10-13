//Type gracefulConn inherits net.Conn and overrides Close method

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
	//After close the connection successfully
	//reduce the server's sync.WaitGroup.
	if err == nil {
		self.srv.waitGroup.Done()
	}

	return err
}
