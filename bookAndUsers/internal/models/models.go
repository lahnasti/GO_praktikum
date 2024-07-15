package models

type Book struct {
	BID    string `json:"bid"`
	Title  string `json:"title" validate:"required"`
	Author string `json:"author" validate:"required"`
	ID string `json:"id"`
}

type User struct {
	ID       string `json:"id"`
	Name     string `json:"name" validate:"required"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
	//Token    string `json:"token"`
}
