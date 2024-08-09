package server

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/golang-jwt/jwt/v4"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/lahnasti/GO_praktikum/internal/models"
	"github.com/rs/zerolog"
	"golang.org/x/crypto/bcrypt"
)

const jwtSecret = "secure_jwt"

type Claims struct {
	UserID string
	jwt.RegisteredClaims
}

type Repository interface {
	GetAllUsers() ([]models.User, error)
	GetUser(int) (models.User, error)
	GetUserByLogin(string) (models.User, error)
	GetAllBooks() ([]models.Book, error)
	GetBooksByUser(int) ([]models.Book, error)
	AddUser(models.User) (int, error)
	SaveBook(models.Book) error
	SaveBooks([]models.Book, int) error
	UpdateUser(int, models.User) error
	DeleteUser() error
	DeleteBooks() error
	SetDeleteStatus(int) error
}

type Server struct {
	Db         Repository
	ErrorChan  chan error
	deleteChan chan int
	Valid      *validator.Validate
	log        zerolog.Logger
}

func NewServer(ctx context.Context, db Repository, zlog *zerolog.Logger) *Server {
	dChan := make(chan int, 5)
	errChan := make(chan error)
	srv := Server{
		Db:         db,
		deleteChan: dChan,
		ErrorChan:  errChan,
		log:        *zlog,
	}
	go srv.deleter(ctx)
	return &Server{
		Db:         db,
		deleteChan: dChan,
		ErrorChan:  errChan,
	}
}

func (s *Server) GetAllUsersHandler(ctx *gin.Context) {
	users, err := s.Db.GetAllUsers()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "List users", "users": users})
}

func (s *Server) GetUserHandler(c *gin.Context) {
	param := c.Query("uid")
	log.Println("Param from url - " + param)
	uid, err := strconv.Atoi(param)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid argument"})
		return
	}
	log.Printf("UID - %v", uid)
	user, err := s.Db.GetUser(uid)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, user)
}

func (s *Server) RegisterUserHandler(ctx *gin.Context) {
	var user models.User
	err := ctx.ShouldBindBodyWithJSON(&user)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Invalid params", "error": err.Error()})
		return
	}

	/*err = s.Valid.Struct(user)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Data has not been validated", "error": err.Error()})
		return
	}*/

	// Хеширование пароля
	hash, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to hash password", "error": err.Error()})
		return
	}
	user.Password = string(hash)
	uid, err := s.Db.AddUser(user)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) { // Приводим ошибку err к типу PgError
			if !pgerrcode.IsIntegrityConstraintViolation(pgErr.Code) { // Проверяем, входит ли ошибка в 23й класс ошибок Sql
				ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()}) // Если нет, то просто выводим ошибку
				return
			}
			ctx.JSON(http.StatusConflict, gin.H{"error": "login already userd"})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	uidStr := strconv.Itoa(uid)
	token, err := GenerateJWT(uidStr)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.Header("Authorization", token)
	ctx.JSON(200, gin.H{"message": "User successfully registered", "user_id": uid})

}

// Обработчик для входа пользователя, который проверяет учетные данные и генерирует JWT токен:
func (s *Server) LoginHandler(ctx *gin.Context) {
	var user models.User
	if err := ctx.ShouldBindJSON(&user); err != nil {
		s.log.Error().Err(err).Msg("failed parse login data from body")
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Проверка учетных данных пользователя
	userFromDb, err := s.Db.GetUserByLogin(user.Login)
	if err != nil {
		s.log.Error().Err(err).Msg("failed get user by login")
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials: user not found"})
		return
	}

	// Проверка пароля
	err = bcrypt.CompareHashAndPassword([]byte(userFromDb.Password), []byte(user.Password))
	if err != nil {
		s.log.Error().Err(err).Msg("failed get user by login")
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials: incorrect password"})
		return
	}

	// Если учетные данные верны, генерируем JWT токен
	uidStr := strconv.Itoa(userFromDb.UID)
	token, err := GenerateJWT(uidStr)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}
	ctx.Header("Authorization", token)
	ctx.String(http.StatusOK, "User %s was logined", user.Name)
}

func (s *Server) UpdateUserHandler(ctx *gin.Context) {
	var user models.User
	if err := ctx.ShouldBindJSON(&user); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	uid := ctx.Param("id")
	uIdInt, err := strconv.Atoi(uid)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	err = s.Db.UpdateUser(uIdInt, user)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "User updated", "user": user})
}

func (s *Server) DeleteUserHandler(ctx *gin.Context) {
	token := ctx.GetHeader("Authorization")
	_, err := ValidateJWT(token)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Bad auth token", "error": err.Error()})
		return
	}

	uId := ctx.Param("id")
	uIdInt, err := strconv.Atoi(uId)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Bad user_id", "error": err.Error()})
		return
	}

	err = s.Db.SetDeleteStatus(uIdInt)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	s.deleteChan <- uIdInt
	ctx.JSON(http.StatusOK, gin.H{"message": "User deleted", "user_id": uIdInt})
}

