package model

//Session ...
type Session struct {
	ID           int    `json:"id"`
	UserID       string `json:"userid"`
	RefrashToken string `json:"-"`
	Useragent    string `json:"useragent"`
}
