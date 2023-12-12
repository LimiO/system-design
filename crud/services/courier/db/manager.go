package db

import (
	"fmt"

	"onlinestore/db"
	"onlinestore/pkg/models"
)

const (
	dbName    = "courier.db"
	initQuery = `
	CREATE TABLE IF NOT EXISTS couriers (
		username VARCHAR(64) NOT NULL PRIMARY KEY,
		status INTEGER NOT NULL
	);
	INSERT OR IGNORE INTO couriers(username, status) VALUES('cour1', 0);
	INSERT OR IGNORE INTO couriers(username, status) VALUES('cour2', 0);`
)

var (
	getCourQuery      = fmt.Sprintf("SELECT * FROM couriers WHERE username = ?;")
	getFreeCourQuery  = fmt.Sprintf("SELECT * FROM couriers WHERE status = 0 LIMIT 1;")
	updateStatusQuery = fmt.Sprintf(`UPDATE couriers
		SET status = ?
		WHERE username = ?;`)
)

type ReserveStatus int

const (
	Free     ReserveStatus = 0
	Reserved ReserveStatus = 1
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

func (m *Manager) GetCourier(username string) (*models.Courier, error) {
	result := models.Courier{}
	scanFields := []interface{}{&result.Username, &result.Status}
	found, err := m.Get(getCourQuery, []interface{}{username}, scanFields)
	if found {
		return &result, nil
	}
	return nil, err
}

func (m *Manager) UpdateStatus(username string, status ReserveStatus) error {
	result, err := m.GetDB().Exec(updateStatusQuery, status, username)
	if err != nil {
		return fmt.Errorf("failed to update status: %v", err)
	}
	if affected, err := result.RowsAffected(); affected == 0 || err != nil {
		return fmt.Errorf("zero rows affected")
	}
	return nil
}

func (m *Manager) GetFreeCourier() (*models.Courier, error) {
	result := models.Courier{}
	scanFields := []interface{}{&result.Username, &result.Status}
	found, err := m.Get(getFreeCourQuery, []interface{}{}, scanFields)
	if found {
		return &result, nil
	}
	return nil, err
}
