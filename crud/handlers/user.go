package handlers

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"

	"user-service/db"
	"user-service/pkg/models"
)

type HandlerManager struct {
	dbManager *db.Manager
}

func DecodeHttpBody[T interface{}](body io.ReadCloser) (*T, error) {
	t := new(T)
	decoder := json.NewDecoder(body)
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(t); err != nil {
		return nil, fmt.Errorf("failed to read data: %v", err)
	}
	return t, nil
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
			WriteNotFound(w)
		} else {
			log.Printf("failed to get user: %v", err)
			WriteBadRequest(w, err.Error())
		}
		return nil, fmt.Errorf("get user error: %v", err)
	}
	return user, nil
}

func (h *HandlerManager) Health(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	WriteStatusOK(w)
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
	WriteJson(w, user)
}

func (h *HandlerManager) PostUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	user, err := DecodeHttpBody[models.User](r.Body)
	if err != nil {
		log.Printf("failed to decode http body: %v", err)
		WriteBadRequest(w, err.Error())
		return
	}

	validationErrors := user.Validate()
	if len(validationErrors) > 0 {
		WriteValidationErrors(w, validationErrors)
		return
	}

	if err = h.dbManager.CreateUser(user); err != nil {
		log.Printf("failed to create user: %v", err)
		WriteBadRequest(w, err.Error())
		return
	}

	fmt.Printf("user %q successfuly created\n", user.Username)
	WriteStatusOK(w)
}

func (h *HandlerManager) PutUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	username := chi.URLParam(r, "username")
	oldUser, err := h.GetUserOrWriteError(w, username)
	if err != nil {
		log.Printf("failed to get user: %v", err)
		return
	}

	user, err := DecodeHttpBody[models.User](r.Body)
	if err != nil {
		log.Printf("failed to decode http body: %v", err)
		WriteBadRequest(w, err.Error())
		return
	}
	validationErrors := user.Validate()
	if len(validationErrors) > 0 {
		WriteValidationErrors(w, validationErrors)
		return
	}

	FillByDefaults(user, oldUser)

	user.Username = username
	if err = h.dbManager.UpdateUser(user); err != nil {
		log.Printf("failed to update user: %v", err)
		WriteBadRequest(w, err.Error())
		return
	}

	fmt.Printf("user %q successfuly updated\n", username)
	WriteStatusOK(w)
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
		WriteBadRequest(w, err.Error())
		return
	}

	fmt.Printf("user %q successfuly deleted\n", username)
	WriteStatusOK(w)
}
