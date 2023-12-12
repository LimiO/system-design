package types

type ReserveRequest struct {
	ProductID int `json:"product_id"`
	Count     int `json:"count"`
}

type ReserveResponse struct {
	ReserveID int `json:"reserve_id"`
}

type AddCountRequest struct {
	ProductID int `json:"product_id"`
	Count     int `json:"count"`
}

type GetCountRequest struct {
	ProductID int `json:"product_id"`
}

type GetCountResponse struct {
	Count int `json:"count"`
}

type CommitRequest struct {
	ReserveID int `json:"reserve_id"`
	Status    int `json:"status"`
}
