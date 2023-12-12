package db

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	_ "github.com/mattn/go-sqlite3"
)

const (
	Unknown   CommitStatus = 0
	Committed CommitStatus = 1
	Cancelled CommitStatus = 2
)

var (
	commitQuery = fmt.Sprintf(`UPDATE reserve SET status = ? WHERE reserve_id = ?;`)
)

type Manager struct {
	conn *sql.Conn
	db   *sql.DB
}

type CommitStatus int

func NewManager(dbName string, initReq string) (*Manager, error) {
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
	if err = m.CreateTables(initReq); err != nil {
		return nil, fmt.Errorf("failed to create tables: %v", err)
	}
	return m, nil
}

func (m *Manager) GetDB() *sql.DB {
	return m.db
}

func (m *Manager) CreateTables(initReq string) error {
	if _, err := m.db.Exec(initReq); err != nil {
		return fmt.Errorf("failed to execute query: %v", err)
	}
	return nil
}

func (m *Manager) Get(query string, params []interface{}, scanFields []interface{}) (bool, error) {
	row := m.GetDB().QueryRow(query, params...)
	if err := row.Err(); err != nil {
		return false, fmt.Errorf("failed to get: %v", err)
	}
	if err := row.Scan(scanFields...); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

func (m *Manager) List(query string, params []interface{}, scanFields []interface{}) (bool, error) {
	row := m.GetDB().QueryRow(query, params...)
	if err := row.Err(); err != nil {
		return false, fmt.Errorf("failed to get: %v", err)
	}
	if err := row.Scan(scanFields...); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

func (m *Manager) Commit(reserveID int, status CommitStatus, rollbackQuery string, rollbackArgs ...any) error {
	_, err := m.GetDB().Exec(commitQuery, status, reserveID)
	if err != nil {
		return fmt.Errorf("failed to commit: %v", err)
	}

	if status == Cancelled {
		_, err = m.GetDB().Exec(rollbackQuery, rollbackArgs...)
		if err != nil {
			return fmt.Errorf("failed to rollback: %v", err)
		}
	}
	return nil
}
