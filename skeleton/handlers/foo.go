package handlers

import (
	"net/http"

	"github.com/boringding/beekeeper/router"
)

type Foo struct {
}

func init() {
	var foo Foo
	router.DefaultRouter.AddRoute(router.Route{
		Handler: &foo,
		Method:  router.MethodPost,
		Path:    "/foo/c/d",
	})
}

func (self *Foo) ServeHTTP(resWriter http.ResponseWriter, req *http.Request) {
	str := req.URL.Path
	resWriter.Write([]byte(str))
}
