package web

import (
	"fmt"
	"net/http"

	"onlinestore/internal/jwt"
	"onlinestore/pkg/web"
	"onlinestore/services/stock/db"
	"onlinestore/services/stock/types"
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

func (h *HandlerManager) Reserve(w http.ResponseWriter, r *http.Request) {
	req, err := web.DecodeHttpBody[types.ReserveRequest](r.Body)
	if err != nil {
		web.WriteBadRequest(w, err.Error())
		return
	}

	reserveID, err := h.dbManager.ReserveCount(req.ProductID, req.Count)
	if err != nil {
		web.WriteBadRequest(w, err.Error())
		return
	}

	web.WriteData(w, types.ReserveResponse{
		ReserveID: reserveID,
	})
}

// TODO(albert-si) add rollback query, not in commit
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

func (h *HandlerManager) AddCount(w http.ResponseWriter, r *http.Request) {
	req, err := web.DecodeHttpBody[types.AddCountRequest](r.Body)
	if err != nil {
		web.WriteBadRequest(w, err.Error())
		return
	}

	err = h.dbManager.AddCount(req.ProductID, req.Count)
	if err != nil {
		web.WriteBadRequest(w, err.Error())
		return
	}

	web.WriteStatusOK(w)
	return
}

func (h *HandlerManager) GetCount(w http.ResponseWriter, r *http.Request) {
	req, err := web.DecodeHttpBody[types.GetCountRequest](r.Body)
	if err != nil {
		web.WriteBadRequest(w, err.Error())
		return
	}

	count, err := h.dbManager.GetCount(req.ProductID)
	if err != nil {
		web.WriteBadRequest(w, err.Error())
		return
	}

	web.WriteData(w, types.GetCountResponse{
		Count: count,
	})
	return
}
