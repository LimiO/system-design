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
	handleManager, err := NewHandlerManager(jwtSecret)
	if err != nil {
		return nil, fmt.Errorf("failed to create handle manager: %v", err)
	}
	server := &Server{
		web.NewServer(addr, port, r),
	}
	mwManager := middlewares.NewMiddlewareManager(jwtSecret)
	r.Use(mwManager.RecoverRequest)

	r.Post("/register", handleManager.Register)
	r.Get("/token", handleManager.GetToken)
	r.Post("/unregister", handleManager.Unregister)
	return server, nil
}
