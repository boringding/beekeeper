//Example command line configure struct definition
//and registration.

package main

import (
	"github.com/boringding/beekeeper"
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

var CommandConf CmdConf

func init() {
	beekeeper.AddCmdConfItem(&CommandConf)
}
