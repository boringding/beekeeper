package main

import (
	"fmt"
	"os"

	"github.com/boringding/beekeeper"
	"github.com/boringding/beekeeper/conf"
	"github.com/boringding/beekeeper/grace"
	"github.com/boringding/beekeeper/proc"
	//_ "github.com/boringding/beekeeper/skeleton/handlers"
)

type CmdConf struct {
	A uint    `usage:"parameter a"`
	B uint32  `usage:"parameter b"`
	C string  `usage:"parameter c"`
	D float64 `usage:"parameter d"`
	E bool    `usage:"parameter e"`
	F int64   `usage:"parameter f"`
	G uint64  `usage:"parameter g"`
	H int32   `usage:"parameter h"`
	I int     `usage:"parameter i"`
}

func main() {
	err := proc.DumpSelfPid("./beekeeper.pid")
	if err != nil {
		fmt.Println(err)
		return
	}

	beekeeper.InitConf(os.Args[1], "../conf/")

	var cmdConf CmdConf
	var frameworkConf conf.FrameworkConf

	beekeeper.AddCmdConfItem(&cmdConf)
	beekeeper.AddConfItem("framework", &frameworkConf)

	err = beekeeper.ParseConf()
	if err != nil {
		fmt.Println(err)
		return
	}

	err = beekeeper.InitLog(frameworkConf.LogConf)
	if err != nil {
		fmt.Println(err)
		return
	}

	beekeeper.LogInfo("log init finished")

	srv, err := grace.NewGracefulSrv(frameworkConf.SrvConf)
	if err != nil {
		beekeeper.LogFatal("create graceful server failed")
		return
	}

	beekeeper.LogInfo("server starting...")

	err = srv.Serve(grace.SrvTypeFcgi, beekeeper.GetRouter())
	if err != nil {
		beekeeper.LogFatal("start serve failed")
		return
	}

	beekeeper.LogInfo("server finished")
}
