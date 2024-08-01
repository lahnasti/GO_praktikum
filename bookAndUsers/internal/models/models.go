package models

type Book struct {
	BID    int `json:"bId"`
	Title  string `json:"title" validate:"required"`
	Author string `json:"author" validate:"required"`
	UID int `json:"uId"`
}

type User struct {
	UID       int `json:"uId"`
	Name     string `json:"name" validate:"required"`
	Login    string `json:"login" validate:"required"`
	Password string `json:"password" validate:"required"`
}
