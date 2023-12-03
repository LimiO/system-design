package db

import (
	"context"
	"fmt"
	"onlinestore/db"
	"onlinestore/pkg/models"
)

type PaidStatus int

var (
	NotPaid   PaidStatus = 0
	Paid      PaidStatus = 1
	Cancelled PaidStatus = 2
)

const (
	dbName    = "purchases.db"
	initQuery = `
	CREATE TABLE IF NOT EXISTS products (
		product_id INTEGER NOT NULL PRIMARY KEY,
		count INTEGER NOT NULL,
		price INTEGER NOT NULL
	);
	CREATE TABLE IF NOT EXISTS orders (
		order_id INTEGER NOT NULL PRIMARY KEY,
		product_id INTEGER NOT NULL,
		price INTEGER NOT NULL,
		count INTEGER NOT NULL,
		username VARCHAR(64) NOT NULL,
		paid INTEGER NOT NULL
	);
	INSERT OR IGNORE INTO products(product_id, count, price) VALUES(0, 10, 50);
	INSERT OR IGNORE INTO products(product_id, count, price) VALUES(1, 20, 30);
	INSERT OR IGNORE INTO products(product_id, count, price) VALUES(2, 30, 30);
	INSERT OR IGNORE INTO products(product_id, count, price) VALUES(3, 40, 30);
	INSERT OR IGNORE INTO products(product_id, count, price) VALUES(4, 50, 30);
	`
	getOrderQuery  = `SELECT * FROM orders WHERE order_id = ?;`
	listOrderQuery = `SELECT * 
	FROM orders
	WHERE username = ? LIMIT ?;`
	getProductQuery = `SELECT * FROM products WHERE product_id = ?;`
	addProductQuery = `UPDATE products
		SET count = count + ?
		WHERE product_id = ? AND (count + ?) >= 0;`
	subProductQuery = `UPDATE products
		SET count = count - ?
		WHERE product_id = ? AND (count - ?) >= 0 AND price = ?;`
	createOrderQuery = `INSERT INTO orders(product_id, price, count, username, paid)
		VALUES(?, (SELECT products.price FROM products WHERE product_id = ?), ?, ?, 0);`
	updateOrderQuery = `UPDATE orders
		SET paid = ?
		WHERE order_id = ?;`
	revertProductsQuery = `UPDATE products
	SET count = count + (
		SELECT count
		FROM orders
		WHERE products.product_id = orders.product_id
		AND orders.order_id = ?
	)
	WHERE product_id IN (
		SELECT product_id
		FROM orders
		WHERE order_id = ?
	);`
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

func (m *Manager) GetOrder(orderID int) (*models.Order, error) {
	result := &models.Order{}
	scanFields := []interface{}{&result.OrderID, &result.ProductID, &result.Price, &result.Count, &result.Username, &result.Paid}
	found, err := m.Get(getOrderQuery, []interface{}{orderID}, scanFields)
	if found {
		return result, nil
	}
	return nil, err
}

func (m *Manager) GetOrders(username string, limit int) ([]*models.Order, error) {
	rows, err := m.GetDB().Query(listOrderQuery, username, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to get orders: %v", err)
	}
	defer rows.Close()

	var result []*models.Order
	for rows.Next() {
		var order models.Order
		err = rows.Scan(&order.OrderID, &order.ProductID, &order.Price, &order.Count, &order.Username, &order.Paid)
		if err != nil {
			return nil, fmt.Errorf("failed to scan row: %v", err)
		}
		result = append(result, &order)
	}

	return result, nil
}

func (m *Manager) GetProduct(productID int) (*models.Product, error) {
	result := &models.Product{}
	scanFields := []any{&result.ProductID, &result.Price, &result.Count}
	found, err := m.Get(getProductQuery, []any{productID}, scanFields)
	if found {
		return result, nil
	}
	return nil, err
}

func (m *Manager) AddProduct(productID int, countToChange int) error {
	_, err := m.GetDB().Exec(addProductQuery, productID, countToChange)
	if err != nil {
		return fmt.Errorf("failed to add product %q count: %v", productID, err)
	}
	return nil
}

func (m *Manager) CreateOrder(productID int, count int, price int, username string) (int, error) {
	var err error
	tx, err := m.GetDB().BeginTx(context.Background(), nil)
	defer func() {
		if err != nil {
			_ = tx.Rollback()
		} else {
			_ = tx.Commit()
		}
	}()

	result, err := m.GetDB().Exec(subProductQuery, count, productID, count, price)
	if err != nil {
		return 0, fmt.Errorf("failed to sub product %q count: %v", productID, err)
	}
	affected, err := result.RowsAffected()
	if err != nil || affected == 0 {
		return 0, fmt.Errorf("can't sub from products")
	}

	result, err = m.GetDB().Exec(createOrderQuery, productID, productID, count, username)
	if err != nil {
		return 0, fmt.Errorf("failed to create order: %v", err)
	}
	id, err := result.LastInsertId()
	if err != nil || id == 0 {
		return 0, fmt.Errorf("failed to get last inserted id: %v", err)
	}
	return int(id), nil
}

func (m *Manager) UpdateOrder(orderID int, status PaidStatus) error {
	result, err := m.GetDB().Exec(updateOrderQuery, status, orderID)
	if err != nil {
		return fmt.Errorf("failed to create order: %v", err)
	}
	affected, err := result.RowsAffected()
	if affected == 0 || err != nil {
		return fmt.Errorf("failed to update order %d, may be order not found", orderID)
	}
	if status == Cancelled {
		_, err = m.GetDB().Exec(revertProductsQuery, orderID, orderID)
		if err != nil {
			return fmt.Errorf("failed to create order: %v", err)
		}
	}
	return nil
}
