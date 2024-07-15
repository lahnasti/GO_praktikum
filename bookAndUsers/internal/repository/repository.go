package repository

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/lahnasti/GO_praktikum/internal/models"
	"golang.org/x/crypto/bcrypt"

	"errors"
)

type Repository struct {
	booksDB map[string]models.Book
	usersDB map[string]models.User
}

func New() *Repository {
	booksDB := make(map[string]models.Book)
	usersDB := make(map[string]models.User)
	return &Repository{
		booksDB: booksDB,
		usersDB: usersDB,
	}
}

func (stor *Repository) AddUser(data models.User) (string, error) {
	// Хеширование пароля
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(data.Password), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("failed to hash password: %w", err)
	}
	data.Password = string(hashedPassword)
	userID := uuid.New().String()
	data.ID = userID
	stor.usersDB[userID] = data
	return userID, nil
}

func (stor *Repository) GetUsers() ([]models.User, error) {
	var users []models.User
	for _, user := range stor.usersDB {
		users = append(users, user)
	}
	return users, nil
}

func (stor *Repository) GetUserByID(id string) (models.User, error) {
	user, exists := stor.usersDB[id]
	if !exists {
		return models.User{}, errors.New("user not found")
	}
	return user, nil
}

func (stor *Repository) UpdateUser(id string, user models.User) error {
	if _, exists := stor.usersDB[id]; !exists {
		return errors.New("user not found")
	}
	user.ID = id
	stor.usersDB[id] = user
	return nil

}

func (stor *Repository) DeleteUser(id string) error {
	if _, exists := stor.usersDB[id]; !exists {
		return errors.New("user not found")
	}
	delete(stor.usersDB, id)
	return nil
}

func (stor *Repository) AddMultipleUsers(users []models.User) ([]string, error) {
	var ids []string
	for _, user := range users {
		userID := uuid.New().String()
		user.ID = userID
		stor.usersDB[userID] = user
		ids = append(ids, userID)
	}
	return ids, nil
}

func (stor *Repository) FindUserByEmail(email string) (models.User, error) {

	for _, user := range stor.usersDB {
		if user.Email == email {
			return user, nil
		}
	}
	return models.User{}, fmt.Errorf("user not found")
}

func (stor *Repository) GetBooks() ([]models.Book, error) {
	var books []models.Book
	for _, book := range stor.booksDB {
		books = append(books, book)
	}
	//stor.log.Debug().Any("db", stor.db).Msg("Check db after get all books")
	return books, nil
}

func (stor *Repository) CreateBook(data models.Book) (string, error) {
	bookID := uuid.New().String()
	data.BID = bookID
	stor.booksDB[bookID] = data
	//stor.log.Debug().Any("db", stor.db).Msg("Check db after add book")
	return bookID, nil
}

func (stor *Repository) CreateMultipleBooks(data []models.Book) ([]string, error) {
	var books []string
	for _, book := range data {
		bookID := uuid.New().String()
		book.BID = bookID
		stor.booksDB[bookID] = book
		books = append(books, bookID)
	}
	return books, nil
}

/*func (stor *Repository) GetBookByID(id string)(models.Book, error) {
	book, exists := stor.db[id]
	if !exists {
		return models.Book{}, errors.New("book not found")
	}
	return book, nil
}*/
