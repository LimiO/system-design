package web

import (
	"fmt"
	"net/http"
	"onlinestore/pkg/models"
	"onlinestore/services/paymentservice/types"

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
	req, err := web.DecodeHttpBody[types.AddBalanceRequest](r.Body)
	if err != nil {
		web.WriteBadRequest(w, err.Error())
		return
	}

	err = h.dbManager.AddBalance(req.Username, req.Amount)
	if err != nil {
		web.WriteBadRequest(w, err.Error())
		return
	}

	web.WriteStatusOK(w)
	return
}

func (h *HandlerManager) SubBalance(w http.ResponseWriter, r *http.Request) {
	req, err := web.DecodeHttpBody[types.SubBalanceRequest](r.Body)
	if err != nil {
		web.WriteBadRequest(w, err.Error())
		return
	}

	err = h.dbManager.SubBalance(req.Username, req.Amount)
	if err != nil {
		web.WriteBadRequest(w, err.Error())
		return
	}

	web.WriteStatusOK(w)
	return
}

func (h *HandlerManager) ReserveBalance(w http.ResponseWriter, r *http.Request) {
	req, err := web.DecodeHttpBody[types.ReserveBalanceRequest](r.Body)
	if err != nil {
		web.WriteBadRequest(w, err.Error())
		return
	}

	reserveID, err := h.dbManager.ReserveBalance(req.Username, req.Amount)
	if err != nil {
		web.WriteBadRequest(w, err.Error())
		return
	}

	web.WriteData(w, types.ReserveBalanceResponse{
		ReserveID: reserveID,
	})
	return
}

func (h *HandlerManager) Commit(w http.ResponseWriter, r *http.Request) {
	req, err := web.DecodeHttpBody[types.CommitRequest](r.Body)
	if err != nil {
		web.WriteBadRequest(w, err.Error())
		return
	}

	err = h.dbManager.Commit(req.ReserveID, req.Status)
	if err != nil {
		web.WriteBadRequest(w, err.Error())
		return
	}
	web.WriteStatusOK(w)
}

func (h *HandlerManager) GetBalance(w http.ResponseWriter, r *http.Request) {
	req, err := web.DecodeHttpBody[types.GetBalanceRequest](r.Body)
	if err != nil {
		web.WriteBadRequest(w, err.Error())
		return
	}

	balanceInfo, err := h.dbManager.GetBalance(req.Username)
	if err != nil {
		web.WriteBadRequest(w, err.Error())
		return
	}
	if balanceInfo == nil {
		balanceInfo = &models.BalanceInfo{
			Balance: 0,
		}
	}
	web.WriteData(w, &types.GetBalanceResponse{Balance: balanceInfo.Balance})
}
