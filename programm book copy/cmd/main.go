package main

import (
	"context"
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/jackc/pgx/v5"

	"github.com/lahnasti/GO_praktikum/internal/server/users/jwt"
	"github.com/lahnasti/GO_praktikum/internal/config"
	
	"github.com/lahnasti/GO_praktikum/internal/repository/booksrepo"
	"github.com/lahnasti/GO_praktikum/internal/repository/usersrepo"

	"github.com/lahnasti/GO_praktikum/internal/server/books"
	"github.com/lahnasti/GO_praktikum/internal/server/users"

	"github.com/lahnasti/GO_praktikum/cmd/routes"
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

	storageUser := usersrepo.NewDB(conn)

	if err := storageUser.CreateTable(context.Background()); err != nil {
		log.Fatalf("Failed to create table: %v", err)
	}

	storageBook := booksrepo.NewDB(conn)

	if err := storageBook.CreateTable(context.Background()); err != nil {
		log.Fatalf("Failed to create table: %v", err)
	}

	validate := validator.New() // Инициализация валидатора

	usersServer := users.Server{
		Db: &storageUser,
		Valid: validate,
	}

	booksServer := books.Server{
		Db:    &storageBook,
		Valid: validate,
	}


	r := gin.Default()

	r.POST("/login", jwt.LoginHandler)

	routes.BookRoutes(r, &booksServer)
	routes.UserRoutes(r, &usersServer)
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
