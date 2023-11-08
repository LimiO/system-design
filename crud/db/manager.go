package db

import (
	"context"
	"database/sql"
	"fmt"

	_ "github.com/mattn/go-sqlite3"

	"user-service/pkg/models"
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
	conn *sql.Conn
	db   *sql.DB
}

func NewManager() (*Manager, error) {
	db, err := sql.Open("sqlite3", dbName)
	if err != nil {
		return nil, fmt.Errorf("failed to open sqlite3 db %q: %v", dbName, err)
	}
	conn, err := db.Conn(context.Background())
	if err != nil {
		return nil, fmt.Errorf("failed to init db connection: %v", err)
	}
	m := &Manager{
		conn: conn,
		db:   db,
	}
	if err = m.CreateTables(); err != nil {
		return nil, fmt.Errorf("failed to create tables: %v", err)
	}
	return &Manager{
		conn: conn,
		db:   db,
	}, nil
}

func (m *Manager) CreateTables() error {
	if _, err := m.db.Exec(createTableQuery); err != nil {
		return fmt.Errorf("failed to execute query: %v", err)
	}
	return nil
}

func (m *Manager) GetUser(user *models.User) (*models.User, error) {
	row := m.db.QueryRow(getUserQuery, user.Username)
	if err := row.Err(); err != nil {
		return nil, fmt.Errorf("failed to get user %q: %v", user.Username, err)
	}
	result := &models.User{}
	if err := row.Scan(&result.Username, &result.FirstName, &result.LastName, &result.Email, &result.Phone); err != nil {
		return nil, err
	}
	return result, nil
}

func (m *Manager) CreateUser(user *models.User) error {
	_, err := m.db.Exec(createUserQuery, user.Username, user.FirstName, user.LastName, user.Email, user.Phone)
	if err != nil {
		return fmt.Errorf("failed to create user %q: %v", user.Username, err)
	}
	return nil
}

func (m *Manager) UpdateUser(user *models.User) error {
	preparation, err := m.db.Prepare(updateUserQuery)
	if err != nil {
		return fmt.Errorf("failed to prepare update user query: %v", err)
	}
	_, err = preparation.Exec(user.FirstName, user.LastName, user.Email, user.Phone, user.Username)
	fmt.Println(user.Username, user.Email)
	if err != nil {
		return fmt.Errorf("failed to update user %q: %v", user.Username, err)
	}
	return nil
}

func (m *Manager) DeleteUser(user *models.User) error {
	_, err := m.db.Exec(deleteUserQuery, user.Username)
	if err != nil {
		return fmt.Errorf("failed to delete user %q: %v", user.Username, err)
	}
	return nil
}
