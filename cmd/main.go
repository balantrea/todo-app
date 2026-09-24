package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/balantrea/todo-app/internal/config"
	"github.com/balantrea/todo-app/internal/handler"
	"github.com/balantrea/todo-app/internal/repository"
	"github.com/balantrea/todo-app/internal/server"
	"github.com/balantrea/todo-app/internal/service"
	_ "github.com/lib/pq"
	"github.com/rs/zerolog"
	"github.com/spf13/viper"
)

func main() {
	logger := zerolog.New(os.Stdout).
		With().
		Timestamp().
		Logger()

	if err := run(logger); err != nil {
		logger.Err(err).
			Msg("todo app is terminated")
		return
	}
}

func run(logger zerolog.Logger) error {
	if err := initConfig(); err != nil {
		return fmt.Errorf("initialize config: %w", err)
	}

	db, err := repository.NewPostgresDB(config.DB{
		Host:     os.Getenv("DB_HOST"),
		Port:     os.Getenv("DB_PORT"),
		Username: os.Getenv("DB_USER"),
		Password: os.Getenv("DB_PASSWORD"),
		DBName:   os.Getenv("DB_NAME"),
		SSLMode:  os.Getenv("SSLMODE"),
	},
		logger,
	)

	if err != nil {
		return fmt.Errorf("initialize database: %w", err)
	}

	repos := repository.NewRepository(db, logger)
	services := service.NewService(repos, logger)
	handlers := handler.NewHandler(services, logger)

	srv := server.NewServer(logger)
	go func() {
		if err = srv.Run(viper.GetString("port"), handlers.InitRouters()); err != nil {
			logger.Err(err).
				Msg("error occurred while running HTTP server")
		}
	}()

	logger.Info().
		Msgf("TodoApp Started in port: %s", viper.GetString("port"))

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)
	<-quit

	logger.Info().
		Msg("TodoApp Shutting Down")

	ctx, cancel := context.WithTimeout(
		context.Background(),
		10*time.Second,
	)
	defer cancel()

	if err = srv.Shutdown(ctx); err != nil {
		logger.Error().
			Err(err).
			Msg("error occurred on server shutting down")
	}

	if err = db.Close(); err != nil {
		logger.Err(err).
			Msg("error occurred on db connection close")
	}

	return nil
}

func initConfig() error {
	viper.AddConfigPath("configs")
	viper.SetConfigName("config")
	return viper.ReadInConfig()
}
