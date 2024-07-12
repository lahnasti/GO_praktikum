package main

import (
	"context"
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/jackc/pgx/v5"

	"github.com/lahnasti/GO_praktikum/internal/config"
	"github.com/lahnasti/GO_praktikum/internal/logger"
	"github.com/lahnasti/GO_praktikum/internal/repository"

	"github.com/lahnasti/GO_praktikum/internal/server"

	"github.com/lahnasti/GO_praktikum/internal/server/routes"
)

func main() {

	fmt.Println("Server started")

	cfg := config.ReadConfig()
	fmt.Println(cfg)

	zlog := logger.SetupLogger(cfg.DebugFlag)
	zlog.Debug().Any("config", cfg).Msg("Check cfg value")


	err := repository.Migrations(cfg.DBAddr, cfg.MPath, zlog)
	if err != nil {
		zlog.Fatal().Err(err).Msg("Init migrations failed")
	}

	conn, err := initDB(cfg.DBAddr)
	if err != nil {
		zlog.Fatal().Err(err).Msg("Connection DB failed")
	}

	storage, err := repository.NewDB(conn)

	if err != nil {
		panic(err)
	}


	validate := validator.New() // Инициализация валидатора

	server := server.Server{
		BooksDB: &storage,
		UsersDB: &storage,
		Valid: validate,
	}


	r := gin.Default()
	routes.BookRoutes(r, &server)
	routes.UserRoutes(r, &server)
	//zlog.Info().Msg("Server was started")

	if err := r.Run(cfg.Addr); err != nil {
		panic(err)
	}

}

func initDB(addr string) (*pgx.Conn, error) {
	for i := 0; i < 7; i++ {
		time.Sleep(2 * time.Second)
		conn, err := pgx.Connect(context.Background(), addr)
		if err == nil {
			return conn, nil
		}
	}
	conn, err := pgx.Connect(context.Background(), addr)
	if err != nil {
		return nil, fmt.Errorf("database initialization error: %w", err)
	}
	return conn, nil
}