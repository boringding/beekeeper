package beekeeper

import (
	"github.com/boringding/beekeeper/router"
)

var defaultRouter = router.NewRouter()

func AddRoute(route router.Route) error {
	return defaultRouter.AddRoute(route)
}

func FindRoute(method int, path string) (router.Route, bool) {
	return defaultRouter.FindRoute(method, path)
}

func GetRouter() *router.Router {
	return defaultRouter
}
