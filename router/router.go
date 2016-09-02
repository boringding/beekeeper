package router

import (
	"errors"
	"net/http"
	"sync"
)

type Router struct {
	mu       sync.RWMutex
	mappings map[string]http.Handler
}

var DefaultRouter = &Router{
	mappings: map[string]http.Handler{},
}

func (self *Router) AddMapping(pattern string, handler http.Handler) error {
	self.mu.Lock()
	defer self.mu.Unlock()

	if len(pattern) <= 0 {
		return errors.New("empty pattern")
	}

	if _, ok := self.mappings[pattern]; ok {
		return errors.New("mapping already exists")
	}

	self.mappings[pattern] = handler

	return nil
}

func (self *Router) ServeHTTP(resWriter http.ResponseWriter, req *http.Request) {
	pattern := req.URL.Path

	self.mu.RLock()
	defer self.mu.RUnlock()

	if v, ok := self.mappings[pattern]; ok {
		v.ServeHTTP(resWriter, req)
	} else {

	}
}
