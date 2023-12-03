package web

import (
	"fmt"
	"github.com/go-chi/chi/v5"
	"net/http"
)

type Server struct {
	http.Server
}

func NewServer(addr string, port int, r chi.Router) *Server {
	return &Server{
		http.Server{
			Addr:    fmt.Sprintf("%s:%d", addr, port),
			Handler: r,
		},
	}
}

func (s *Server) Start() error {
	return s.ListenAndServe()
}
