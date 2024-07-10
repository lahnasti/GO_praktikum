package booksrepo

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/lahnasti/GO_praktikum/internal/models"
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
	_, err := s.conn.Exec(ctx, `CREATE TABLE IF NOT EXISTS books
	(
		bid UUID DEFAULT uuid_generate_v4() PRIMARY KEY,
    	title TEXT NOT NULL,
    	author TEXT NOT NULL,
		user_id UUID NOT NULL,
        FOREIGN KEY (user_id) REFERENCES users(id)
	);
	`)
	if err != nil {
		return fmt.Errorf("failed to create table: %v", err)
	}
	log.Println("Table 'books' created succesfully")
	return nil

}

func (db *DBstorage) GetBooks() ([]models.Book, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	rows, err := db.conn.Query(ctx, "SELECT bid, title, author, user_id FROM books")
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
	query := "INSERT INTO books (title, author, user_id) VALUES ($1, $2, $3) RETURNING bid"
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

	query := "INSERT INTO books (title, author, user_id) VALUES ($1, $2, $3) RETURNING bid"

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

func (db *DBstorage) GetBooksByUser(user_id string) ([]models.Book, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	rows, err := db.conn.Query(ctx, "SELECT bid, title, author FROM books WHERE user_id=$1", user_id)
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

/*func (db *DBstorage) GetBookByID(id string) (models.Book, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	row := db.conn.QueryRow(ctx, "SELECT title, author FROM books WHERE id=$1", id)
	var book models.Book
	if err := row.Scan(&book.Title, &book.Author); err != nil {
		return models.Book{}, err
	}
	return book, nil
}*/
