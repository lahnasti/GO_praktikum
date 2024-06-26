package main

import (
	"github.com/gin-gonic/gin"
	"github.com/lahnasti/GO_praktikum/internal/logger"
	"github.com/lahnasti/GO_praktikum/internal/repository"
	"github.com/lahnasti/GO_praktikum/internal/server"
)

func main() {
	zlog := logger.SetupLogger(true)
	zlog.Debug().Msg("Logger was inited")
	repo := repository.New(zlog)
	server := server.New(repo, zlog)

	r := gin.Default()
	r.GET("/tasks", server.GetTasksHandler) 
	r.POST("/tasks", server.AddTaskHandler)
	r.GET("/tasks/:id", server.GetTaskByIDHandler)
	r.PUT("/tasks/:id", server.UpdateTaskHandler)
	r.DELETE("/tasks/:id", server.DeleteTaskHandler)

	zlog.Info().Msg("Server was started")
	r.Run(":8080")
}
