package web

import (
	"fmt"
	"net/http"
	"onlinestore/internal/jwt"
	"onlinestore/pkg/web"
	"onlinestore/services/purchaseservice/db"
)

type HandlerManager struct {
	dbManager    *db.Manager
	tokenManager *jwt.TokenManager
}

func NewHandlerManager(jwtSecret string) (*HandlerManager, error) {
	tokenManager := jwt.NewTokenManager(jwtSecret)
	dbManager, err := db.NewManager()
	if err != nil {
		return nil, fmt.Errorf("failed to create db manager: %v", err)
	}
	return &HandlerManager{
		dbManager:    dbManager,
		tokenManager: tokenManager,
	}, nil
}

func (h *HandlerManager) Buy(w http.ResponseWriter, r *http.Request) {
	req, err := web.DecodeHttpBody[BuyRequest](r.Body)
	if err != nil {
		web.WriteError(w, err.Error(), http.StatusBadRequest)
		return
	}

	username := web.GetLogin(r.Context())
	if err = ValidateBuyRequest(req); err != nil {
		web.WriteError(w, err.Error(), http.StatusBadRequest)
		return
	}

	orderID, err := h.dbManager.CreateOrder(req.ProductID, req.Count, req.Price, username)
	if err != nil {
		web.WriteError(w, err.Error(), http.StatusBadRequest)
		return
	}
	web.WriteData(w, BuyResponse{OrderID: orderID, Total: req.Count * req.Price})
	return
}

func (h *HandlerManager) CommitOrder(w http.ResponseWriter, r *http.Request) {
	req, err := web.DecodeHttpBody[CommitOrderRequest](r.Body)
	if err != nil {
		web.WriteError(w, err.Error(), http.StatusBadRequest)
		return
	}
	err = h.dbManager.UpdateOrder(req.OrderID, db.Paid)
	if err != nil {
		web.WriteError(w, err.Error(), http.StatusBadRequest)
		return
	}

	web.WriteData(w, CommitOrderResponse{})
	return
}

func (h *HandlerManager) GetOrders(w http.ResponseWriter, r *http.Request) {
	req, err := web.DecodeHttpBody[GetOrderRequest](r.Body)
	if err != nil {
		web.WriteError(w, err.Error(), http.StatusBadRequest)
		return
	}

	// TODO(albert-si) add username to all requests and filter orders by username
	order, err := h.dbManager.GetOrder(req.OrderID)
	if err != nil {
		web.WriteError(w, err.Error(), http.StatusBadRequest)
		return
	}
	if order == nil {
		web.WriteError(w, fmt.Sprintf("order not found"), http.StatusNotFound)
		return
	}
	web.WriteData(w, &GetOrderResponse{Order: order})
}
