package service

import (
	"github.com/balantrea/todo-app"
	"github.com/balantrea/todo-app/internal/pkg/repository"
	"github.com/rs/zerolog"
)

type TodoListService struct {
	repo   repository.TodoList
	logger zerolog.Logger
}

func NewTodoListService(repo repository.TodoList, logger zerolog.Logger) *TodoListService {
	return &TodoListService{repo: repo, logger: logger}
}

func (s *TodoListService) Create(userId int, list todo.TodoList) (int, error) {
	return s.repo.Create(userId, list)
}
