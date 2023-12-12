package types

import "onlinestore/pkg/models"

type BuyRequest struct {
	Count     int `json:"count"`
	Price     int `json:"price"`
	ProductID int `json:"product_id"`
}

type BuyResponse struct {
	OrderID int `json:"order_id"`
	Total   int `json:"total_price"`
}

type CommitOrderRequest struct {
	OrderID int `json:"order_id"`
	Status  int `json:"status"`
}

type CommitOrderResponse struct {
}

type GetOrderRequest struct {
	OrderID  int    `json:"order_id"`
	Username string `json:"username"`
}

type GetOrderResponse struct {
	Order *models.Order `json:"order"`
}

type GetOrdersRequest struct {
	Count int `json:"count"`
}

type GetOrdersResponse struct {
	Orders []*models.Order `json:"orders"`
}

type SubBalanceRequest struct {
	Username string `json:"username"`
	Amount   int    `json:"amount"`
}

type SubBalanceResponse struct {
}
