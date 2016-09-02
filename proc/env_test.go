package proc

import (
	"fmt"
	"testing"
)

func Test_GetEnv(t *testing.T) {
	SetEnv("FFF", "t6")
	fmt.Println(GetEnv("FFF"))
	SetEnv("GGG", "t7")
	fmt.Println(GetEnv("GGG"))
}
