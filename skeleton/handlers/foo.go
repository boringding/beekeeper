//Example handler.

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
	beekeeper.LogInfo("handle foo request")
	beekeeper.MetricsAdd("FOO_REQ", 1)

	str := req.URL.Path
	resWriter.Write([]byte(str))
}
