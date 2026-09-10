package main

import (
	"fmt"
	"os"

	"github.com/balantrea/todo-app"
	"github.com/balantrea/todo-app/internal/pkg/handler"
	"github.com/balantrea/todo-app/internal/pkg/repository"
	"github.com/balantrea/todo-app/internal/pkg/service"
	"github.com/joho/godotenv"
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

	if err := godotenv.Load(); err != nil {
		return fmt.Errorf("load environment variables: %w", err)
	}

	db, err := repository.NewPostgresDB(repository.Config{
		Host:     viper.GetString("db.host"),
		Port:     viper.GetString("db.port"),
		Username: viper.GetString("db.username"),
		Password: os.Getenv("DB_PASSWORD"),
		DBName:   viper.GetString("db.dbname"),
		SSLMode:  viper.GetString("db.sslmode"),
	},
		logger,
	)

	if err != nil {
		return fmt.Errorf("initialize database: %w", err)
	}

	repos := repository.NewRepository(db, logger)
	services := service.NewService(repos, logger)
	handlers := handler.NewHandler(services, logger)

	srv := todo.NewServer(logger)
	if err := srv.Run(viper.GetString("port"), handlers.InitRouters()); err != nil {
		return fmt.Errorf("run HTTP server: %w", err)
	}

	return nil
}

func initConfig() error {
	viper.AddConfigPath("configs")
	viper.SetConfigName("config")
	return viper.ReadInConfig()
}