func (s *Server) GetAllBooksHandler(ctx *gin.Context) {
	books, err := s.Db.GetAllBooks()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "Storage books", "books": books})
}

func (s *Server) SaveBookHandler(ctx *gin.Context) {
	token := ctx.GetHeader("Authorization")
	uid, err := ValidateJWT(token)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Bad auth token", "error": err.Error()})
		return
	}

	var book models.Book
	if err := ctx.ShouldBindBodyWithJSON(&book); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	uidInt, err := strconv.Atoi(uid)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	book.UID = uidInt
	if err := s.Db.SaveBook(book); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "Book saved"})

}

func (s *Server) SaveBooksHandler(ctx *gin.Context) {
	token := ctx.GetHeader("Authorization")
	if token == "" {
		ctx.JSON(http.StatusUnauthorized, gin.H{"message": "Authorization header missing"})
		return
	}
	uid, err := ValidateJWT(token)
	if err != nil {
		s.log.Error().Err(err).Msg("get uid failed")
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Bad auth token", "error": err.Error()})
		return
	}
	var books []models.Book
	if err := ctx.ShouldBindBodyWithJSON(&books); err != nil {
		s.log.Error().Err(err).Msg("parse body failed")
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	log.Println(books)
	uidInt, err := strconv.Atoi(uid)
	if err != nil {
		s.log.Error().Err(err).Msg("parse uid from str to int failed")
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if err := s.Db.SaveBooks(books, uidInt); err != nil {
		s.log.Error().Err(err).Msg("save books failed")
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "All books saved"})
}

func (s *Server) GetBooksByUserHandler(ctx *gin.Context) {
	token := ctx.GetHeader("Authorization")
	if token == "" {
		ctx.JSON(http.StatusUnauthorized, gin.H{"message": "Authorization header missing"})
		return
	}
	uIdStr, err := ValidateJWT(token)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"message": "Invalid token", "error": err.Error()})
		return
	}
	uId, err := strconv.Atoi(uIdStr)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	books, err := s.Db.GetBooksByUser(uId)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"message": "Books not found", "error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "Books retrieved", "books": books})
}

func (s *Server) DeleteBookHandler(ctx *gin.Context) {
	token := ctx.GetHeader("Authorization")
	_, err := ValidateJWT(token)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Bad auth token", "error": err.Error()})
		return
	}

	bId := ctx.Param("id")
	bIdInt, err := strconv.Atoi(bId)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Bad book_id", "error": err.Error()})
		return
	}

	err = s.Db.SetDeleteStatus(bIdInt)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	s.deleteChan <- bIdInt
	ctx.JSON(http.StatusOK, gin.H{"message": "Book deleted", "book_id": bIdInt})
}

func (s *Server) deleter(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		default:
			if len(s.deleteChan) == 5 {
				for i := 0; i < 5; i++ {
					<-s.deleteChan
				}
				if err := s.Db.DeleteBooks(); err != nil {
					s.ErrorChan <- err
					return
				}
			}
		}
	}
}

// Функция для генерации JWT токена для пользователя:
func GenerateJWT(uid string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 3)),
		},
		UserID: uid,
	})
	key := []byte(jwtSecret)
	tokenStr, err := token.SignedString(key)
	if err != nil {
		return "", err
	}
	return tokenStr, nil
}

func ValidateJWT(tokenStr string) (string, error) {
	claim := &Claims{}

	token, err := jwt.ParseWithClaims(tokenStr, claim, func(t *jwt.Token) (interface{}, error) {
		return []byte(jwtSecret), nil
	})
	if err != nil {
		return "", err
	}

	if !token.Valid {
		return "", fmt.Errorf("invalid token")
	}

	return claim.UserID, nil
}

// Создание JWTAuthMiddleware для проверки токена
/*func (s *Server) JWTAuthMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		authHeader := ctx.GetHeader("Authorization")
		if authHeader == "" {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header missing"})
			ctx.Abort()
			return
		}

		tokenStr := strings.TrimPrefix(authHeader, "Bearer ")
		if tokenStr == authHeader {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header format must be Bearer {token}"})
			ctx.Abort()
			return
		}

		fmt.Printf("Token String: %s\n", tokenStr)

		userID, err := ValidateJWT(tokenStr)
		if err != nil {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token: " + err.Error()})
			ctx.Abort()
			return
		}
		fmt.Printf("UserID: %s\n", userID)

		// Сохраняем userID в контексте запроса
		ctx.Set("userID", userID)
		ctx.Next()
	}
}
*/
