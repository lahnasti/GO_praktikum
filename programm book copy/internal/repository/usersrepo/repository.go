package usersrepo

import (
	"github.com/google/uuid"
	"github.com/lahnasti/GO_praktikum/internal/models"

	"errors"
)

type Repository struct {
	db map[string]models.User
}

func New() *Repository {
	db := make(map[string]models.User)
	return &Repository{
		db: db,
	}
}

func (stor *Repository) AddUser(data models.User) (string, error) {
	userID := uuid.New().String()
	data.ID = userID
	stor.db[userID] = data
	return userID, nil
}

func (stor *Repository) GetUsers() ([]models.User, error) {
	var users []models.User
	for _, user := range stor.db {
		users = append(users, user)
	}
	return users, nil
}

func (stor *Repository) GetUserByID(id string) (models.User, error) {
	user, exists := stor.db[id]
	if !exists {
		return models.User{}, errors.New("user not found")
	}
	return user, nil
}

func (stor *Repository) UpdateUser(id string, user models.User) error {
	if _, exists := stor.db[id]; !exists {
		return errors.New("user not found")
	}
	user.ID = id
	stor.db[id] = user
	return nil

}

func (stor *Repository) DeleteUser(id string) error {
	if _, exists := stor.db[id]; !exists {
		return errors.New("user not found")
	}
	delete(stor.db, id)
	return nil
}

func (stor *Repository) AddMultipleUsers(users []models.User) ([]string, error) {
	var ids []string
	for _, user := range users {
		userID := uuid.New().String()
		user.ID = userID
		stor.db[userID] = user
		ids = append(ids, userID)
	}
	return ids, nil
}
