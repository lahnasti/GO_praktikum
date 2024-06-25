package main

import (
	"github.com/gin-gonic/gin"
	"github.com/lahnasti/GO_praktikum/internal/repository"
	"github.com/lahnasti/GO_praktikum/internal/server"
	"github.com/rs/zerolog/log"
)

func main() {
	log.Info().Msg("Service started")
	repo := repository.New()
	log.Debug().Any("repo", repo).Msg("Check new repo")
	server := server.New(repo)
	log.Debug().Any("server", server).Msg("Check new server")

	r := gin.Default()
	r.POST("/users", server.RegisterUser)
	r.GET("/users", server.GetUsersHandler)
	r.GET("/users/:id", server.GetUserByIDHandler)
	r.PUT("/users/:id", server.UpdateUserHandler)
	r.DELETE("/users/:id", server.DeleteUserHandler)
	
	r.Run(":8080")
}