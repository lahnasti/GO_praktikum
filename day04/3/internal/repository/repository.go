package repository

import (
	"errors"
	"github.com/lahnasti/GO_praktikum/internal/domain/models"
)

var users = map[string]string{}

func Register(username, password string) error {
	if _, exists := users[username]; exists {
		return errors.New("user already exists")
	}
	users[username] = password
	return nil
}

func Authenticate(username, password string) error {
	if pass, ok := users[username]; ok && pass == password {
		return nil
	}
	return errors.New("invalid credentials")
}

func GetUser(username string) (*models.User, error) {
	if _, ok := users[username]; ok {
		return &models.User{Username: username}, nil
	}
	return nil, errors.New("user not found")
}