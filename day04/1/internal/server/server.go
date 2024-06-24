package server

import (
	"net/http"
	"github.com/google/uuid"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/lahnasti/GO_praktikum/day04/1/internal/domain/models"
)

type Repository interface {
	GetTasks() ([]models.Task)
	AddTasks(models.Task) error
	GetTasksById(id string) (models.Task, error)
	UpdateTask(id string) (models.Task, error)
	DeleteTask(id string) (models.Task, error)
}

type Server struct {
	list Repository
	validate *validator.Validate
}

func New(list Repository) *Server {
	validate := validator.New()
	return &Server{
		list: list,
		validate: validate,
	}
}

func notBlank(fl validator.FieldLevel) bool {
	return strings.TrimSpace(fl.Field().String()) != ""
}

// для проверки статуса
func validateStatus(fl validator.FieldLevel) bool {
	status := fl.Field().String()
	for _, s := range models.AllowedStatus {
		if status == s {
			return true
		}
	}
	return false
}

func (s *Server) GetTasks(ctx *gin.Context) {
	ctx.JSON(200, models.Task)
}

func (s *Server) CreateTasks(ctx *gin.Context) {
	var task models.Task

	if err := ctx.ShouldBindBodyWithJSON(&task); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if task.Status == "" {
		task.Status = models.DefaultStatus
	}

	if err := s.validate.Struct(task); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Генерация уникального идентификатора для задачи
	task.ID = uuid.New()

	//idString := task.ID.String()

	s.list.AddTasks(task)
}

/*

func getTasksId(c *gin.Context) {
	id := c.Param("id")
	task := taskMap.List[id]
	c.JSON(http.StatusOK, gin.H{"message": "Task retrieved", "task": task})

}


func updateTask(c *gin.Context) {
	var task Task
	if err := c.ShouldBindBodyWithJSON(&task); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	id := c.Param("id")
	taskMap.List[id] = task
	c.JSON(200, gin.H{"message": "Task updated", "id": id})
}

func deleteTask(c *gin.Context) {
	id := c.Param("id")
	delete(taskMap.List, id)
	c.JSON(200, gin.H{"message": "Task deleted", "id": id})
}
*/