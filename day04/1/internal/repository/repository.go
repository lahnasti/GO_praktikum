package repository

import (
	"errors"

	"github.com/google/uuid"
	"github.com/lahnasti/GO_praktikum/internal/domain/models"
	"github.com/rs/zerolog"
)

type Storage struct {
	db map[string]models.Task
	log *zerolog.Logger
}

func New(zlog *zerolog.Logger) *Storage {
	db :=  make(map[string]models.Task)
	return &Storage{
		db: db,
		log: zlog,
	}
}

func (stor *Storage) AddTask(data models.Task)(string, error) {
	taskID := uuid.New().String()
	data.ID = taskID
	stor.db[taskID] = data
	stor.log.Debug().Any("db", stor.db).Msg("Check db after add task")
	return taskID, nil
}

func (stor *Storage) GetAllTasks()([]models.Task, error) {
	var tasks []models.Task
	for _, task := range stor.db {
		tasks = append(tasks, task)
	}
	stor.log.Debug().Any("db", stor.db).Msg("Check db after get all tasks")
	return tasks, nil
}

func (stor *Storage) GetTaskByID(id string) (models.Task, error) {
	task, exists := stor.db[id]
	if !exists {
		return models.Task{}, errors.New("task not found")
	}
	stor.log.Debug().Any("db", stor.db).Msg("Check db after get task by ID")
	return task, nil
}

func (stor *Storage) UpdateTask(id string, task models.Task) error {
	if _, exists := stor.db[id]; !exists {
		return errors.New("task not found")
	}
	task.ID = id
	stor.db[id] = task
	stor.log.Debug().Any("db", stor.db).Msg("Check db after update task")
	return nil
}

func (stor *Storage) DeleteTask(id string) error {
	if _, exists := stor.db[id]; !exists {
		return errors.New("task not found")
	}
	delete(stor.db, id)
	stor.log.Debug().Any("db", stor.db).Msg("Check db after delete task")
	return nil
}



