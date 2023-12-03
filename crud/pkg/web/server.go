package server

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"

	"onlinestore/pkg/server/middlewares"
	"onlinestore/services/userservice/handlers"
)

type Server struct {
	http.Server
}

func NewServer(addr string, port int) (*Server, error) {
	r := chi.NewRouter()
	handleManager, err := handlers.NewHandlerManager()
	if err != nil {
		return nil, fmt.Errorf("failed to create handle manager: %v", err)
	}
	server := &Server{
		http.Server{
			Addr:    fmt.Sprintf("%s:%d", addr, port),
			Handler: r,
		},
	}
	r.Use(middlewares.RecoverRequest)
	r.Use(middlewares.Authorize)

	r.Get("/user/{username}", handleManager.GetUser)
	r.Put("/user/{username}", handleManager.PutUser)
	r.Delete("/user/{username}", handleManager.DeleteUser)
	r.Get("/health", handleManager.Health)
	r.Post("/user", handleManager.PostUser)
	return server, nil
}

func (s *Server) Start() error {
	return s.ListenAndServe()
}
