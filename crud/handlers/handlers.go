package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"user-service/pkg/validation"

	"github.com/go-chi/chi/v5"
)

var (
	NotFoundJson []byte
	OkJson       []byte
)

type RequestStatus struct {
	Status string `json:"status"`
	Error  string `json:"error,omitempty"`
}

type HTTPValidationError struct {
	Details []*validation.ValidationErrorItem `json:"details"`
}

func NewHTTPValidationError(details []*validation.ValidationErrorItem) *HTTPValidationError {
	return &HTTPValidationError{
		Details: details,
	}
}

func WriteBadRequest(w http.ResponseWriter, msg string) {
	w.WriteHeader(http.StatusBadRequest)
	data, _ := json.Marshal(&RequestStatus{"error", fmt.Sprintf("bad request: %v", msg)})
	if _, err := w.Write(data); err != nil {
		log.Printf("failed to send data to user: %v", err)
	}
}

func WriteNotFound(w http.ResponseWriter) {
	w.WriteHeader(http.StatusNotFound)
	if _, err := w.Write(NotFoundJson); err != nil {
		log.Printf("failed to send data to user: %v", err)
	}
}

func WriteStatusOK(w http.ResponseWriter) {
	w.WriteHeader(http.StatusOK)
	if _, err := w.Write(OkJson); err != nil {
		log.Printf("failed to send data to user: %v", err)
	}
}

func WriteValidationErrors(w http.ResponseWriter, validationErrors []*validation.ValidationErrorItem) {
	w.WriteHeader(http.StatusUnprocessableEntity)
	validationError := NewHTTPValidationError(validationErrors)
	_ = json.NewEncoder(w).Encode(validationError)
}

func WriteJson(w http.ResponseWriter, value any) {
	err := json.NewEncoder(w).Encode(value)
	if err != nil {
		log.Printf("failed to write json: %v", err)
	}
}

func DumpErrors() error {
	var err error
	if NotFoundJson, err = json.Marshal(&RequestStatus{"error", "User not found"}); err != nil {
		return err
	}
	if OkJson, err = json.Marshal(&RequestStatus{Status: "ok"}); err != nil {
		return err
	}
	return nil
}

func NewRouter() (http.Handler, error) {
	r := chi.NewRouter()
	handleManager, err := NewHandlerManager()
	if err != nil {
		return nil, fmt.Errorf("failed to create handle manager: %v", err)
	}
	r.Get("/user/{username}", handleManager.GetUser)
	r.Put("/user/{username}", handleManager.PutUser)
	r.Delete("/user/{username}", handleManager.DeleteUser)
	r.Get("/health", handleManager.Health)
	r.Post("/user", handleManager.PostUser)
	return r, nil
}
