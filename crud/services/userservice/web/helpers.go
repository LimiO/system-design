package handlers

import (
	"onlinestore/pkg/models"
)

func FillByDefaults(user *models.User, oldUser *models.User) {
	if user.Email == "" {
		user.Email = oldUser.Email
	}
	if user.FirstName == "" {
		user.FirstName = oldUser.FirstName
	}
	if user.LastName == "" {
		user.LastName = oldUser.LastName
	}
	if user.Phone == 0 {
		user.Phone = oldUser.Phone
	}
	if user.Password == "" {
		user.Password = oldUser.Password
	}
	if user.Username == "" {
		user.Username = oldUser.Username
	}
}
