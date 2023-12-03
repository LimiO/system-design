package db

import (
	"fmt"
	"onlinestore/db"
	"onlinestore/pkg/models"
)

const (
	dbName           = "users.db"
	createTableQuery = `
	CREATE TABLE IF NOT EXISTS users (
		username VARCHAR(64) NOT NULL PRIMARY KEY,
		first_name VARCHAR(64) NOT NULL,
		last_name VARCHAR(64) NOT NULL,
		email VARCHAR(64) NOT NULL, 
		phone INT NOT NULL,
		password VARCHAR(64)
	)`
	getUserQuery    = `SELECT * FROM users WHERE username = ?;`
	createUserQuery = `INSERT INTO users(username, first_name, last_name, email, phone, password) VALUES (?, ?, ?, ?, ?, ?);`
	updateUserQuery = `UPDATE users
	SET first_name=?,
		last_name=?,
		email=?,
		phone=?,
		password=?
	WHERE username=?;
	`
	deleteUserQuery = `DELETE FROM users WHERE username = ?;`
)

type Manager struct {
	*db.Manager
}

func NewManager() (*Manager, error) {
	manager, err := db.NewManager(dbName, createTableQuery)
	if err != nil {
		return nil, fmt.Errorf("failed to create base manager: %v", err)
	}
	return &Manager{
		Manager: manager,
	}, nil
}

func (m *Manager) GetUser(user *models.User) (*models.User, error) {
	row := m.GetDB().QueryRow(getUserQuery, user.Username)
	if err := row.Err(); err != nil {
		return nil, fmt.Errorf("failed to get user %q: %v", user.Username, err)
	}
	result := &models.User{}
	if err := row.Scan(&result.Username, &result.FirstName, &result.LastName, &result.Email, &result.Phone, &result.Password); err != nil {
		return nil, err
	}
	return result, nil
}

func (m *Manager) CreateUser(user *models.User) error {
	_, err := m.GetDB().Exec(createUserQuery, user.Username, user.FirstName, user.LastName, user.Email, user.Phone, user.Password)
	if err != nil {
		return fmt.Errorf("failed to create user %q: %v", user.Username, err)
	}
	return nil
}

func (m *Manager) UpdateUser(user *models.User) error {
	preparation, err := m.GetDB().Prepare(updateUserQuery)
	if err != nil {
		return fmt.Errorf("failed to prepare update user query: %v", err)
	}
	_, err = preparation.Exec(user.FirstName, user.LastName, user.Email, user.Phone, user.Password, user.Username)
	fmt.Println(user.Username, user.Email)
	if err != nil {
		return fmt.Errorf("failed to update user %q: %v", user.Username, err)
	}
	return nil
}

func (m *Manager) DeleteUser(user *models.User) error {
	_, err := m.GetDB().Exec(deleteUserQuery, user.Username)
	if err != nil {
		return fmt.Errorf("failed to delete user %q: %v", user.Username, err)
	}
	return nil
}
