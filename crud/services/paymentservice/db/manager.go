package db

import (
	"fmt"

	"onlinestore/db"
	"onlinestore/pkg/models"
)

const (
	dbName    = "payment.db"
	initQuery = `
	CREATE TABLE IF NOT EXISTS balance (
		username VARCHAR(64) NOT NULL PRIMARY KEY,
		balance INTEGER NOT NULL
	);
	CREATE TABLE IF NOT EXISTS history (
		tx_id INTEGER NOT NULL PRIMARY KEY,
		username VARCHAR(64) NOT NULL,
		diff INTEGER NOT NULL
	);`
	insertOrIgnore = `INSERT OR IGNORE INTO balance (username, balance) VALUES (?, 0)`
)

var (
	getBalanceQuery = fmt.Sprintf("%s; SELECT * FROM balance WHERE username = ?;", insertOrIgnore)
	addBalanceQuery = fmt.Sprintf(`%s; UPDATE balance
		SET balance = balance + ?
		WHERE username = ?;`, insertOrIgnore)
	subBalanceQuery = fmt.Sprintf(`%s; UPDATE balance
		SET balance = balance - ?
		WHERE username = ? AND balance - ? >= 0;`, insertOrIgnore)
)

type Manager struct {
	*db.Manager
}

func NewManager() (*Manager, error) {
	manager, err := db.NewManager(dbName, initQuery)
	if err != nil {
		return nil, fmt.Errorf("failed to create base manager: %v", err)
	}
	return &Manager{
		Manager: manager,
	}, nil
}

func (m *Manager) GetBalance(username string) (*models.BalanceInfo, error) {
	result := &models.BalanceInfo{}
	scanFields := []interface{}{&result.Username, &result.Balance}
	found, err := m.Get(getBalanceQuery, []interface{}{username, username}, scanFields)
	if found {
		return result, nil
	}
	return nil, err
}

func (m *Manager) AddBalance(username string, amount int) error {
	result, err := m.GetDB().Exec(addBalanceQuery, username, amount, username)
	if err != nil {
		return fmt.Errorf("failed to add balance: %v", err)
	}
	if affected, err := result.RowsAffected(); affected == 0 || err != nil {
		return fmt.Errorf("zero rows affected")
	}
	return nil
}

func (m *Manager) SubBalance(username string, amount int) error {
	result, err := m.GetDB().Exec(subBalanceQuery, username, amount, username, amount)
	if err != nil {
		return fmt.Errorf("failed to sub balance: %v", err)
	}
	if affected, err := result.RowsAffected(); affected == 0 || err != nil {
		return fmt.Errorf("zero rows affected")
	}
	return nil
}
