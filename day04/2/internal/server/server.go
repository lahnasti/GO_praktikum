package server

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/lahnasti/GO_praktikum/internal/domain/models"
)

type Repository interface {
	AdduUser(models.User) (string, error)
	GetUserByID(id string) (models.User, error)
	GetUsers()([]models.User, error)
	UpdateUser(id string, user models.User) error
	DeleteUser(id string) error
}

type Server struct {
	db Repository
	valid *validator.Validate
}

func New(db Repository) *Server {
	valid := validator.New()
	return &Server {
		db: db,
		valid: valid,
	}
}

func (s *Server) GetUsersHandler(ctx *gin.Context) {
	users, err := s.db.GetUsers()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "List users", "users": users})
}

func (s *Server) RegisterUser(ctx *gin.Context) {
	var user models.User
	err := ctx.ShouldBindBodyWithJSON(&user)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Invalid params", "error": err.Error()})
		return
	}
	err = s.valid.Struct(user)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Data has not been validated", "error": err.Error()})
		return
	}
	userID, err := s.db.AdduUser(user)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to save user"})
		return
	}
	ctx.JSON(200, gin.H{"message": "User successfully registered", "user_id": userID})
}

func (s *Server) GetUserByIDHandler(ctx *gin.Context) {
	id := ctx.Param("id")
	user, err := s.db.GetUserByID(id)
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
	err := s.db.UpdateUser(id, user)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "User updated", "user": user})
}

func (s *Server) DeleteUserHandler(ctx *gin.Context) {
	id := ctx.Param("id")
	err := s.db.DeleteUser(id)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(200, gin.H{"message": "User deleted", "user_id": id})
}