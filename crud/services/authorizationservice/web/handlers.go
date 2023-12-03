package web

import (
	"fmt"
	"log"
	"net/http"

	"onlinestore/internal/jwt"
	"onlinestore/pkg/web"
	"onlinestore/services/authorizationservice/db"
)

type HandlerManager struct {
	dbManager    *db.Manager
	tokenManager *jwt.TokenManager
}

// TODO(albert-si) replace all internal errors with log fatal

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

func (h *HandlerManager) Register(w http.ResponseWriter, r *http.Request) {
	req, err := web.DecodeHttpBody[TokenRequest](r.Body)
	if err != nil {
		web.WriteError(w, err.Error(), http.StatusBadRequest)
		return
	}
	if req.Password == "" || req.Username == "" {
		web.WriteError(w, fmt.Sprintf("empty password or username"), http.StatusBadRequest)
		return
	}

	passInfo, err := h.dbManager.GetUserPassword(req.Username)
	if err != nil {
		log.Fatalf("failed to get user password: %v", err)
	}
	if passInfo != nil {
		web.WriteError(w, fmt.Sprintf("user already exists"), http.StatusConflict)
		return
	}

	hashed := jwt.MakeMD5Hash(req.Password)
	if err = h.dbManager.CreateUserPassword(req.Username, hashed); err != nil {
		web.WriteError(w, err.Error(), http.StatusBadRequest)
		return
	}

	token, err := h.tokenManager.CreateToken(req.Username)
	if err != nil {
		web.WriteError(w, err.Error(), http.StatusBadRequest)
		return
	}
	web.WriteData(w, &TokenResponse{Token: token})
}

func (h *HandlerManager) GetToken(w http.ResponseWriter, r *http.Request) {
	req, err := web.DecodeHttpBody[TokenRequest](r.Body)
	if err != nil {
		web.WriteError(w, err.Error(), http.StatusBadRequest)
		return
	}
	passInfo, err := h.dbManager.GetUserPassword(req.Username)
	if err != nil {
		web.WriteError(w, err.Error(), http.StatusBadRequest)
		return
	}
	if passInfo == nil {
		web.WriteError(w, fmt.Sprintf("user not found"), http.StatusNotFound)
		return
	}

	if passInfo.Passhash != jwt.MakeMD5Hash(req.Password) {
		web.WriteError(w, fmt.Sprintf("wrong password"), http.StatusForbidden)
		return
	}

	token, err := h.tokenManager.CreateToken(req.Username)
	if err != nil {
		web.WriteError(w, err.Error(), http.StatusBadRequest)
		return
	}
	web.WriteData(w, &TokenResponse{Token: token})
}

func (h *HandlerManager) Unregister(w http.ResponseWriter, r *http.Request) {
	req, err := web.DecodeHttpBody[DeleteTokenRequest](r.Body)
	if err != nil {
		web.WriteError(w, err.Error(), http.StatusBadRequest)
		return
	}
	passInfo, err := h.dbManager.GetUserPassword(req.Username)
	if err != nil {
		web.WriteError(w, err.Error(), http.StatusBadRequest)
		return
	}
	if passInfo == nil {
		web.WriteError(w, fmt.Sprintf("user not found"), http.StatusNotFound)
		return
	}

	if passInfo.Passhash != jwt.MakeMD5Hash(req.Password) {
		web.WriteError(w, fmt.Sprintf("wrong password"), http.StatusForbidden)
		return
	}

	if err = h.dbManager.DeleteUserPassword(req.Username); err != nil {
		web.WriteError(w, fmt.Sprintf("failed to delete user password: %v", err), http.StatusForbidden)
		return
	}

	web.WriteData(w, &DeleteTokenResponse{})
}
