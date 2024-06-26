package models

type Task struct {
	ID string
	Title string `json:"title" validate:"required"`
	Description string `json:"description" validate:"required"`
}