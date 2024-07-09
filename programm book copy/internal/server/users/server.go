package users

import (
	"net/http"
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/lahnasti/GO_praktikum/internal/models"
	"github.com/lahnasti/GO_praktikum/internal/server/users/jwt"
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
}

type Server struct {
	Db    Repository
	Valid *validator.Validate
}

func (s *Server) GetUsersHandler(ctx *gin.Context) {
	users, err := s.Db.GetUsers()
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

	userID, err := s.Db.AddUser(user)
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
	user, err := s.Db.GetUserByID(id)
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
	err := s.Db.UpdateUser(id, user)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "User updated", "user": user})
}

func (s *Server) DeleteUserHandler(ctx *gin.Context) {
	id := ctx.Param("id")
	err := s.Db.DeleteUser(id)
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

	userID, err := s.Db.AddMultipleUsers(users)
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
	dbUser, err := s.Db.FindUserByEmail(user.Email)
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
	tokenString, err := jwt.GenerateJWT(dbUser)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"token": tokenString})
}
