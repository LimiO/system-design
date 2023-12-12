package types

import "onlinestore/pkg/models"

type PostUserResponse struct {
}

type SubBalanceRequest struct {
	Username string `json:"username"`
	Amount   int    `json:"amount"`
}

type SubBalanceResponse struct {
}

type GetBalanceRequest struct {
	Username string `json:"username"`
}

type GetBalanceResponse struct {
	Balance int `json:"balance"`
}

type AddBalanceRequest struct {
	Username string `json:"username"`
	Amount   int    `json:"amount"`
}

type AddBalanceResponse struct {
}

type TokenRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type TokenResponse struct {
	Token string `json:"token"`
}

type BuyRequest struct {
	Price     int `json:"price"`
	ProductID int `json:"product_id"`
	Count     int `json:"count"`
}

type BuyResponse struct {
	OrderID int `json:"order_id"`
	Total   int `json:"total_price"`
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
