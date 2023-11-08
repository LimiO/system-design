package models

import (
	"fmt"
	"user-service/pkg/validation"
)

var (
	InvalidLengthType = "invalid length"
	InvalidValue      = "invalid value"
)

type User struct {
	Username  string `db:"username" json:"username"`
	FirstName string `db:"first_name" json:"first_name"`
	LastName  string `db:"last_name" json:"last_name"`
	Email     string `db:"email" json:"email"`
	Phone     int    `db:"phone" json:"phone"`
	Password  string `db:"password" json:"password,omitempty"`
}

func (u *User) GetUsername() string {
	return u.Username
}

func (u *User) GetFirstName() string {
	return u.FirstName
}

func (u *User) GetLastName() string {
	return u.LastName
}

func (u *User) GetEmail() string {
	return u.Email
}

func (u *User) GetPhone() int {
	return u.Phone
}

func (u *User) GetPassword() string {
	return u.Password
}

func (u *User) Validate() []*validation.ValidationErrorItem {
	var result []*validation.ValidationErrorItem
	if err := validation.ValidateEmail(u.Email); err != nil {
		result = append(result, validation.NewValidationErrorItem(fmt.Sprintf("invalid email: %v", err), InvalidValue))
	}
	if len(u.Email) > 64 {
		result = append(result, validation.NewValidationErrorItem("email len can't be greater, than 64", InvalidLengthType))
	}
	if len(u.Username) > 64 {
		result = append(result, validation.NewValidationErrorItem("username len can't be greater, than 64", InvalidLengthType))
	}
	if len(u.FirstName) > 64 {
		result = append(result, validation.NewValidationErrorItem("first name len can't be greater, than 64", InvalidLengthType))
	}
	if len(u.LastName) > 64 {
		result = append(result, validation.NewValidationErrorItem("first name len can't be greater, than 64", InvalidLengthType))
	}
	if len(u.Password) > 64 {
		result = append(result, validation.NewValidationErrorItem("password len can't be greater, than 64", InvalidLengthType))
	}
	return result
}
