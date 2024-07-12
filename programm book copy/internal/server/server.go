package server

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/lahnasti/GO_praktikum/internal/models"
	"golang.org/x/crypto/bcrypt"
)

type Repository interface {
	AddUser(models.User) (string, error)
	GetUserByID(id string) (models.User, error)
	GetUsers() ([]models.User, error)
	UpdateUser(id string, user models.User) error
	DeleteUser(id string) error
	AddMultipleUsers([]models.User) ([]string, error)
	FindUserByEmail(email string) (models.User, error)

	CreateBook(models.Book) (string, error)
	GetBooks() ([]models.Book, error)
	CreateMultipleBooks([]models.Book) ([]string, error)
	GetBooksByUser(id string) ([]models.Book, error)
}

type Server struct {
	BooksDB Repository
	UsersDB Repository
	Valid   *validator.Validate
}

func (s *Server) GetUsersHandler(ctx *gin.Context) {
	users, err := s.UsersDB.GetUsers()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "List users", "users": users})
}

func (s *Server) RegisterUserHandler(ctx *gin.Context) {
	var user models.User
	err := ctx.ShouldBindBodyWithJSON(&user)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Invalid params", "error": err.Error()})
		return
	}

	err = s.Valid.Struct(user)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Data has not been validated", "error": err.Error()})
		return
	}

	// Хеширование пароля
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to hash password", "error": err.Error()})
		return
	}

	user.Password = string(hashedPassword)

	userID, err := s.UsersDB.AddUser(user)
	if err != nil {
		// Проверка на ошибку доступа к таблице users
		// if strings.Contains(err.Error(), "permission denied for table users") {
		//    ctx.JSON(http.StatusForbidden, gin.H{"message": "Permission denied for table users"})
		//} else {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to save user"})
		return
	}

	ctx.JSON(200, gin.H{"message": "User successfully registered", "user_id": userID})

}

func (s *Server) GetUserByIDHandler(ctx *gin.Context) {
	id := ctx.Param("id")
	user, err := s.UsersDB.GetUserByID(id)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "User retrieved", "user": user})
}

func (s *Server) UpdateUserHandler(ctx *gin.Context) {
	var user models.User
	if err := ctx.ShouldBindJSON(&user); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	id := ctx.Param("id")
	user.ID = id
	err := s.UsersDB.UpdateUser(id, user)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "User updated", "user": user})
}

func (s *Server) DeleteUserHandler(ctx *gin.Context) {
	id := ctx.Param("id")
	err := s.UsersDB.DeleteUser(id)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(200, gin.H{"message": "User deleted", "user_id": id})
}

func (s *Server) RegisterMultipleUsersHadler(ctx *gin.Context) {
	var users []models.User

	err := ctx.ShouldBindBodyWithJSON(&users)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Invalid params", "error": err.Error()})
		return
	}

	for i := range users {
		user := &users[i] // Получаем указатель на каждого пользователя

		err = s.Valid.Struct(user)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"message": "Data has not been validated", "error": err.Error()})
			return
		}
		// Хеширование пароля
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to hash password", "error": err.Error()})
			return
		}

		user.Password = string(hashedPassword)

	}

	userID, err := s.UsersDB.AddMultipleUsers(users)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to save users", "error": err.Error()})
		return
	}

	ctx.JSON(200, gin.H{"message": "Users successfully registered", "users_ID": userID})
}

// Обработчик для входа пользователя, который проверяет учетные данные и генерирует JWT токен:
func (s *Server) LoginHandler(ctx *gin.Context) {
	var user models.User
	if err := ctx.ShouldBindJSON(&user); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Проверка учетных данных пользователя
	dbUser, err := s.UsersDB.FindUserByEmail(user.Email)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials: user not found"})
		return
	}

	fmt.Printf("User from DB: %+v\n", dbUser)

	// Проверка пароля
	err = bcrypt.CompareHashAndPassword([]byte(dbUser.Password), []byte(user.Password))
	if err != nil {
		// Логирование пароля и хеша для отладки
		fmt.Printf("Provided password: %s\n", user.Password)
		fmt.Printf("Stored hash: %s\n", dbUser.Password)
		fmt.Printf("Error: %v\n", err)
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials: incorrect password"})
		return
	}

	// Если учетные данные верны, генерируем JWT токен
	tokenString, err := GenerateJWT(dbUser)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"token": tokenString})
}

