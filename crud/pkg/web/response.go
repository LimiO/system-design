package web

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"onlinestore/pkg/validation"
)

func WriteError(w http.ResponseWriter, msg string, code int) {
	w.WriteHeader(code)
	WriteData(w, &RequestStatus{"error", msg})
}

func WriteBadRequest(w http.ResponseWriter, msg string) {
	w.WriteHeader(http.StatusBadRequest)
	WriteData(w, &RequestStatus{"error", fmt.Sprintf("bad request: %v", msg)})
}

func WriteStatusOK(w http.ResponseWriter) {
	w.WriteHeader(http.StatusOK)
	WriteData(w, &RequestStatus{Status: "ok"})
}

func WriteValidationErrors(w http.ResponseWriter, validationErrors []*validation.ErrorItem) {
	w.WriteHeader(http.StatusUnprocessableEntity)
	validationError := NewHTTPValidationError(validationErrors)
	WriteData(w, validationError)
}

func WriteData(w http.ResponseWriter, value any) {
	err := json.NewEncoder(w).Encode(value)
	if err != nil {
		log.Printf("failed to write json: %v", err)
	}
}
