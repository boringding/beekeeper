package beekeeper

import (
	"strings"

	"github.com/boringding/beekeeper/router"
)

const MethodSeperator = ","

var defaultRouter = router.NewRouter()

func InitRouter(pathPrefix string) error {
	return defaultRouter.Init(pathPrefix)
}

func PathPrefix() string {
	return defaultRouter.PathPrefix()
}

func SetPathPrefix(pathPrefix string) {
	defaultRouter.SetPathPrefix(pathPrefix)
}

func AddRoute(path string, method string, timeoutSeconds int, handle router.HandleFunc) error {
	methods := strings.Split(method, MethodSeperator)

	m := 0
	for _, v1 := range methods {
		if v2, ok := router.MethodMap[v1]; ok {
			m = m | v2
		}
	}

	return defaultRouter.AddRoute(router.Route{
		Handle:         handle,
		Method:         m,
		Path:           path,
		TimeoutSeconds: timeoutSeconds,
	})
}

func FindRoute(method int, path string) (router.Route, bool) {
	return defaultRouter.FindRoute(method, path)
}

func GetRouter() *router.Router {
	return defaultRouter
}
