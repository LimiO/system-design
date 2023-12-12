package db

import (
	"fmt"

	"onlinestore/db"
	"onlinestore/pkg/models"
)

type CommitStatus int

const (
	Unknown   CommitStatus = 0
	Committed CommitStatus = 1
	Cancelled CommitStatus = 2
)

const (
	dbName    = "stock.db"
	initQuery = `
	CREATE TABLE IF NOT EXISTS products (
		product_id INTEGER NOT NULL PRIMARY KEY,
		count INTEGER NOT NULL
	);
	CREATE TABLE IF NOT EXISTS reserve (
		reserve_id INTEGER NOT NULL PRIMARY KEY,
		product_id INTEGER NOT NULL,
		count INTEGER NOT NULL,
		status INTEGER NOT NULL
	);
	INSERT OR IGNORE INTO products(product_id, count) VALUES(0, 50);
	INSERT OR IGNORE INTO products(product_id, count) VALUES(1, 50);
	INSERT OR IGNORE INTO products(product_id, count) VALUES(2, 50);`
)

var (
	getProductsQuery = fmt.Sprintf("SELECT * FROM products WHERE product_id = ?;")
	addProductsQuery = fmt.Sprintf(`UPDATE products
		SET count = count + ?
		WHERE product_id = ?;`)
	subProductsQuery = fmt.Sprintf(`UPDATE products
		SET count = count - ?
		WHERE product_id = ? AND count - ? >= 0;`)
	reserveProductsQuery = fmt.Sprintf(`INSERT INTO reserve(product_id, count, status) VALUES(?, ?, 0);`)
	commitQuery          = fmt.Sprintf(`UPDATE reserve SET status = ? WHERE reserve_id = ?;`)
	selectReserveQuery   = fmt.Sprintf(`SELECT product_id, count FROM reserve WHERE reserve_id = ?`)
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

func (m *Manager) GetCount(productID int) (int, error) {
	result := models.ProductInfo{}
	scanFields := []interface{}{&result.ProductID, &result.Count}
	found, err := m.Get(getProductsQuery, []interface{}{productID}, scanFields)
	if found {
		return result.Count, nil
	}
	return 0, err
}

func (m *Manager) AddCount(productID int, count int) error {
	result, err := m.GetDB().Exec(addProductsQuery, productID, count)
	if err != nil {
		return fmt.Errorf("failed to add count: %v", err)
	}
	if affected, err := result.RowsAffected(); affected == 0 || err != nil {
		return fmt.Errorf("zero rows affected")
	}
	return nil
}

func (m *Manager) SubCount(productID int, count int) error {
	result, err := m.GetDB().Exec(subProductsQuery, productID, count)
	if err != nil {
		return fmt.Errorf("failed to sub count: %v", err)
	}
	if affected, err := result.RowsAffected(); affected == 0 || err != nil {
		return fmt.Errorf("zero rows affected")
	}
	return nil
}

func (m *Manager) ReserveCount(productID int, count int) (int, error) {
	// TODO(albert-si) add lock
	result, err := m.GetDB().Exec(subProductsQuery, count, productID, count)
	if err != nil {
		return 0, fmt.Errorf("failed to sub count: %v", err)
	}
	if affected, err := result.RowsAffected(); affected == 0 || err != nil {
		return 0, fmt.Errorf("zero rows affected")
	}

	result, err = m.GetDB().Exec(reserveProductsQuery, productID, count)
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
	var productID int
	var count int
	scanFields := []interface{}{&productID, &count}
	found, err := m.Get(selectReserveQuery, []interface{}{reserveID}, scanFields)
	if err != nil {
		return fmt.Errorf("failed to find reserve_id %d: %v", reserveID, err)
	}
	if !found {
		return nil
	}
	return m.Manager.Commit(reserveID, db.CommitStatus(status), addProductsQuery, count, productID)
}
