package main

import (
	"log"

	"github.com/balantrea/todo-app"
	handler "github.com/balantrea/todo-app/pkg/handler"
	"github.com/balantrea/todo-app/pkg/repository"
	service "github.com/balantrea/todo-app/pkg/service"
)

func main() {
	repos := repository.NewRepository()
	service := service.NewService(repos)
	handlers := handler.NewHandler(service)

	srv := new(todo.Server)
	if err := srv.Run("8000", handlers.InitRouters()); err != nil {
		log.Fatalf("error occurred while running http server: %s\n", err.Error())
	}
}
