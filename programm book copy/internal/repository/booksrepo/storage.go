package booksrepo

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/lahnasti/GO_praktikum/internal/models"
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
	rows, err := db.conn.Query(ctx, "SELECT bid, title, author, iduser FROM books")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var books []models.Book
	for rows.Next() {
		var book models.Book
		if err := rows.Scan(&book.BID, &book.Title, &book.Author, &book.IDuser); err != nil {
			return nil, err
		}
		book.Title = strings.TrimSpace(book.Title)
		book.Author = strings.TrimSpace(book.Author)
		book.IDuser = strings.TrimSpace(book.IDuser)
		books = append(books, book)
	}
	return books, nil
}

func (db *DBstorage) CreateBook(book models.Book) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	query := "INSERT INTO books (title, author) VALUES ($1, $2) RETURNING bid"
	var bookBID string
	err := db.conn.QueryRow(ctx, query, book.Title, book.Author).Scan(&bookBID)
	if err != nil {
		return "", fmt.Errorf("failed to insert book: %w", err)
	}
	if bookBID == "" {
		return "", fmt.Errorf("failed to get bookID after insert")
	}
	return bookBID, nil
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
