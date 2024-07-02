package repository

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
	return DBstorage{
		conn: conn,
	}
}

func (s *DBstorage) CreateTable(ctx context.Context) error {
	_, err := s.conn.Exec(ctx, `CREATE TABLE IF NOT EXISTS books
	(
		id UUID DEFAULT uuid_generate_v4() PRIMARY KEY,
    	title TEXT NOT NULL,
    	author TEXT NOT NULL
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
	rows, err := db.conn.Query(ctx, "SELECT id, title, author FROM books")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var books []models.Book
	for rows.Next() {
		var book models.Book
		if err := rows.Scan(&book.ID, &book.Title, &book.Author); err != nil {
			return nil, err
		}
		book.Title = strings.TrimSpace(book.Title)
		book.Author = strings.TrimSpace(book.Author)
		books = append(books, book)
	}
	return books, nil
}

func (db *DBstorage) CreateBook(book models.Book) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	query := "INSERT INTO books (title, author) VALUES ($1, $2) RETURNING id"
	var bookID string
	err := db.conn.QueryRow(ctx, query, book.Title, book.Author).Scan(&bookID)
	if err != nil {
		return "", fmt.Errorf("failed to insert book: %w", err)
	}
	if bookID == "" {
		return "", fmt.Errorf("failed to get bookID after insert")
	}
	return bookID, nil
}