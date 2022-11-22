package server

import (
	"context"
	"fmt"
	"html/template"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

func New(ctx context.Context, ip string, port int) (*Server, error) {
	srv := &Server{
		IP:     ip,
		Port:   port,
		router: mux.NewRouter().StrictSlash(true),
	}

	srv.api = srv.newRouter("/api")
	srv.sites = srv.newRouter("/public")
	srv.sites.mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		markup := `
		<!doctype html>

		<html>
			<head>
				<title>Ygg</title>
			</head>

			<body>
				<ul>
					{{ range . }}
					<li><a href="/public{{ . }}">{{ . }}</a></li>
					{{ end }}
				</ul>
			</body>
		</html>
		`
		tmpl := template.Must(template.New("index").Parse(markup))
		if err := tmpl.Execute(w, srv.sites.services); err != nil {
			panic(err)
		}
	})

	return srv, nil
}

type Server struct {
	IP   string
	Port int

	router *mux.Router

	api   *Router
	sites *Router
}

func (s *Server) Addr() string {
	return fmt.Sprintf("%s:%d", s.IP, s.Port)
}

func (s *Server) Listen() error {
	srv := &http.Server{
		Addr:              s.Addr(),
		Handler:           s.router,
		ReadTimeout:       30 * time.Second,
		ReadHeaderTimeout: 5 * time.Second,
		WriteTimeout:      30 * time.Second,
	}
	return srv.ListenAndServe()
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

func (s *Server) URL() string {
	return fmt.Sprintf("http://%s", s.Addr())
}

func (s *Server) newRouter(base string) *Router {
	mux := s.router.PathPrefix(base).Subrouter()
	return &Router{
		base:     base,
		mux:      mux,
		services: []string{},
	}
}

func walkFn(route *mux.Route, router *mux.Router, ancestors []*mux.Route) error {
	path, err := route.GetPathTemplate()
	if err != nil {
		return err
	}
	fmt.Println(path)
	return nil
}

type Router struct {
	base     string
	mux      *mux.Router
	services []string
}

func (r *Router) Mount(mountPoint string, routes []Route) {
	mux := r.mux.PathPrefix(mountPoint).Subrouter()
	for _, route := range routes {
		mux.HandleFunc(route.Path, route.Func).Methods(route.Methods...)
	}
	r.services = append(r.services, mountPoint)
}

type Route struct {
	Func    http.HandlerFunc
	Path    string
	Methods []string
}
