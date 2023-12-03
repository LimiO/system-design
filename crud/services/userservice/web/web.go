package web

import (
	"fmt"
	"github.com/go-chi/chi/v5"

	"onlinestore/pkg/web"
	"onlinestore/pkg/web/middlewares"
)

type Server struct {
	*web.Server
}

func NewServer(addr string, port int, jwtSecret string) (*Server, error) {
	r := chi.NewRouter()
	handleManager, err := NewHandlerManager()
	if err != nil {
		return nil, fmt.Errorf("failed to create handle manager: %v", err)
	}
	server := &Server{
		web.NewServer(addr, port, r),
	}
	mwManager := middlewares.NewMiddlewareManager(jwtSecret)
	r.Use(mwManager.RecoverRequest)
	r.Use(mwManager.Authenticate)

	// TODO(albert-si) authorize user to do post requests

	r.Get("/user/{username}", handleManager.GetUser)
	r.Put("/user/{username}", handleManager.PutUser)
	r.Delete("/user/{username}", handleManager.DeleteUser)
	r.Get("/health", handleManager.Health)
	r.Post("/user", handleManager.PostUser)
	return server, nil
}
