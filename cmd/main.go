package main

import (
	"log"

	"github.com/balantrea/todo-app"
	"github.com/balantrea/todo-app/pkg/handler"
	"github.com/balantrea/todo-app/pkg/repository"
	"github.com/balantrea/todo-app/pkg/service"
	"github.com/spf13/viper"
)

func main() {
	if err := initConfig(); err != nil {
		log.Fatalf("error initializing configs: %s", err.Error())
	}

	repos := repository.NewRepository()
	services := service.NewService(repos)
	handlers := handler.NewHandler(services)

	srv := new(todo.Server)
	if err := srv.Run(viper.GetString("8000"), handlers.InitRouters()); err != nil {
		log.Fatalf("error occurred while running http server: %s\n", err.Error())
	}
}

func initConfig() error {
	viper.AddConfigPath("configs")
	viper.SetConfigName("config")
	return viper.ReadInConfig()
}
