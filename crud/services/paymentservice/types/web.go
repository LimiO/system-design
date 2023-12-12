package types

type AddBalanceRequest struct {
	Username string `json:"username"`
	Amount   int    `json:"amount"`
}

type AddBalanceResponse struct {
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

type ReserveBalanceRequest struct {
	Username string `json:"username"`
	Amount   int    `json:"amount"`
}

type ReserveBalanceResponse struct {
	ReserveID int `json:"reserve_id"`
}

type CommitRequest struct {
	ReserveID int `json:"reserve_id"`
	Status    int `json:"status"`
}