func (s *Server) GetBooksHandler(ctx *gin.Context) {
	books, err := s.BooksDB.GetBooks()
	if err != nil {
		//s.log.Error().Err(err).Msg("Failed inquiry")
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(200, gin.H{"message": "Storage books", "books": books})
}

func (s *Server) CreateBookHandler(ctx *gin.Context) {
	token := ctx.GetHeader("Authorization")
	if token == "" {
		ctx.JSON(http.StatusUnauthorized, gin.H{"message": "Authorization header missing"})
		return
	}
	token = strings.TrimPrefix(token, "Bearer ")
	id, err := ValidateJWT(token)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "bad auth token", "error": err.Error()})
		return
	}

	var book models.Book
	err = ctx.ShouldBindBodyWithJSON(&book)
	if err != nil {
		//s.log.Error().Err(err).Msg("Failed unmarshal body")
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Invalid params", "error": err.Error()})
		return
	}
	err = s.Valid.Struct(book)
	if err != nil {
		//s.log.Error().Err(err).Msg("Failed validation")
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Data has not been validated", "error": err.Error()})
		return
	}

	book.ID = id
	bookID, err := s.BooksDB.CreateBook(book)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to save book"})
		return
	}
	ctx.JSON(200, gin.H{"message": "Book successfully created", "book_id": bookID})

}

func (s *Server) CreateMultipleBooksHandler(ctx *gin.Context) {
	token := ctx.GetHeader("Authorization")
	if token == "" {
		ctx.JSON(http.StatusUnauthorized, gin.H{"message": "Authorization header missing"})
		return
	}
	token = strings.TrimPrefix(token, "Bearer ")
	id, err := ValidateJWT(token)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "bad auth token", "error": err.Error()})
		return
	}

	var books []models.Book

	err = ctx.ShouldBindBodyWithJSON(&books)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Invalid params", "error": err.Error()})
	}

	for i, book := range books {
		books[i].ID = id // Устанавливаем userID для каждой книги
		err = s.Valid.Struct(book)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"message": "Data has not been validated", "error": err.Error()})
			return
		}
	}

	bookID, err := s.BooksDB.CreateMultipleBooks(books)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to save books", "error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "Books successfully registered", "book_id": bookID})
}

func (s *Server) GetBooksByUserHandler(ctx *gin.Context) {
	token := ctx.GetHeader("Authorization")
	if token == "" {
		ctx.JSON(http.StatusUnauthorized, gin.H{"message": "Authorization header missing"})
		return
	}

	token = strings.TrimPrefix(token, "Bearer ")
	id, err := ValidateJWT(token)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"message": "Invalid token", "error": err.Error()})
		return
	}
	books, err := s.BooksDB.GetBooksByUser(id)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"message": "Books not found", "error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "Books retrieved", "books": books})
}

/*func (s *Server) GetBookByIDHandler(ctx *gin.Context) {
	param := ctx.Query("id")
	log.Println("Param from url - " + param)
	id, err := strconv.Atoi(param)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Invalid argument", "error": err.Error()})
		return
	}
	log.Printf("ID - %v", id)
	book, err := s.Db.GetBookByID(id)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(200, gin.H{"message": "Book retrieved", "book": book})
}*/

var jwtSecret = []byte("secure_jwt")

type Claims struct {
	User models.User `json:"user"`
	jwt.StandardClaims
}

// Функция для генерации JWT токена для пользователя:
func GenerateJWT(user models.User) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":    user.ID,
		"name":  user.Name,
		"email": user.Email,
		"exp":   time.Now().Add(time.Hour * 72).Unix(),
	})

	tokenString, err := token.SignedString(jwtSecret)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func ValidateJWT(tokenString string) (string, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(jwtSecret), nil
	})

	if err != nil || !token.Valid {
		return "", fmt.Errorf("invalid token: %v", err)
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return "", fmt.Errorf("invalid token claims")
	}

	userID, ok := claims["id"].(string)
	if !ok {
		return "", fmt.Errorf("user id not found in token")
	}

	return userID, nil
}

// Создание JWTAuthMiddleware для проверки токена
func (s *Server) JWTAuthMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		authHeader := ctx.GetHeader("Authorization")
		if authHeader == "" {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header missing"})
			ctx.Abort()
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return jwtSecret, nil
		})

		if err != nil || !token.Valid {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			ctx.Abort()
			return
		}

		ctx.Next()
	}
}
