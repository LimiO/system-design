package web

import (
	"fmt"
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

func (h *HandlerManager) GetUserOrStatusCode(username string) (*models.User, int, error) {
	user, err := h.dbManager.GetUser(username)
	if err != nil {
		log.Printf("failed to get user: %v", err)
		return nil, http.StatusBadRequest, fmt.Errorf("get user error: %v", err)
	}
	if user == nil {
		log.Printf("user not found: %v", err)
		return nil, http.StatusNotFound, fmt.Errorf("user not found")
	}
	return user, http.StatusOK, nil
}

func (h *HandlerManager) Health(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	web.WriteStatusOK(w)
}

func (h *HandlerManager) GetUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	username := web.GetLogin(r.Context())
	user, code, err := h.GetUserOrStatusCode(username)
	if err != nil {
		log.Printf("failed to get user: %v", err)
		web.WriteError(w, err.Error(), code)
		return
	}
	user.Password = ""

	log.Printf("user %q successfuly getted\n", user.Username)
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

	if err = h.dbManager.CreateUser(user.FirstName, user.LastName, user.Email, user.Phone, user.Username); err != nil {
		log.Printf("failed to create user: %v", err)
		web.WriteBadRequest(w, err.Error())
		return
	}

	log.Printf("user %q successfuly created\n", user.Username)
	web.WriteStatusOK(w)
}

func (h *HandlerManager) PutUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	username := web.GetLogin(r.Context())
	oldUser, code, err := h.GetUserOrStatusCode(username)
	if err != nil {
		log.Printf("failed to get user: %v", err)
		web.WriteError(w, err.Error(), code)
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
	if err = h.dbManager.UpdateUser(newData.FirstName, newData.LastName, newData.Email, newData.Phone, newData.Username); err != nil {
		log.Printf("failed to update user: %v", err)
		web.WriteBadRequest(w, err.Error())
		return
	}

	log.Printf("user %q successfuly updated\n", username)
	web.WriteStatusOK(w)
}

func (h *HandlerManager) DeleteUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	username := web.GetLogin(r.Context())
	user, code, err := h.GetUserOrStatusCode(username)
	if err != nil {
		log.Printf("failed to get user: %v", err)
		web.WriteError(w, err.Error(), code)
		return
	}

	if err = h.dbManager.DeleteUser(user.Username); err != nil {
		log.Printf("failed to delete user: %v", err)
		web.WriteBadRequest(w, err.Error())
		return
	}

	log.Printf("user %q successfuly deleted\n", username)
	web.WriteStatusOK(w)
}
