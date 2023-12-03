package models

type UnhashedPasswordInfo struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type PasswordInfo struct {
	Username string `db:"username"`
	Passhash string `db:"passhash"`
}

func (p *PasswordInfo) GetUsername() string {
	return p.Username
}

func (p *PasswordInfo) GetPasshash() string {
	return p.Passhash
}
