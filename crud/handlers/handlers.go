package handlers

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"

	"user-service/pkg/validation"
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
