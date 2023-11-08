package validation

import (
	"fmt"
	"regexp"
)

type ValidationErrorItem struct {
	Msg            string `json:"Msg"`
	ValidationType string `json:"type"`
}

func NewValidationErrorItem(msg, validationType string) *ValidationErrorItem {
	return &ValidationErrorItem{
		Msg:            msg,
		ValidationType: validationType,
	}
}

func ValidateEmail(e string) error {
	emailRegex := regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")
	if emailRegex.MatchString(e) {
		return nil
	}
	return fmt.Errorf("email %q invalid", e)
}
