package router

import (
	"fmt"
	"reflect"
)

type Handler struct {
}

func (self *Handler) Register(concreteHandler interface{}) {
	fmt.Println(reflect.TypeOf(concreteHandler))
}
