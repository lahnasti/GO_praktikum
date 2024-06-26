package main

import (
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	"github.com/lahnasti/GO_praktikum/internal/server"
	"github.com/lahnasti/GO_praktikum/internal/repository"
)

func main() {
	log.Info().Msg("Service started")
	repo := repository.New()
	log.Debug().Any("repo", repo).Msg("Check new repo")
	server := server.New(repo)
	log.Debug().Any("server", server).Msg("Check new server")

	r := gin.Default()
	r.GET("/tasks", server.GetTasksHandler) 
	r.POST("/tasks", server.AddTaskHandler)
	r.GET("/tasks/:id", server.GetTaskByIDHandler)
	r.PUT("/tasks/:id", server.UpdateTaskHandler)
	r.DELETE("/tasks/:id", server.DeleteTaskHandler)

	r.Run(":8080")
}
