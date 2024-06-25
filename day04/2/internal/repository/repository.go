package repository

import (
	"github.com/google/uuid"
	"github.com/lahnasti/GO_praktikum/internal/domain/models"

	"errors"
)

type Storage struct {
	db map[string]models.User
}

func New() *Storage {
	db := make(map[string]models.User)
	return &Storage{
		db: db,
	}
}

func (stor *Storage) AdduUser(data models.User)(string, error) {
	userID := uuid.New().String()
	data.ID = userID
	stor.db[userID] = data
	return userID, nil
}

func (stor *Storage) GetUsers()([]models.User, error) {
	var users []models.User
	for _, user := range stor.db {
		users = append(users, user)
	}
	return users, nil
}

func (stor *Storage) GetUserByID(id string) (models.User, error) {
	user, exists := stor.db[id]
	if !exists {
		return models.User{}, errors.New("user not found")
	}
	return user, nil
}

func (stor *Storage) UpdateUser (id string, user models.User) error {
	if _, exists := stor.db[id]; !exists {
		return errors.New("user not found")
	}
	user.ID = id
	stor.db[id] = user
	return nil
}

func (stor *Storage) DeleteUser (id string) error {
	if _, exists := stor.db[id]; !exists {
		return errors.New("user not found")
	}
	delete(stor.db, id)
	return nil
}

