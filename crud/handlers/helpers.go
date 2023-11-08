package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"user-service/pkg/validation"
)

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
