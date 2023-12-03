package web

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
