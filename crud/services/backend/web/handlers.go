package web

import (
	"fmt"
	"io"
	"net/http"

	"onlinestore/internal/jwt"
	"onlinestore/pkg/models"
	"onlinestore/pkg/web"
	"onlinestore/services/backend/types"

	auth "onlinestore/services/authorizationservice/pkg/client"
	payment "onlinestore/services/paymentservice/pkg/client"
	purchases "onlinestore/services/purchaseservice/pkg/client"
	user "onlinestore/services/userservice/pkg/client"
)

type HandlerManager struct {
	tokenManager    *jwt.TokenManager
	paymentClient   *payment.Client
	userClient      *user.Client
	authClient      *auth.Client
	purchasesClient *purchases.Client
}

func (h *HandlerManager) RetranslateRequest(
	w http.ResponseWriter,
	r *http.Request,
	callback func(data []byte, header http.Header) (*http.Response, error),
) {
	bodyBytes, _ := io.ReadAll(r.Body)
	resp, err := callback(bodyBytes, r.Header)
	if err != nil {
		web.WriteError(w, err.Error(), resp.StatusCode)
		return
	}
	respBody, _ := io.ReadAll(resp.Body)
	w.WriteHeader(resp.StatusCode)
	_, _ = w.Write(respBody)
	return
}

func NewHandlerManager(jwtSecret string, PaymentAddr, AuthorizationAddr, PurchasesAddr, UserAddr string) (*HandlerManager, error) {
	tokenManager := jwt.NewTokenManager(jwtSecret)
	paymentClient := payment.NewClient(PaymentAddr)
	userClient := user.NewClient(UserAddr)
	authClient := auth.NewClient(AuthorizationAddr)
	purchasesClient := purchases.NewClient(PurchasesAddr)
	return &HandlerManager{
		tokenManager:    tokenManager,
		paymentClient:   paymentClient,
		userClient:      userClient,
		authClient:      authClient,
		purchasesClient: purchasesClient,
	}, nil
}

func (h *HandlerManager) AddBalance(w http.ResponseWriter, r *http.Request) {
	login := web.GetLogin(r.Context())

	req, err := web.DecodeHttpBody[types.AddBalanceRequest](r.Body)
	if err != nil {
		web.WriteBadRequest(w, err.Error())
		return
	}

	err = h.paymentClient.AddBalance(login, req.Amount, r.Header)
	if err != nil {
		web.WriteError(w, fmt.Sprintf("failed to add balance: %v", err), http.StatusBadRequest)
		return
	}
	web.WriteData(w, types.AddBalanceResponse{})
}

func (h *HandlerManager) SubBalance(w http.ResponseWriter, r *http.Request) {
	login := web.GetLogin(r.Context())

	req, err := web.DecodeHttpBody[types.SubBalanceRequest](r.Body)
	if err != nil {
		web.WriteBadRequest(w, err.Error())
		return
	}

	err = h.paymentClient.SubBalance(login, req.Amount, r.Header)
	if err != nil {
		web.WriteError(w, fmt.Sprintf("failed to sub balance: %v", err), http.StatusBadRequest)
		return
	}
	web.WriteData(w, types.SubBalanceResponse{})
}

func (h *HandlerManager) GetBalance(w http.ResponseWriter, r *http.Request) {
	login := web.GetLogin(r.Context())
	balance, err := h.paymentClient.GetBalance(login, r.Header)
	if err != nil {
		web.WriteError(w, fmt.Sprintf("failed to get balance: %v", err), http.StatusBadRequest)
		return
	}
	web.WriteData(w, types.GetBalanceResponse{
		Balance: balance,
	})
}

func (h *HandlerManager) GetUser(w http.ResponseWriter, r *http.Request) {
	login := web.GetLogin(r.Context())

	u, err := h.userClient.GetUser(login, r.Header)
	if err != nil {
		web.WriteError(w, fmt.Sprintf("failed to get user %q: %v", login, err), http.StatusBadRequest)
		return
	}
	web.WriteData(w, u)
}

