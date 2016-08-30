package proc

import (
	"fmt"
	"testing"
)

func Test_Pid(t *testing.T) {
	fmt.Println(GetSelfPid())
	fmt.Println(GetParentPid())
	fmt.Println(GetSelfPid())
	fmt.Println(GetParentPid())
}
