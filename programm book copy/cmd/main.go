package main

import (
	"fmt"
	"context"


	

	"github.com/go-playground/validator/v10"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"

	//"github.com/lahnasti/GO_praktikum/internal/logger"
	"github.com/lahnasti/GO_praktikum/internal/config"
	"github.com/lahnasti/GO_praktikum/internal/repository"
	"github.com/lahnasti/GO_praktikum/internal/server"
)

func main() {

	//zlog := logger.SetupLogger(true)
	//zlog.Debug().Msg("Logger was invited")
	// Получение строки подключения из переменной окружения
	

	cfg := config.ReadConfig()
	fmt.Println(cfg)

	conn, err := initDB(cfg.DBAddr)
	if err != nil {
		panic(err)
	}

	storage := repository.NewDB(conn)

	validate := validator.New()  // Инициализация валидатора

	server := server.Server{
		Db: &storage,
		Valid: validate,
	}
	r := gin.Default()
	r.GET("/books", server.GetBooksHandler)
	r.POST("/books", server.CreateBookHandler)

	//zlog.Info().Msg("Server was started")

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