func (h *HandlerManager) PostUser(w http.ResponseWriter, r *http.Request) {
	u, err := web.DecodeHttpBody[models.User](r.Body)
	if err != nil {
		web.WriteBadRequest(w, err.Error())
		return
	}

	token, err := h.authClient.Register(u.Username, u.Password, r.Header)
	if err != nil {
		web.WriteError(w, fmt.Sprintf("failed to register: %v", err), http.StatusBadRequest)
		return
	}

	r.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))

	err = h.userClient.PostUser(u, r.Header)
	if err != nil {
		if unrerr := h.authClient.Unregister(u.Username, u.Password, r.Header); unrerr != nil {
			web.WriteError(w, fmt.Sprintf("failed to unregister user: %v", unrerr), http.StatusInternalServerError)
			return
		}
		web.WriteError(w, fmt.Sprintf("failed to make user: %v", err), http.StatusBadRequest)
		return
	}

	web.WriteData(w, types.TokenResponse{Token: token})
}

func (h *HandlerManager) PutUser(w http.ResponseWriter, r *http.Request) {
	u, err := web.DecodeHttpBody[models.User](r.Body)
	if err != nil {
		web.WriteBadRequest(w, err.Error())
		return
	}
	err = h.userClient.PutUser(u, r.Header)
	if err != nil {
		web.WriteError(w, fmt.Sprintf("failed to put user: %v", err), http.StatusBadRequest)
		return
	}
	web.WriteStatusOK(w)
}

func (h *HandlerManager) DeleteUser(w http.ResponseWriter, r *http.Request) {
	login := web.GetLogin(r.Context())
	err := h.userClient.DeleteUser(login, r.Header)
	if err != nil {
		web.WriteError(w, fmt.Sprintf("failed to delete user: %v", err), http.StatusBadRequest)
		return
	}
	web.WriteStatusOK(w)
}

func (h *HandlerManager) GetToken(w http.ResponseWriter, r *http.Request) {
	req, err := web.DecodeHttpBody[types.TokenRequest](r.Body)
	if err != nil {
		web.WriteBadRequest(w, err.Error())
		return
	}

	token, err := h.authClient.GetToken(req.Username, req.Password, r.Header)
	if err != nil {
		web.WriteError(w, fmt.Sprintf("failed to get token: %v", err), http.StatusBadRequest)
		return
	}
	web.WriteData(w, types.TokenResponse{
		Token: token,
	})
}

func (h *HandlerManager) GetOrder(w http.ResponseWriter, r *http.Request) {
	req, err := web.DecodeHttpBody[types.GetOrderRequest](r.Body)
	if err != nil {
		web.WriteBadRequest(w, err.Error())
		return
	}

	order, err := h.purchasesClient.GetOrder(req.Username, req.OrderID, r.Header)
	if err != nil {
		web.WriteError(w, fmt.Sprintf("failed to get order: %v", err), http.StatusBadRequest)
		return
	}
	web.WriteData(w, types.GetOrderResponse{
		Order: order,
	})
}

func (h *HandlerManager) GetOrders(w http.ResponseWriter, r *http.Request) {
	req, err := web.DecodeHttpBody[types.GetOrdersRequest](r.Body)
	if err != nil {
		web.WriteBadRequest(w, err.Error())
		return
	}

	orders, err := h.purchasesClient.GetOrders(req.Count, r.Header)
	if err != nil {
		web.WriteError(w, fmt.Sprintf("failed to get orders: %v", err), http.StatusBadRequest)
		return
	}
	web.WriteData(w, types.GetOrdersResponse{
		Orders: orders,
	})
}

func (h *HandlerManager) Buy(w http.ResponseWriter, r *http.Request) {
	buyReq, err := web.DecodeHttpBody[types.BuyRequest](r.Body)
	if err != nil {
		web.WriteBadRequest(w, err.Error())
		return
	}

	orderID, total, err := h.purchasesClient.Buy(buyReq.Count, buyReq.Price, buyReq.ProductID, r.Header)
	if err != nil {
		web.WriteBadRequest(w, err.Error())
		return
	}

	web.WriteData(w, types.BuyResponse{
		OrderID: orderID,
		Total:   total,
	})
}
