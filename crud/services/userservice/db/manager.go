package db

import (
	"database/sql"
	"errors"
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
		phone INT NOT NULL
	)`
	getUserQuery    = `SELECT * FROM users WHERE username = ?;`
	createUserQuery = `INSERT INTO users(username, first_name, last_name, email, phone) VALUES (?, ?, ?, ?, ?);`
	updateUserQuery = `UPDATE users
	SET first_name=?,
		last_name=?,
		email=?,
		phone=?
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

func (m *Manager) GetUser(username string) (*models.User, error) {
	row := m.GetDB().QueryRow(getUserQuery, username)
	if err := row.Err(); err != nil {
		return nil, fmt.Errorf("failed to get user %q: %v", username, err)
	}
	result := &models.User{}
	if err := row.Scan(&result.Username, &result.FirstName, &result.LastName, &result.Email, &result.Phone); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return result, nil
}

func (m *Manager) CreateUser(firstName, lastName, email string, phone int, username string) error {
	_, err := m.GetDB().Exec(createUserQuery, username, firstName, lastName, email, phone)
	if err != nil {
		return fmt.Errorf("failed to create user %q: %v", username, err)
	}
	return nil
}

func (m *Manager) UpdateUser(firstName, lastName, email string, phone int, username string) error {
	preparation, err := m.GetDB().Prepare(updateUserQuery)
	if err != nil {
		return fmt.Errorf("failed to prepare update user query: %v", err)
	}
	_, err = preparation.Exec(firstName, lastName, email, phone, username)
	if err != nil {
		return fmt.Errorf("failed to update user %q: %v", username, err)
	}
	return nil
}

func (m *Manager) DeleteUser(username string) error {
	_, err := m.GetDB().Exec(deleteUserQuery, username)
	if err != nil {
		return fmt.Errorf("failed to delete user %q: %v", username, err)
	}
	return nil
}
