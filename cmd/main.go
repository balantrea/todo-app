package main

import (
	"log"

	"github.com/balantrea/todo-app"
)

func main() {
	srv := new(todo.Server)
	if err := srv.Run("8000"); err != nil {
		log.Fatalf("error occurred while running http server: %s\n", err.Error())
	}

}
