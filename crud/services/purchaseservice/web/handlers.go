package web

import (
	"fmt"
	"net/http"

	"onlinestore/internal/jwt"
	"onlinestore/pkg/web"
	courier "onlinestore/services/courier/pkg/client"
	payment "onlinestore/services/paymentservice/pkg/client"
	"onlinestore/services/purchaseservice/db"
	"onlinestore/services/purchaseservice/types"
	stock "onlinestore/services/stock/pkg/client"
)

type HandlerManager struct {
	dbManager *db.Manager

	tokenManager  *jwt.TokenManager
	paymentClient *payment.Client
	courClient    *courier.Client
	stockClient   *stock.Client
}

func NewHandlerManager(jwtSecret string, PaymentAddr string, CourAddr string, StockAddr string) (*HandlerManager, error) {
	tokenManager := jwt.NewTokenManager(jwtSecret)
	dbManager, err := db.NewManager()
	if err != nil {
		return nil, fmt.Errorf("failed to create db manager: %v", err)
	}
	paymentClient := payment.NewClient(PaymentAddr)
	courClient := courier.NewClient(CourAddr)
	stockClient := stock.NewClient(StockAddr)
	return &HandlerManager{
		dbManager:     dbManager,
		tokenManager:  tokenManager,
		paymentClient: paymentClient,
		courClient:    courClient,
		stockClient:   stockClient,
	}, nil
}

func (h *HandlerManager) Buy(w http.ResponseWriter, r *http.Request) {
	req, err := web.DecodeHttpBody[types.BuyRequest](r.Body)
	if err != nil {
		web.WriteBadRequest(w, err.Error())
		return
	}
	username := web.GetLogin(r.Context())

	commit := func(w http.ResponseWriter, rollbackPayID int, rollbackStockID int, status int) {
		if err = h.stockClient.Commit(rollbackStockID, status, r.Header); err != nil {
			web.WriteError(w, err.Error(), http.StatusInternalServerError)
			return
		}
		if err = h.paymentClient.Commit(rollbackPayID, status, r.Header); err != nil {
			web.WriteError(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	rollback := func(w http.ResponseWriter, orderID int, rollbackPayID int, rollbackStockID int) {
		if upderr := h.dbManager.UpdateOrder(orderID, db.Cancelled); upderr != nil {
			web.WriteError(w, upderr.Error(), http.StatusInternalServerError)
		}
		commit(w, rollbackPayID, rollbackStockID, 2)
	}

	orderID, err := h.dbManager.CreateOrder(req.ProductID, req.Count, req.Price, username)
	if err != nil {
		web.WriteBadRequest(w, err.Error())
		return
	}

	BalanceReserveID, err := h.paymentClient.ReserveBalance(web.GetLogin(r.Context()), req.Count*req.Price, r.Header)
	if err != nil {
		web.WriteBadRequest(w, err.Error())
		rollback(w, orderID, -1, -1)
		return
	}

	StockReserveID, err := h.stockClient.ReserveCount(req.ProductID, req.Count, r.Header)
	if err != nil {
		web.WriteBadRequest(w, err.Error())
		rollback(w, orderID, BalanceReserveID, StockReserveID)
		return
	}

	if err = h.courClient.ReserveCourier(r.Header); err != nil {
		web.WriteBadRequest(w, err.Error())
		rollback(w, orderID, BalanceReserveID, StockReserveID)
		return
	}

	commit(w, BalanceReserveID, StockReserveID, 1)
	if err = h.dbManager.UpdateOrder(orderID, db.Paid); err != nil {
		return
	}

	web.WriteData(w, types.BuyResponse{OrderID: orderID, Total: req.Count * req.Price})
	return
}

func (h *HandlerManager) CommitOrder(w http.ResponseWriter, r *http.Request) {
	req, err := web.DecodeHttpBody[types.CommitOrderRequest](r.Body)
	if err != nil {
		web.WriteBadRequest(w, err.Error())
		return
	}
	err = h.dbManager.UpdateOrder(req.OrderID, db.PaidStatus(req.Status))
	if err != nil {
		web.WriteBadRequest(w, err.Error())
		return
	}

	web.WriteData(w, types.CommitOrderResponse{})
	return
}

func (h *HandlerManager) GetOrder(w http.ResponseWriter, r *http.Request) {
	req, err := web.DecodeHttpBody[types.GetOrderRequest](r.Body)
	if err != nil {
		web.WriteBadRequest(w, err.Error())
		return
	}

	order, err := h.dbManager.GetOrder(req.OrderID)
	if err != nil {
		web.WriteBadRequest(w, err.Error())
		return
	}
	if order == nil {
		web.WriteError(w, fmt.Sprintf("order not found"), http.StatusNotFound)
		return
	}
	web.WriteData(w, &types.GetOrderResponse{Order: order})
}

func (h *HandlerManager) GetOrders(w http.ResponseWriter, r *http.Request) {
	req, err := web.DecodeHttpBody[types.GetOrdersRequest](r.Body)
	if err != nil {
		web.WriteBadRequest(w, err.Error())
		return
	}

	if req.Count == 0 {
		req.Count = 50
	}
	orders, err := h.dbManager.GetOrders(web.GetLogin(r.Context()), req.Count)
	if err != nil {
		web.WriteBadRequest(w, err.Error())
		return
	}
	web.WriteData(w, &types.GetOrdersResponse{Orders: orders})
}
