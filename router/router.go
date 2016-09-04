package router

import (
	"errors"
	"net/http"
	"sync"
)

const (
	MethodGet = 1 << iota
	MethodHead
	MethodPost
	MethodPut
	MethodPatch
	MethodDelete
	MethodConnect
	MethodOptions
	MethodTrace
)

var MethodMap = map[string]int{
	http.MethodGet:     MethodGet,
	http.MethodHead:    MethodHead,
	http.MethodPost:    MethodPost,
	http.MethodPut:     MethodPut,
	http.MethodPatch:   MethodPatch,
	http.MethodDelete:  MethodDelete,
	http.MethodConnect: MethodConnect,
	http.MethodOptions: MethodOptions,
	http.MethodTrace:   MethodTrace,
}

type Route struct {
	http.Handler
	Method int
	Path   string
}

type router struct {
	mu     sync.RWMutex
	routes map[string]map[int]Route
}

func newRouter() *router {
	return &router{
		routes: map[string]map[int]Route{},
	}
}

var DefaultRouter = newRouter()

func (self *router) AddRoute(route Route) error {
	self.mu.Lock()
	defer self.mu.Unlock()

	pathOk := false
	if _, pathOk = self.routes[route.Path]; pathOk {
		if _, methodOk := self.routes[route.Path][route.Method]; methodOk {
			return errors.New("route already exists")
		}
	}

	if pathOk == false {
		self.routes[route.Path] = map[int]Route{}
	}

	self.routes[route.Path][route.Method] = route

	return nil
}

func (self *router) FindRoute(method int, path string) (Route, bool) {
	self.mu.RLock()
	defer self.mu.RUnlock()

	var route Route

	if _, ok := self.routes[path]; ok {
		for _, v := range self.routes[path] {
			if v.Method&method > 0 {
				return v, true
			}
		}
	}

	return route, false
}

func (self *router) ServeHTTP(resWriter http.ResponseWriter, req *http.Request) {
	if v, ok := self.FindRoute(MethodMap[req.Method], req.URL.Path); ok {
		v.Handler.ServeHTTP(resWriter, req)
	} else {
		http.NotFound(resWriter, req)
	}
}
