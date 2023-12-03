package db

import (
	"database/sql"
	"errors"
	"fmt"

	"onlinestore/db"
	"onlinestore/pkg/models"
)

const (
	dbName           = "passwords.db"
	createTableQuery = `
	CREATE TABLE IF NOT EXISTS passwords (
		username VARCHAR(64) NOT NULL PRIMARY KEY,
		passhash VARCHAR(64) NOT NULL
	)`
	getUserPasswordQuery = `SELECT * FROM passwords WHERE username = ?;`
	createPasswordQuery  = `INSERT INTO passwords(username, passhash) VALUES (?, ?);`
	updatePasswordQuery  = `UPDATE passwords
	SET username=?,
		passhash=?
	WHERE username=?;`
	deleteUserPasswordQuery = `DELETE FROM passwords WHERE username = ?;`
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

func (m *Manager) GetUserPassword(username string) (*models.PasswordInfo, error) {
	row := m.GetDB().QueryRow(getUserPasswordQuery, username)
	if err := row.Err(); err != nil {
		return nil, fmt.Errorf("failed to get user %q: %v", username, err)
	}
	result := &models.PasswordInfo{}
	if err := row.Scan(&result.Passhash, &result.Passhash); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return result, nil
}

func (m *Manager) CreateUserPassword(username string, passhash string) error {
	_, err := m.GetDB().Exec(createPasswordQuery, username, passhash)
	if err != nil {
		return fmt.Errorf("failed to create user %q: %v", username, err)
	}
	return nil
}

func (m *Manager) UpdateUserPassword(username string, passhash string) error {
	preparation, err := m.GetDB().Prepare(updatePasswordQuery)
	if err != nil {
		return fmt.Errorf("failed to prepare update user query: %v", err)
	}
	_, err = preparation.Exec(username, passhash)
	if err != nil {
		return fmt.Errorf("failed to update user %q: %v", username, err)
	}
	return nil
}

func (m *Manager) DeleteUserPassword(username string) error {
	_, err := m.GetDB().Exec(deleteUserPasswordQuery, username)
	if err != nil {
		return fmt.Errorf("failed to delete user %q: %v", username, err)
	}
	return nil
}
