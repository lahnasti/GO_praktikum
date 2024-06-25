package repository

import (
	"errors"
	"github.com/lahnasti/GO_praktikum/internal/domain/models"
)

var users = map[string]string{
	"admin": "password",
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