//Example handler.

package handlers

import (
	"context"
	"net/http"

	"github.com/boringding/beekeeper"
)

func init() {
	beekeeper.AddRoute("/example/a/b", "GET,POST", 0, Example)
}

func Example(ctx context.Context, resWriter http.ResponseWriter, req *http.Request) {
	beekeeper.LogInfo("handle example request")
	beekeeper.MetricsAdd("EXAMPLE_REQ", 1)

	str := req.URL.Path
	resWriter.Write([]byte(str))
}
