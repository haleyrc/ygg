package server

import (
	"context"
	"fmt"
	"net/http"
	"time"
)

func New(ctx context.Context, ip string, port int, handler http.Handler) (*Server, error) {
	srv := &Server{
		httpServer: &http.Server{
			Addr:              fmt.Sprintf("%s:%d", ip, port),
			Handler:           handler,
			ReadTimeout:       30 * time.Second,
			ReadHeaderTimeout: 5 * time.Second,
			WriteTimeout:      30 * time.Second,
		},
	}

	return srv, nil
}

type Server struct {
	httpServer *http.Server
}

func (s *Server) Addr() string {
	return s.httpServer.Addr
}

func (s *Server) ListenAndServe() error {
	return s.httpServer.ListenAndServe()
}
