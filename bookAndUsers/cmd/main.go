package main

import (
	"context"
	"fmt"
	"time"
	"os"
	"os/signal"
	"syscall"
	"log"
	"net/http"


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
	// Канал для получения системных сигналов
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)


	r := gin.Default()
	routes.BookRoutes(r, &server)
	routes.UserRoutes(r, &server)
	zlog.Info().Msg("Server was started")

	httpServer := &http.Server{
		Addr:    cfg.Addr,
		Handler: r,
	}

	if err := r.Run(cfg.Addr); err != nil {
		panic(err)
	}

	// Запуск сервера в отдельной горутине
	go func() {
		if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Could not listen on %s: %v\n", cfg.Addr, err)
		}
	}()

	log.Printf("Server is ready to handle requests at %s", cfg.Addr)

	// Блокируемся, ожидая сигнала завершения
	<-stop
	log.Println("Server is shutting down...")

	// Создаем контекст с таймаутом для graceful shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Завершаем сервер
	if err := httpServer.Shutdown(ctx); err != nil {
		log.Fatalf("Server Shutdown Failed:%+v", err)
	}
	log.Println("Server exited properly")
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