package beekeeper

import (
	"fmt"
	"testing"
)

type TestParam struct {
	A string `param:"a"`
	B string `param:"b"`
	C string `param:"c"`
}

func Test_ParseQueryStr(t *testing.T) {
	var testParam TestParam
	err := ParseQueryStr("a=this&b=is&c=test", &testParam)
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(testParam)
	}
}
