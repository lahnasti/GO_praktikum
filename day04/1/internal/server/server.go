package server

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/lahnasti/GO_praktikum/internal/domain/models"
)

type Repository interface {
	GetAllTasks()([]models.Task, error)
	AddTask(models.Task) (string, error)
	GetTaskByID(id string) (models.Task, error)
	UpdateTask(id string, task models.Task) error
	DeleteTask(id string) error
}

type Server struct {
	db Repository
	valid *validator.Validate
}

func New(db Repository) *Server {
	valid := validator.New()
	return &Server{
		db: db,
		valid: valid,
	}
}

func (s *Server) GetTasksHandler(ctx *gin.Context) {
	tasks, err := s.db.GetAllTasks()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError,  gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "List tasks", "users": tasks})

}

func (s *Server) AddTaskHandler(ctx *gin.Context) {
	var task models.Task
	err := ctx.ShouldBindBodyWithJSON(&task)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Invalid params", "error": err.Error()})
		return
	}
	err = s.valid.Struct(task)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Data has not been validated", "error": err.Error()})
		return
	}
	taskID, err := s.db.AddTask(task)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to save task"})
		return
	}
	ctx.JSON(200, gin.H{"message": "Task successfully added", "task_id": taskID})
}

func (s *Server) GetTaskByIDHandler (ctx *gin.Context) {
	id := ctx.Param("id")
	task, err := s.db.GetTaskByID(id)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(200, gin.H{"message": "Task retrieved", "task": task})
}

func (s *Server) UpdateTaskHandler (ctx *gin.Context) {
	var task models.Task
	if err := ctx.ShouldBindJSON(&task); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	id := ctx.Param("id")
	task.ID = id
	err := s.db.UpdateTask(id, task)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Task updated", "task": task})

}

func (s *Server) DeleteTaskHandler(ctx *gin.Context) {
	id := ctx.Param("id")
	err := s.db.DeleteTask(id)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(200, gin.H{"message": "Task deleted", "task_id": id})
}