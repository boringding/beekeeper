package beekeeper

import (
	"fmt"
	"testing"
)

type TestParam struct {
	A string  `param:"a"`
	B string  `param:"b"`
	C string  `param:"c"`
	D int     `param:"d"`
	E bool    `param:"e"`
	F float32 `param:"f"`
	G uint64  `param:"g"`
}

func Test_ParseQueryStr(t *testing.T) {
	var testParam TestParam
	err := ParseQueryStr("a=this&b=is&c=test&d=-35&e=true&f=4.55a&g=23423", &testParam)
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(testParam)
	}
}
