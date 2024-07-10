package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/jackc/pgx/v5"

	"github.com/lahnasti/GO_praktikum/internal/config"
	"github.com/lahnasti/GO_praktikum/internal/logger"

	"github.com/lahnasti/GO_praktikum/internal/repository/booksrepo"
	"github.com/lahnasti/GO_praktikum/internal/repository/usersrepo"

	"github.com/lahnasti/GO_praktikum/internal/server/books"
	"github.com/lahnasti/GO_praktikum/internal/server/users"

	"github.com/lahnasti/GO_praktikum/cmd/routes"
)

func main() {

	fmt.Println("Server started")

	cfg := config.ReadConfig()
	fmt.Println(cfg)

	zlog := logger.SetupLogger(cfg.DebugFlag)
	zlog.Debug().Any("config", cfg).Msg("Check cfg value")

	err := usersrepo.Migrations(cfg.DBAddr, cfg.MPath, zlog)
	if err != nil {
		zlog.Fatal().Err(err).Msg("Init migrations failed")
	}

	err = booksrepo.Migrations(cfg.DBAddr, cfg.MPath, zlog)
	if err != nil {
		zlog.Fatal().Err(err).Msg("Init migrations failed")
	}

	conn, err := initDB(cfg.DBAddr)
	if err != nil {
		zlog.Fatal().Err(err).Msg("Connection DB failed")
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

	routes.BookRoutes(r, &booksServer)
	routes.UserRoutes(r, &usersServer)
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