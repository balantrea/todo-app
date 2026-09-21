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

func (s *TodoListService) GetAll(userId int) ([]todo.TodoList, error) {
	return s.repo.GetAll(userId)
}

func (s *TodoListService) GetById(userId, listId int) (todo.TodoList, error) {
	return s.repo.GetById(userId, listId)
}

func (s *TodoListService) UpdateList(userId, listId int) (int, error) {
	return s.repo.UpdateList(userId, listId)
}

func (s *TodoListService) DeleteList(userId, listId int) error {
	return s.repo.DeleteList(userId, listId)
}
