package handlers

import (
	"net/http"

	"github.com/boringding/beekeeper"
)

type Foo struct {
}

func init() {
	beekeeper.AddRoute("/foo/c/d", "GET", &Foo{})
}

func (self *Foo) ServeHTTP(resWriter http.ResponseWriter, req *http.Request) {
	str := req.URL.Path
	resWriter.Write([]byte(str))
}
