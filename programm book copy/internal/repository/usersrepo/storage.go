package usersrepo

import (
	"context"
	"fmt"
	"strings"
	"time"
	"log"

	"github.com/jackc/pgx/v5"
	"github.com/lahnasti/GO_praktikum/internal/models"
)

type DBstorage struct {
	conn *pgx.Conn
}

func NewDB(conn *pgx.Conn) DBstorage {
	return DBstorage {
		conn: conn,
	}
}

func (s *DBstorage) CreateTable(ctx context.Context) error {
	_, err := s.conn.Exec(ctx, `CREATE TABLE IF NOT EXISTS users
	(
		id UUID DEFAULT uuid_generate_v4() PRIMARY KEY,
    	name TEXT NOT NULL,
    	email TEXT UNIQUE NOT NULL,
		password TEXT NOT NULL
	);
	`)
	if err != nil {
		return fmt.Errorf("failed to create table: %v", err)
	}
	log.Println("Table 'users' created succesfully")
	return nil

}

func (db *DBstorage) GetUsers() ([]models.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	rows, err := db.conn.Query(ctx, "SELECT id, name, email, password FROM users")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var users []models.User
	for rows.Next() {
		var user models.User
		if err := rows.Scan(&user.ID, &user.Name, &user.Email, &user.Password); err != nil {
			return nil, err
		}
		user.Name = strings.TrimSpace(user.Name)
		user.Email = strings.TrimSpace(user.Email)
		user.Password = strings.TrimSpace(user.Password)
		users = append(users, user)
	}
	return users, nil
}

func (db *DBstorage) GetUserByID(id string) (models.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	row := db.conn.QueryRow(ctx, "SELECT name, email, password FROM users WHERE id=$1", id)
	var user models.User
	if err := row.Scan(&user.Name, &user.Email, &user.Password); err != nil {
		return models.User{}, err
	}
	return user, nil
}

func (db *DBstorage) AddUser(user models.User) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	query := "INSERT INTO users (name, email, password) VALUES ($1, $2, $3) RETURNING id"
	var userID string
	err := db.conn.QueryRow(ctx, query, user.Name, user.Email, user.Password).Scan(&userID)
	if err != nil {
		return "", fmt.Errorf("failed to insert user: %w", err)
	}
	// Проверка, что userID не пустой
	if userID == "" {
		return "", fmt.Errorf("failed to get userID after insert")
	}
	return userID, nil
}

func (db *DBstorage) UpdateUser(id string, user models.User) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	_, err := db.conn.Exec(ctx, "UPDATE users SET name=$1, email=$2, password=$3 WHERE id=$4", user.Name, user.Email, user.Password, id)
	if err != nil {
		return fmt.Errorf("update user failed: %w", err)
	}
	return nil
}

func (db *DBstorage) DeleteUser(id string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	_, err := db.conn.Exec(ctx, "DELETE FROM users WHERE id=$1", id)
	if err != nil {
		return fmt.Errorf("delete user failed: %w", err)
	}
	return nil
}

func (db *DBstorage) AddMultipleUsers(users []models.User) ([]string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	//открытие транзакции 
	tx, err := db.conn.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(ctx)

    query := "INSERT INTO users (name, email, password) VALUES ($1, $2, $3) RETURNING id"
	//принимает слайс пользователей
	allUsers := make([]string, len(users))
	
	//для каждого ползователя вып запрос на вставку с возвратом ид в виде строки
	// который сохраняется в слайс 
	for i, user := range users {
		var userID string
		err := tx.QueryRow(ctx, query, user.Name, user.Email, user.Password).Scan(&userID)
		if err != nil {
			return nil, err
		}
		allUsers[i] = userID
	}
	// если все ок, транзакция фиксируется с ->
	err = tx.Commit(ctx)
	if err != nil {
		return nil, err
	} // возвращение ид всех доб пользователей в виде строк
	return allUsers, nil
}
	
	




