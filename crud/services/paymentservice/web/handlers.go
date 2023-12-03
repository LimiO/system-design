package web

import (
	"fmt"
	"net/http"
	"onlinestore/pkg/models"

	"onlinestore/internal/jwt"
	"onlinestore/pkg/web"
	"onlinestore/services/paymentservice/db"
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

func (h *HandlerManager) AddBalance(w http.ResponseWriter, r *http.Request) {
	req, err := web.DecodeHttpBody[AddBalanceRequest](r.Body)
	if err != nil {
		web.WriteError(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = h.dbManager.AddBalance(req.Username, req.Amount)
	if err != nil {
		web.WriteError(w, err.Error(), http.StatusBadRequest)
		return
	}

	web.WriteStatusOK(w)
	return
}

func (h *HandlerManager) SubBalance(w http.ResponseWriter, r *http.Request) {
	req, err := web.DecodeHttpBody[SubBalanceRequest](r.Body)
	if err != nil {
		web.WriteError(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = h.dbManager.SubBalance(req.Username, req.Amount)
	if err != nil {
		web.WriteError(w, err.Error(), http.StatusBadRequest)
		return
	}

	web.WriteStatusOK(w)
	return
}

func (h *HandlerManager) GetBalance(w http.ResponseWriter, r *http.Request) {
	req, err := web.DecodeHttpBody[GetBalanceRequest](r.Body)
	if err != nil {
		web.WriteError(w, err.Error(), http.StatusBadRequest)
		return
	}

	balanceInfo, err := h.dbManager.GetBalance(req.Username)
	if err != nil {
		web.WriteError(w, err.Error(), http.StatusBadRequest)
		return
	}
	if balanceInfo == nil {
		balanceInfo = &models.BalanceInfo{
			Balance: 0,
		}
	}
	web.WriteData(w, &GetBalanceResponse{Balance: balanceInfo.Balance})
	return
}
