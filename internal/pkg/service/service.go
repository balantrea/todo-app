package service

import (
	"github.com/balantrea/todo-app"
	"github.com/balantrea/todo-app/internal/pkg/repository"
	"github.com/rs/zerolog"
)

type Authorization interface {
	CreateUser(user todo.User) (int, error)
	GenerateToken(username, password string) (string, error)
	ParseToken(token string) (int, error)
}

type TodoList interface {
}

type TodoItem interface {
}

type Service struct {
	Authorization
	TodoList
	TodoItem
	logger zerolog.Logger
}

func NewService(repos *repository.Repository, logger zerolog.Logger) *Service {
	return &Service{
		logger:        logger,
		Authorization: NewAuthService(repos.Authorization),
	}
}
