package models

import "github.com/google/uuid"

var (
	AllowedStatus = []string{"New", "End", "In progress"}
	DefaultStatus = "New"
)
type Task struct {
	ID          uuid.UUID `json:"id"`
	Title       string    `json:"title" validate:"notblank"`
	Description string    `json:"description" validate:"notblank"`
	Status      string    `json:"status" validate:"status"`
}

