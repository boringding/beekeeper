package handlers

import (
	"net/http"

	"github.com/boringding/beekeeper"
)

type Example struct {
}

func init() {
	beekeeper.AddRoute("/example/a/b", "GET,POST", &Example{})
}

func (self *Example) ServeHTTP(resWriter http.ResponseWriter, req *http.Request) {
	str := req.URL.Path
	resWriter.Write([]byte(str))
}
