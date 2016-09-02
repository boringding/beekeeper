package main

import (
	"fmt"
	"net/http"
	"reflect"
	_ "time"

	"github.com/boringding/beekeeper/conf"
	_ "github.com/boringding/beekeeper/grace"
	"github.com/boringding/beekeeper/router"
)

func init() {
	var handlerA HandlerA
	handlerA.Register()
}

func main() {
	srvConf := conf.SrvConf{
		Name:                   "FUCK",
		Host:                   "127.0.0.1",
		Port:                   8900,
		KeepAlive:              true,
		KeepAliveSeconds:       60,
		ReadTimeoutSeconds:     30,
		WriteTimeoutSeconds:    30,
		MaxHeaderBytes:         1000,
		ShutdownTimeoutSeconds: 30,
	}

	fmt.Println(reflect.TypeOf(srvConf).Name())
	fmt.Println(reflect.TypeOf(srvConf).PkgPath())

	/*
		srv, err := grace.NewGracefulSrv(srvConf)
		if err != nil {
			fmt.Println(err)
			return
		}

		var defaultHandler DefaultHandler

		err = srv.Serve(grace.SrvTypeHttp, &defaultHandler)
		if err != nil {
			fmt.Println(err)
		}
	*/
}

type HandlerA struct {
	handler router.Handler
}

func (self *HandlerA) Register() {
	self.handler.Register(self)
}

type DefaultHandler struct {
}

func (self *DefaultHandler) ServeHTTP(resWriter http.ResponseWriter, req *http.Request) {
	//time.Sleep(10 * time.Second)
	str := req.URL.Path
	resWriter.Write([]byte(str))
}
