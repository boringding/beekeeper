package main

import (
	"fmt"

	"github.com/boringding/beekeeper/conf"
	"github.com/boringding/beekeeper/grace"
	"github.com/boringding/beekeeper/router"
	_ "github.com/boringding/beekeeper/skeleton/handlers"
)

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

	srv, err := grace.NewGracefulSrv(srvConf)
	if err != nil {
		fmt.Println(err)
		return
	}

	err = srv.Serve(grace.SrvTypeHttp, router.DefaultRouter)
	if err != nil {
		fmt.Println(err)
	}
}
