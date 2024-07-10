package usersrepo

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/lahnasti/GO_praktikum/internal/models"
	"golang.org/x/crypto/bcrypt"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/rs/zerolog"

)

type DBstorage struct {
	conn *pgx.Conn
}

func NewDB(conn *pgx.Conn) DBstorage {
	return DBstorage{
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
		// Хеширование пароля
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
		if err != nil {
			return nil, err
		}
		var userID string
		err = tx.QueryRow(ctx, query, user.Name, user.Email, string(hashedPassword)).Scan(&userID)
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

// FindUserByEmail ищет пользователя по email в базе данных и возвращает его, если найден.
func (db *DBstorage) FindUserByEmail(email string) (models.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var user models.User
	query := "SELECT id, name, email, password FROM users WHERE email = $1"
	err := db.conn.QueryRow(ctx, query, email).Scan(&user.ID, &user.Name, &user.Email, &user.Password)
	if err != nil {
		return models.User{}, fmt.Errorf("failed to find user by email: %w", err)
	}

	return user, nil
}

func Migrations(dbAddr, migrationsPath string, zlog *zerolog.Logger) error {
	migratePath := fmt.Sprintf("file://%s", migrationsPath)
	m, err := migrate.New(migratePath, dbAddr)
	if err != nil {
		return err
	}
	if err := m.Up(); err != nil {
		if errors.Is(err, migrate.ErrNoChange) {
			zlog.Debug().Msg("No migrations apply")
			return nil
		}
		return err
	}
	zlog.Debug().Msg("Migrate complete")
	return nil
}