package main

import (
	"context"
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
	"github.com/go-playground/validator/v10"

	"github.com/lahnasti/GO_praktikum/internal/config"
	"github.com/lahnasti/GO_praktikum/internal/repository"
	"github.com/lahnasti/GO_praktikum/internal/server"
	"github.com/rs/zerolog/log"
)

func main() {
	log.Info().Msg("Service started")

	cfg := config.ReadConfig()
	fmt.Println(cfg)

	conn, err := initDB(cfg.DBAddr)
	if err != nil {
		panic(err)
	}
	//repo := repository.New()
	storage := repository.NewDB(conn)

	log.Debug().Any("storage", storage).Msg("Check new storage")

	validate := validator.New()  // Инициализация валидатора

	server := server.Server{
		Db: &storage,
		Valid: validate,
	}
	log.Debug().Any("server", server).Msg("Check new server")

	r := gin.Default()

	r.POST("/users", server.RegisterUser)
	r.GET("/users", server.GetUsersHandler)
	r.GET("/users/:id", server.GetUserByIDHandler)
	r.PUT("/users/:id", server.UpdateUserHandler)
	r.DELETE("/users/:id", server.DeleteUserHandler)

	if err := r.Run(cfg.Addr); err != nil {
		panic(err)
	}
}

func initDB(addr string) (*pgx.Conn, error) {
	conn, err := pgx.Connect(context.Background(), addr)
	if err != nil {
		return nil,
			fmt.Errorf("database initialization error: %w", err)
	}
	return conn, nil
}
