package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
	"golang.org/x/sync/errgroup"

	"github.com/lahnasti/GO_praktikum/internal/config"
	"github.com/lahnasti/GO_praktikum/internal/logger"
	"github.com/lahnasti/GO_praktikum/internal/repository"

	"github.com/lahnasti/GO_praktikum/internal/server"

	"github.com/lahnasti/GO_praktikum/internal/server/routes"
)

func main() {

	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT)

		<-c
		cancel()
	}()
	fmt.Println("Server starting")
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
	group, gCtx := errgroup.WithContext(ctx)
	srv := server.NewServer(gCtx, &storage, zlog)
	//validate := validator.New() // Инициализация валидатора
	group.Go(func() error {
		r := gin.Default()
		routes.BookRoutes(r, srv)
		routes.UserRoutes(r, srv)
		zlog.Info().Msg("Server was started")

		if err := r.Run(cfg.Addr); err != nil {
			return err
		}
		return nil
	})

	group.Go(func() error {
		err := <-srv.ErrorChan
		return err
	})
	group.Go(func() error {
		<-gCtx.Done()
		return gCtx.Err()
	})

	if err := group.Wait(); err != nil {
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
