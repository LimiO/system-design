package models

type Order struct {
	OrderID   int    `db:"order_id" json:"order_id"`
	ProductID int    `db:"product_id" json:"product_id"`
	Price     int    `db:"price" json:"price"`
	Count     int    `db:"count" json:"count"`
	Username  string `db:"username" json:"username"`
	Paid      bool   `db:"paid" json:"paid"`
}

type Product struct {
	ProductID int `db:"product_id" json:"product_id"`
	Price     int `db:"price" json:"price"`
	Count     int `db:"count" json:"count"`
}
