package proc

import (
	"fmt"
	"testing"
)

func Test_GetEnv(t *testing.T) {
	SetEnv("A", "12")
	fmt.Println(GetEnv("A"))
	SetEnv("B", "34")
	fmt.Println(GetEnv("B"))
}
