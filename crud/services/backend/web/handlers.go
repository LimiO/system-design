package web

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"onlinestore/pkg/models"
	web2 "onlinestore/services/purchaseservice/web"

	"onlinestore/services/backend/clients/auth"
	"onlinestore/services/backend/clients/payment"
	"onlinestore/services/backend/clients/purchases"
	"onlinestore/services/backend/clients/user"

	"onlinestore/internal/jwt"
	"onlinestore/pkg/web"
	authweb "onlinestore/services/authorizationservice/web"
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
	h.RetranslateRequest(w, r, h.paymentClient.AddBalance)
}

func (h *HandlerManager) SubBalance(w http.ResponseWriter, r *http.Request) {
	h.RetranslateRequest(w, r, h.paymentClient.SubBalance)
}

func (h *HandlerManager) GetBalance(w http.ResponseWriter, r *http.Request) {
	h.RetranslateRequest(w, r, h.paymentClient.GetBalance)
}

func (h *HandlerManager) GetUser(w http.ResponseWriter, r *http.Request) {
	h.RetranslateRequest(w, r, h.userClient.GetUser)
}

func (h *HandlerManager) PostUser(w http.ResponseWriter, r *http.Request) {
	req, err := web.DecodeHttpBody[models.User](r.Body)
	if err != nil {
		web.WriteError(w, err.Error(), http.StatusBadRequest)
		return
	}

	data := map[string]string{
		"username": req.Username,
		"password": req.Password,
	}
	registerBody, _ := json.Marshal(data)

	registerResp, err := h.authClient.Register(registerBody, r.Header)
	if err != nil {
		web.WriteError(w, err.Error(), http.StatusBadRequest)
		return
	}
	if registerResp.StatusCode != http.StatusOK {
		web.WriteError(w, "failed to register", registerResp.StatusCode)
		return
	}
	tokenResp, _ := web.DecodeHttpBody[authweb.TokenResponse](registerResp.Body)

	r.Header.Set("Authorization", fmt.Sprintf("Bearer %s", tokenResp.Token))
	d, _ := json.Marshal(req)
	userResp, err := h.userClient.PostUser(d, r.Header)
	if err != nil {
		web.WriteError(w, err.Error(), http.StatusBadRequest)
		return
	}
	if userResp.StatusCode != http.StatusOK {
		_, _ = h.authClient.Unregister(registerBody, r.Header)
		web.WriteError(w, "failed to make user", userResp.StatusCode)
		return
	}

	web.WriteData(w, tokenResp)
}

func (h *HandlerManager) PutUser(w http.ResponseWriter, r *http.Request) {
	h.RetranslateRequest(w, r, h.userClient.PutUser)
}

func (h *HandlerManager) DeleteUser(w http.ResponseWriter, r *http.Request) {
	h.RetranslateRequest(w, r, h.userClient.DeleteUser)
}

func (h *HandlerManager) GetToken(w http.ResponseWriter, r *http.Request) {
	h.RetranslateRequest(w, r, h.authClient.GetToken)
}

func (h *HandlerManager) GetOrder(w http.ResponseWriter, r *http.Request) {
	h.RetranslateRequest(w, r, h.purchasesClient.GetOrder)
}

func (h *HandlerManager) GetOrders(w http.ResponseWriter, r *http.Request) {
	h.RetranslateRequest(w, r, h.purchasesClient.GetOrders)
}

func (h *HandlerManager) Buy(w http.ResponseWriter, r *http.Request) {
	buyReq, err := web.DecodeHttpBody[web2.BuyRequest](r.Body)
	if err != nil {
		web.WriteError(w, err.Error(), http.StatusBadRequest)
		return
	}
	body, _ := json.Marshal(buyReq)

	purchResp, err := h.purchasesClient.Buy(body, r.Header)
	if err != nil {
		web.WriteError(w, err.Error(), http.StatusBadRequest)
		return
	}
	if purchResp.StatusCode != http.StatusOK {
		purchRespBytes, _ := io.ReadAll(purchResp.Body)
		web.WriteError(w, fmt.Sprintf("failed to buy: %v", string(purchRespBytes)), purchResp.StatusCode)
		return
	}
	purchBuyResp, _ := web.DecodeHttpBody[web2.BuyResponse](purchResp.Body)

	subBalanceReq := &SubBalanceRequest{
		Username: web.GetLogin(r.Context()),
		Amount:   purchBuyResp.Total,
	}
	subBalanceReqData, _ := json.Marshal(subBalanceReq)
	subBalanceResp, err := h.paymentClient.SubBalance(subBalanceReqData, r.Header)
	commitOrderReq := web2.CommitOrderRequest{
		OrderID: purchBuyResp.OrderID,
	}
	if subBalanceResp.StatusCode != http.StatusOK {
		commitOrderReq.Status = 2
		commitOrderReqData, _ := json.Marshal(commitOrderReq)
		_, _ = h.purchasesClient.Commit(commitOrderReqData, r.Header)

		subBalanceRespData, _ := io.ReadAll(subBalanceResp.Body)
		web.WriteError(w, string(subBalanceRespData), subBalanceResp.StatusCode)
		return
	}

	commitOrderReq.Status = 1
	commitOrderReqData, _ := json.Marshal(commitOrderReq)
	_, _ = h.purchasesClient.Commit(commitOrderReqData, r.Header)
	web.WriteData(w, purchBuyResp)
}
