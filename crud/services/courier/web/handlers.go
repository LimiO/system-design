package web

import (
	"fmt"
	"net/http"

	"onlinestore/internal/jwt"
	"onlinestore/pkg/web"
	"onlinestore/services/courier/db"
	"onlinestore/services/courier/types"
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

func (h *HandlerManager) ReserveCourier(w http.ResponseWriter, r *http.Request) {
	// TODO(albert-si) add lock or do it in one method
	free, err := h.dbManager.GetFreeCourier()
	if err != nil {
		web.WriteBadRequest(w, fmt.Sprintf("cour not found: %v", err))
		return
	}
	err = h.dbManager.UpdateStatus(free.Username, db.Reserved)
	if err != nil {
		web.WriteBadRequest(w, err.Error())
		return
	}

	web.WriteData(w, types.ReserveCourierResponse{
		Username: free.Username,
	})
}

func (h *HandlerManager) UnreserveCourier(w http.ResponseWriter, r *http.Request) {
	req, err := web.DecodeHttpBody[types.UnreserveCourierRequest](r.Body)
	if err != nil {
		web.WriteBadRequest(w, err.Error())
		return
	}

	err = h.dbManager.UpdateStatus(req.Username, db.Free)
	if err != nil {
		web.WriteBadRequest(w, err.Error())
		return
	}

	web.WriteStatusOK(w)
	return
}

func (h *HandlerManager) GetCourier(w http.ResponseWriter, r *http.Request) {
	req, err := web.DecodeHttpBody[types.GetCourierRequest](r.Body)
	if err != nil {
		web.WriteBadRequest(w, err.Error())
		return
	}

	cour, err := h.dbManager.GetCourier(req.Username)
	if err != nil {
		web.WriteBadRequest(w, err.Error())
		return
	}

	web.WriteData(w, cour)
	return
}
