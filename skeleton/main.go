package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/boringding/beekeeper"
	"github.com/boringding/beekeeper/conf"
	"github.com/boringding/beekeeper/grace"
	"github.com/boringding/beekeeper/proc"
	_ "github.com/boringding/beekeeper/skeleton/handlers"
)

var (
	ENV      = "dev"
	CONF_DIR = "../conf/"
)

func main() {
	err := proc.DumpSelfPid(fmt.Sprintf("%s.pid", os.Args[0]))
	if err != nil {
		fmt.Println("dump self pid failed: ", err)
		return
	}

	beekeeper.InitConf(ENV, CONF_DIR)

	var frameworkConf conf.FrameworkConf

	beekeeper.AddConfItem("framework", &frameworkConf)

	err = beekeeper.ParseConf()
	if err != nil {
		fmt.Println("parse configure failed: ", err)
		return
	}

	err = beekeeper.InitLog(frameworkConf.LogConf)
	if err != nil {
		fmt.Println("initialize log failed: ", err)
		return
	}

	srv, err := grace.NewGracefulSrv(frameworkConf.SrvConf)
	if err != nil {
		beekeeper.LogFatal("create graceful server failed: %v", err)
		return
	}

	err = beekeeper.InitRouter("/" + filepath.Base(os.Args[0]))
	if err != nil {
		beekeeper.LogFatal("initialize router failed: %v", err)
		return
	}

	err = beekeeper.InitMonitor(frameworkConf.MonConf)
	if err != nil {
		beekeeper.LogFatal("initialize monitor failed: %v", err)
		return
	}

	beekeeper.LogInfo("server starting...")

	err = srv.Serve(grace.SrvTypeFcgi, beekeeper.GetRouter())
	if err != nil {
		beekeeper.LogInfo("server finished: %v", err)
		return
	}
}
