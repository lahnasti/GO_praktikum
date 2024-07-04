package repository

import (
	"github.com/google/uuid"
	"github.com/lahnasti/GO_praktikum/internal/models"
	"errors"
	//"github.com/rs/zerolog"
)

type Repository struct {
	db map[string]models.Book
	//log *zerolog.Logger
}

func New() *Repository {
	db := make(map[string]models.Book)
	return &Repository{
		db: db,
		//log: zlog,
	}
}


func (stor *Repository) GetBooks() ([]models.Book, error) {
	var books []models.Book
	for _, book := range stor.db {
		books = append(books, book)
	}
	//stor.log.Debug().Any("db", stor.db).Msg("Check db after get all books")
	return books, nil
}

func (stor *Repository) CreateBook(data models.Book) (string, error) {
	bookID := uuid.New().String()
	data.ID = bookID
	stor.db[bookID] = data
	//stor.log.Debug().Any("db", stor.db).Msg("Check db after add book")
	return bookID, nil
}

func (stor *Repository) GetBookByID(id string)(models.Book, error) {
	book, exists := stor.db[id]
	if !exists {
		return models.Book{}, errors.New("book not found")
	}
	return book, nil
}