package server

import (
	"context"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

func New(ctx context.Context, ip string, port int) (*Server, error) {
	srv := &Server{
		router: mux.NewRouter().StrictSlash(true),
	}

	srv.api = srv.newRouter("/api")
	srv.sites = srv.newRouter("/public")

	return srv, nil
}

type Server struct {
	router *mux.Router

	api   *Router
	sites *Router
}

func (s *Server) API() *Router {
	return s.api
}

func (s *Server) Sites() *Router {
	return s.sites
}

func (s *Server) Routes() {
	s.router.Walk(walkFn)
}

func walkFn(route *mux.Route, router *mux.Router, ancestors []*mux.Route) error {
	path, err := route.GetPathTemplate()
	if err != nil {
		return err
	}
	fmt.Println(path)
	return nil
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.router.ServeHTTP(w, r)
}

func (s *Server) newRouter(base string) *Router {
	mux := s.router.PathPrefix(base).Subrouter()
	return &Router{mux: mux}
}

type Router struct {
	mux *mux.Router
}

func (r *Router) Mount(mountPoint string, routes []Route) {
	mux := r.mux.PathPrefix(mountPoint).Subrouter()
	for _, route := range routes {
		mux.HandleFunc(route.Path, route.Func).Methods(route.Methods...)
	}
}

type Route struct {
	Func    http.HandlerFunc
	Path    string
	Methods []string
}
