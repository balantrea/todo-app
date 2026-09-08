package repository

import (
	"github.com/balantrea/todo-app"
	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog"
)

type Authorization interface {
	CreateUser(user todo.User) (int, error)
	GetUser(username, password string) (todo.User, error)
}

type TodoList interface {
}

type TodoItem interface {
}

type Repository struct {
	Authorization
	TodoList
	TodoItem
	logger zerolog.Logger
}

func NewRepository(db *sqlx.DB, logger zerolog.Logger) *Repository {
	return &Repository{
		Authorization: NewAuthPostgres(db),
		logger:        logger,
	}
}
