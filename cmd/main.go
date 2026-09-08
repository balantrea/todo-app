package main

import (
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

	if err := initConfig(); err != nil {
		logger.Err(err).
			Msg("error initializing configs")
		return
	}

	if err := godotenv.Load(); err != nil {
		logger.Err(err).
			Msg("loading env variables")
		return
	}

	db, err := repository.NewPostgresDB(repository.Config{
		Host:     viper.GetString("db.host"),
		Port:     viper.GetString("db.port"),
		Username: viper.GetString("db.username"),
		Password: os.Getenv("DB_PASSWORD"),
		DBName:   viper.GetString("db.dbname"),
		SSLMode:  viper.GetString("db.sslmode"),
	})

	if err != nil {
		logger.Err(err).
			Msg("failed to initialize db")
		return
	}

	repos := repository.NewRepository(db, logger)
	services := service.NewService(repos, logger)
	handlers := handler.NewHandler(services, logger)

	srv := todo.NewServer(logger)
	if err := srv.Run(viper.GetString("port"), handlers.InitRouters()); err != nil {
		logger.Err(err).
			Msg("error occurred while running http server")
		return
	}
}

func initConfig() error {
	viper.AddConfigPath("configs")
	viper.SetConfigName("config")
	return viper.ReadInConfig()
}
