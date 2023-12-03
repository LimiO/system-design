package web

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/go-chi/chi/v5"
	"log"
	"net/http"
	"onlinestore/pkg/models"
	"onlinestore/pkg/web"

	"onlinestore/services/userservice/db"
)

type HandlerManager struct {
	dbManager *db.Manager
}

func NewHandlerManager() (*HandlerManager, error) {
	dbManager, err := db.NewManager()
	if err != nil {
		return nil, fmt.Errorf("failed to create db manager: %v", err)
	}
	return &HandlerManager{
		dbManager: dbManager,
	}, nil
}

func (h *HandlerManager) GetUserOrWriteError(w http.ResponseWriter, username string) (*models.User, error) {
	user, err := h.dbManager.GetUser(&models.User{Username: username})
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			log.Printf("user not found: %v", err)
			web.WriteNotFound(w, fmt.Sprintf("user not found: %v", err))
		} else {
			log.Printf("failed to get user: %v", err)
			web.WriteBadRequest(w, err.Error())
		}
		return nil, fmt.Errorf("get user error: %v", err)
	}
	return user, nil
}

func (h *HandlerManager) Health(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	web.WriteStatusOK(w)
}

func (h *HandlerManager) GetUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	username := chi.URLParam(r, "username")
	user, err := h.GetUserOrWriteError(w, username)
	if err != nil {
		log.Printf("failed to get user: %v", err)
		return
	}
	user.Password = ""

	fmt.Printf("user %q successfuly getted\n", user.Username)
	web.WriteData(w, user)
}

func (h *HandlerManager) PostUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	user, err := web.DecodeHttpBody[models.User](r.Body)
	if err != nil {
		log.Printf("failed to decode http body: %v", err)
		web.WriteBadRequest(w, err.Error())
		return
	}

	validationErrors := user.Validate()
	if len(validationErrors) > 0 {
		web.WriteValidationErrors(w, validationErrors)
		return
	}

	if err = h.dbManager.CreateUser(user); err != nil {
		log.Printf("failed to create user: %v", err)
		web.WriteBadRequest(w, err.Error())
		return
	}

	fmt.Printf("user %q successfuly created\n", user.Username)
	web.WriteStatusOK(w)
}

func (h *HandlerManager) PutUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	username := chi.URLParam(r, "username")
	oldUser, err := h.GetUserOrWriteError(w, username)
	if err != nil {
		log.Printf("failed to get user: %v", err)
		return
	}

	newData, err := web.DecodeHttpBody[models.User](r.Body)
	if err != nil {
		log.Printf("failed to decode http body: %v", err)
		web.WriteBadRequest(w, err.Error())
		return
	}
	validationErrors := newData.Validate()
	if len(validationErrors) > 0 {
		web.WriteValidationErrors(w, validationErrors)
		return
	}

	FillByDefaults(newData, oldUser)

	newData.Username = username
	if err = h.dbManager.UpdateUser(newData); err != nil {
		log.Printf("failed to update user: %v", err)
		web.WriteBadRequest(w, err.Error())
		return
	}

	fmt.Printf("user %q successfuly updated\n", username)
	web.WriteStatusOK(w)
}

func (h *HandlerManager) DeleteUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	username := chi.URLParam(r, "username")
	user, err := h.GetUserOrWriteError(w, username)
	if err != nil {
		log.Printf("failed to get user: %v", err)
		return
	}

	if err = h.dbManager.DeleteUser(user); err != nil {
		log.Printf("failed to delete user: %v", err)
		web.WriteBadRequest(w, err.Error())
		return
	}

	fmt.Printf("user %q successfuly deleted\n", username)
	web.WriteStatusOK(w)
}
