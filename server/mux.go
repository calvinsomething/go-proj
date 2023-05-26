package main

import (
	"net/http"
)

type (
	methodRouter struct {
		get    http.HandlerFunc
		put    http.HandlerFunc
		post   http.HandlerFunc
		delete http.HandlerFunc
	}

	mux struct {
		*http.ServeMux
		middleware []http.HandlerFunc
		routes     map[string]methodRouter
	}
)

func newMux(middleware ...http.HandlerFunc) *mux {
	return &mux{
		ServeMux:   http.NewServeMux(),
		middleware: middleware,
		routes:     map[string]methodRouter{},
	}
}

var invalidMethod = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusMethodNotAllowed)
})

func (m *mux) post(route string, handler http.HandlerFunc) {
	if r, ok := m.routes[route]; ok {
		r.post = handler
	} else {
		m.routes[route] = methodRouter{post: handler}
	}
}

func (m *mux) get(route string, handler http.HandlerFunc) {
	if r, ok := m.routes[route]; ok {
		r.get = handler
	} else {
		m.routes[route] = methodRouter{get: handler}
	}
}

func (m *mux) put(route string, handler http.HandlerFunc) {
	if r, ok := m.routes[route]; ok {
		r.put = handler
	} else {
		m.routes[route] = methodRouter{put: handler}
	}
}

func (m *mux) delete(route string, handler http.HandlerFunc) {
	if r, ok := m.routes[route]; ok {
		r.delete = handler
	} else {
		m.routes[route] = methodRouter{delete: handler}
	}
}

func (m *mux) ListenAndServe(addr string) error {
	for k, v := range m.routes {
		// all unassigned request methods should return status 405
		for _, h := range []*http.HandlerFunc{&v.get, &v.put, &v.post, &v.delete} {
			if *h == nil {
				*h = invalidMethod
			}
		}
		// assign each route to a methodRouter
		m.HandleFunc(k, routeByMethod(k, v))
	}
	return http.ListenAndServe(addr, m)
}

func (m *mux) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	for _, mw := range m.middleware {
		mw(w, r)
	}
	m.ServeMux.ServeHTTP(w, r)
}

func routeByMethod(route string, router methodRouter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "GET":
			router.get(w, r)
		case "PUT":
			router.put(w, r)
		case "POST":
			router.post(w, r)
		case "DELETE":
			router.delete(w, r)
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	}
}
