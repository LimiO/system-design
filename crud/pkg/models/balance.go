package models

type BalanceInfo struct {
	Balance  int    `db:"balance" json:"balance"`
	Username string `db:"username" json:"username"`
}
