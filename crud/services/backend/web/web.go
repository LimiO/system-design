package web

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"

	"onlinestore/pkg/web"
	"onlinestore/pkg/web/middlewares"
)

type Server struct {
	*web.Server
}

func NewServer(addr string, port int, jwtSecret string, PaymentAddr, AuthorizationAddr, PurchasesAddr, UserAddr string) (*Server, error) {
	r := chi.NewRouter()
	handleManager, err := NewHandlerManager(jwtSecret, PaymentAddr, AuthorizationAddr, PurchasesAddr, UserAddr)
	if err != nil {
		return nil, fmt.Errorf("failed to create handle manager: %v", err)
	}
	server := &Server{
		web.NewServer(addr, port, r),
	}
	mwManager := middlewares.NewMiddlewareManager(jwtSecret)
	r.Use(mwManager.RecoverRequest)

	// TODO(albert-si) add rate limiter

	// TODO(albert-si) add username to all requests
	r.Post("/billing/balance/add", mwManager.Authenticate(http.HandlerFunc(handleManager.AddBalance)).ServeHTTP)
	r.Post("/billing/balance/sub", mwManager.Authenticate(http.HandlerFunc(handleManager.SubBalance)).ServeHTTP)
	r.Get("/billing/balance", mwManager.Authenticate(http.HandlerFunc(handleManager.GetBalance)).ServeHTTP)
	r.Get("/user", mwManager.Authenticate(http.HandlerFunc(handleManager.GetUser)).ServeHTTP)
	r.Post("/user", handleManager.PostUser)
	r.Put("/user", mwManager.Authenticate(http.HandlerFunc(handleManager.PutUser)).ServeHTTP)
	r.Delete("/user", mwManager.Authenticate(http.HandlerFunc(handleManager.DeleteUser)).ServeHTTP)
	r.Get("/token", handleManager.GetToken)

	r.Get("/order", mwManager.Authenticate(http.HandlerFunc(handleManager.GetOrder)).ServeHTTP)
	r.Get("/orders", mwManager.Authenticate(http.HandlerFunc(handleManager.GetOrders)).ServeHTTP)
	r.Post("/buy", mwManager.Authenticate(http.HandlerFunc(handleManager.Buy)).ServeHTTP)

	return server, nil
}
