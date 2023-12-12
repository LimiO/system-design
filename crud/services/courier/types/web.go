package types

type ReserveCourierResponse struct {
	Username string `json:"username"`
}

type UnreserveCourierRequest struct {
	Username string `json:"username"`
}

type GetCourierRequest struct {
	Username string `json:"username"`
}
