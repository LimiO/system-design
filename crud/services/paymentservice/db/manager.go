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
	CREATE TABLE IF NOT EXISTS reserve (
		reserve_id INTEGER NOT NULL PRIMARY KEY,
		username VARCHAR(64) NOT NULL,
		amount INTEGER NOT NULL,
		status INTEGER NOT NULL
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
	reserveBalanceQuery = fmt.Sprintf(`INSERT INTO reserve(username, amount, status) VALUES(?, ?, 0);`)
	selectReserveQuery  = fmt.Sprintf(`SELECT username, amount FROM reserve WHERE reserve_id = ?`)
)

type CommitStatus db.CommitStatus

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

func (m *Manager) ReserveBalance(username string, amount int) (int, error) {
	// TODO(albert-si) add lock
	if err := m.SubBalance(username, amount); err != nil {
		return 0, fmt.Errorf("failed to sub balance: %v", err)
	}

	result, err := m.GetDB().Exec(reserveBalanceQuery, username, amount)
	if err != nil {
		return 0, fmt.Errorf("failed to reserve count: %v", err)
	}
	inserted, err := result.LastInsertId()
	if inserted == 0 || err != nil {
		return 0, fmt.Errorf("zero rows affected")
	}
	return int(inserted), nil
}

func (m *Manager) Commit(reserveID int, status int) error {
	// TODO(albert-si) add lock
	var username string
	var amount int
	scanFields := []interface{}{&username, &amount}
	found, err := m.Get(selectReserveQuery, []interface{}{reserveID}, scanFields)
	if err != nil {
		return fmt.Errorf("failed to find reserve_id %d: %v", reserveID, err)
	}
	if !found {
		return nil
	}

	return m.Manager.Commit(reserveID, db.CommitStatus(status), addBalanceQuery, username, amount, username)
}
