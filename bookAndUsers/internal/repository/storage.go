package repository

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jackc/pgx/v5"
	"github.com/lahnasti/GO_praktikum/internal/models"
	"github.com/rs/zerolog"
	"golang.org/x/crypto/bcrypt"
)

type DBstorage struct {
	conn *pgx.Conn
}

func NewDB(conn *pgx.Conn) (DBstorage, error) {
	return DBstorage{
		conn: conn,
	}, nil
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

func (db *DBstorage) GetBooks() ([]models.Book, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	rows, err := db.conn.Query(ctx, "SELECT bid, title, author, id FROM books")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var books []models.Book
	for rows.Next() {
		var book models.Book
		if err := rows.Scan(&book.BID, &book.Title, &book.Author, &book.ID); err != nil {
			return nil, err
		}
		book.Title = strings.TrimSpace(book.Title)
		book.Author = strings.TrimSpace(book.Author)
		book.ID = strings.TrimSpace(book.ID)
		books = append(books, book)
	}
	return books, nil
}

func (db *DBstorage) CreateBook(book models.Book) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	query := "INSERT INTO books (title, author, id) VALUES ($1, $2, $3) RETURNING bid"
	var bookBID string
	err := db.conn.QueryRow(ctx, query, book.Title, book.Author, book.ID).Scan(&bookBID)
	if err != nil {
		return "", fmt.Errorf("failed to insert book: %w", err)
	}
	if bookBID == "" {
		return "", fmt.Errorf("failed to get bookID after insert")
	}
	return bookBID, nil
}

func (db *DBstorage) CreateMultipleBooks(books []models.Book) ([]string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	tx, err := db.conn.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(ctx)

	query := "INSERT INTO books (title, author, id) VALUES ($1, $2, $3) RETURNING bid"

	allBooks := make([]string, len(books))

	for i, book := range books {
		var bookBID string
		err := tx.QueryRow(ctx, query, book.Title, book.Author, book.ID).Scan(&bookBID)
		if err != nil {
			return nil, err
		}
		allBooks[i] = bookBID
	}
	err = tx.Commit(ctx)
	if err != nil {
		return nil, err
	}
	return allBooks, nil
}

func (db *DBstorage) GetBooksByUser(id string) ([]models.Book, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	rows, err := db.conn.Query(ctx, "SELECT bid, title, author FROM books WHERE id=$1", id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var books []models.Book
	for rows.Next() {
		var book models.Book
		if err := rows.Scan(&book.BID, &book.Title, &book.Author); err != nil {
			return nil, err
		}
		books = append(books, book)

	}
	if err := rows.Err(); err != nil {
		return nil, err

	}
	return books, nil

}

func (db *DBstorage) SetDeleteStatus(bid int) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if _, err := db.conn.Exec(ctx, "UPDATE books SET delete = true WHERE bid = $1"); err != nil {
		return fmt.Errorf("update delete status failed: %w", err)
	}
	return nil
}

func (db *DBstorage) DeleteBooks(bIds []int) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	tx, err := db.conn.Begin(ctx)
	if err != nil {
		return fmt.Errorf("create transaction failed: %w", err)
	}
	defer tx.Rollback(ctx)

	if _, err := tx.Prepare(ctx, "delete book", "DELETE FROM books WHERE bid = $1"); err != nil {
		return fmt.Errorf("create prepare sql str failed: %w", err)
	}

	for _, bid := range bIds {
		if _, err := tx.Exec(ctx, "delete book", bid); err != nil {
			return fmt.Errorf("failed delete book: %w", err)
		}
	}
	return tx.Commit(ctx)
}