//Type Router distributes requests according to
//their path.

package router

import (
	"context"
	"errors"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/boringding/beekeeper/mon"
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

type HandleFunc func(ctx context.Context, resWriter http.ResponseWriter, req *http.Request)

type Route struct {
	Handle         HandleFunc
	Method         int
	Path           string
	TimeoutSeconds int
}

type Router struct {
	mu         sync.RWMutex
	routes     map[string]map[int]Route
	pathPrefix string
	//The number of requests received by the Router.
	totalReq *mon.Metrics
	//Total elapsed time on handling requests in millisecond.
	totalTime *mon.Metrics
}

func NewRouter() *Router {
	return &Router{
		routes:     map[string]map[int]Route{},
		pathPrefix: "",
		totalReq:   nil,
		totalTime:  nil,
	}
}

func (self *Router) Init(pathPrefix string) error {
	self.pathPrefix = pathPrefix

	err, totalReq := mon.NewMetrics("TOTAL_REQ", int64(0))
	err, totalTime := mon.NewMetrics("TOTAL_TIME", int64(0))
	if err != nil {
		return err
	}

	self.totalReq = totalReq
	self.totalTime = totalTime

	return nil
}

func (self *Router) PathPrefix() string {
	self.mu.RLock()
	defer self.mu.RUnlock()

	return self.pathPrefix
}

func (self *Router) SetPathPrefix(pathPrefix string) {
	self.mu.Lock()
	defer self.mu.Unlock()

	self.pathPrefix = pathPrefix
}

//The route.Path must NOT include prefix.
func (self *Router) AddRoute(route Route) error {
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

//The parameter path should include prefix.
func (self *Router) FindRoute(method int, path string) (Route, bool) {
	self.mu.RLock()
	defer self.mu.RUnlock()

	path = strings.TrimPrefix(path, self.pathPrefix)

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

//Type Router implements ServeHTTP method
//so it is an implementation of http.Handler.
//See server.go.
func (self *Router) ServeHTTP(resWriter http.ResponseWriter, req *http.Request) {
	self.totalReq.Add(int64(1))

	if v, ok := self.FindRoute(MethodMap[req.Method], req.URL.Path); ok {
		var ctx context.Context
		var cancel context.CancelFunc

		if v.TimeoutSeconds > 0 {
			ctx, cancel = context.WithTimeout(req.Context(), time.Duration(v.TimeoutSeconds)*time.Second)
		} else {
			ctx, cancel = context.WithCancel(req.Context())
		}

		defer cancel()

		start := time.Now()
		v.Handle(ctx, resWriter, req)
		duration := time.Since(start)

		self.totalTime.Add(int64(duration / time.Millisecond))

	} else {
		http.NotFound(resWriter, req)
	}
}
