package models

type Book struct {
	ID string `json:"id"`
	Title string `json:"title" validate:"required"`
	Author string `json:"author" validate:"required"`
}