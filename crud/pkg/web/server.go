package web

import (
	"fmt"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
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
	log.Println("start server!")
	return s.ListenAndServe()
}
