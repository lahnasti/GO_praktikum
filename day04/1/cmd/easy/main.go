package main

import (
	"github.com/gin-gonic/gin"
	"github.com/lahnasti/GO_praktikum/day04/1/internal/repository"
	"github.com/lahnasti/GO_praktikum/day04/1/internal/server"
	"github.com/rs/zerolog/log"
)

func main() {
	log.Info().Msg("SERVICE STARTED")
	repo := repository.New()
	log.Debug().Any("repo", repo).Msg("Check new repo")
	server := server.New(repo)
	log.Debug().Any("server", server).Msg("Check new server")

	r := gin.Default()

	taskRoutes := r.Group("/tasks")
	{
		taskRoutes.GET("/", getTasks)
		taskRoutes.POST("/", createTasks)

		taskRoutes.GET("/:id", getTasksId)
		taskRoutes.PUT("/:id", updateTask)
		taskRoutes.DELETE("/:id", deleteTask)
	}

}