package repository

import "github.com/lahnasti/GO_praktikum/day04/1/internal/domain/models"

// структура для хранения задач
type TaskMap struct {
	list map[string]models.Task
}

func New() *TaskMap {
	list := make(map[string]models.Task)
	return &TaskMap {
	list: list,
	}
}