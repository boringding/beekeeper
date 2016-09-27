package handlers

import (
	"context"
	"net/http"

	"github.com/boringding/beekeeper"
)

func init() {
	beekeeper.AddRoute("/foo/c/d", "GET", 0, Foo)
}

func Foo(ctx context.Context, resWriter http.ResponseWriter, req *http.Request) {
	str := req.URL.Path
	resWriter.Write([]byte(str))
}
