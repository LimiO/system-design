package server

import "onlinestore/pkg/validation"

type ResponseMetadata struct {
	Code  int    `json:"code"`
	Error string `json:"error"`
}

type RequestStatus struct {
	Status string `json:"status"`
	Error  string `json:"error,omitempty"`
}

type HTTPValidationError struct {
	Details []*validation.ErrorItem `json:"details"`
}

func NewHTTPValidationError(details []*validation.ErrorItem) *HTTPValidationError {
	return &HTTPValidationError{
		Details: details,
	}
}
