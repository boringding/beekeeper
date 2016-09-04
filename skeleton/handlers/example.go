package handlers

import (
	"net/http"

	"github.com/boringding/beekeeper/router"
)

type Example struct {
}

func init() {
	var example Example
	router.DefaultRouter.AddRoute(router.Route{
		Handler: &example,
		Method:  router.MethodGet | router.MethodPost,
		Path:    "/example/a/b",
	})
}

func (self *Example) ServeHTTP(resWriter http.ResponseWriter, req *http.Request) {
	str := req.URL.Path
	resWriter.Write([]byte(str))
}
