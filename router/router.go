package router

import (
	"context"
	"net/http"
	"strings"
)

type ContextKey string

type Router struct {
	routes map[string]map[string]http.HandlerFunc
}

func NewRouter() *Router {
	return &Router{
		routes: make(map[string]map[string]http.HandlerFunc),
	}
}

func (r *Router) addRoute(method, path string, handler http.HandlerFunc) {
	if r.routes[path] == nil {
		r.routes[path] = make(map[string]http.HandlerFunc)
	}

	r.routes[path][method] = handler
}

func (r *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	path := req.URL.Path
	method := req.Method

	if methods, ok := r.routes[path]; ok {
		if handler, ok := methods[method]; ok {
			handler(w, req)
			return
		}
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	for pattern, methods := range r.routes {
		if match, params := matchDynamic(pattern, path); match {
			ctx := req.Context()
			for k, v := range params {
				ctx = context.WithValue(ctx, ContextKey(k), v)
			}
			req = req.WithContext(ctx)

			if handler, ok := methods[method]; ok {
				handler(w, req)
				return
			}
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}
	}

	http.NotFound(w, req)
}

func matchDynamic(pattern, path string) (bool, map[string]string) {
	patternParts := strings.Split(strings.Trim(pattern, "/"), "/")
	pathParts := strings.Split(strings.Trim(path, "/"), "/")

	if len(patternParts) != len(pathParts) {
		return false, nil
	}

	params := make(map[string]string)
	for i, pp := range patternParts {
		if strings.HasPrefix(pp, "{") && strings.HasSuffix(pp, "}") {
			// dynamic parameter
			key := strings.TrimSuffix(strings.TrimPrefix(pp, "{"), "}")
			params[key] = pathParts[i]
		} else if pp != pathParts[i] {
			return false, nil
		}
	}

	return true, params
}

func (r *Router) GET(path string, handler http.HandlerFunc) {
	r.addRoute(http.MethodGet, path, handler)
}

func (r *Router) POST(path string, handler http.HandlerFunc) {
	r.addRoute(http.MethodPost, path, handler)
}

func (r *Router) PUT(path string, handler http.HandlerFunc) {
	r.addRoute(http.MethodPut, path, handler)
}

func (r *Router) DELETE(path string, handler http.HandlerFunc) {
	r.addRoute(http.MethodDelete, path, handler)
}
