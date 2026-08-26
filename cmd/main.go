package main

import (
	"log"

	"github.com/balantrea/todo-app"
	handler "github.com/balantrea/todo-app/pkg/handler"
)

func main() {
	handlers := new(handler.Handler)
	srv := new(todo.Server)
	if err := srv.Run("8000", handlers.InitRouters()); err != nil {
		log.Fatalf("error occurred while running http server: %s\n", err.Error())
	}
}
