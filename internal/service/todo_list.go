package service

import (
	"github.com/balantrea/todo-app/internal/model"
	"github.com/balantrea/todo-app/internal/repository"
	"github.com/rs/zerolog"
)

type TodoListService struct {
	repo   repository.TodoList
	logger zerolog.Logger
}

func NewTodoListService(repo repository.TodoList, logger zerolog.Logger) *TodoListService {
	return &TodoListService{repo: repo, logger: logger}
}

func (s *TodoListService) Create(userId int, list model.TodoList) (int, error) {
	return s.repo.Create(userId, list)
}

func (s *TodoListService) GetAll(userId int) ([]model.TodoList, error) {
	return s.repo.GetAll(userId)
}

func (s *TodoListService) GetById(userId, listId int) (model.TodoList, error) {
	return s.repo.GetById(userId, listId)
}

func (s *TodoListService) UpdateList(userId, listId int, input model.UpdateListInput) error {
	if err := input.Validate(); err != nil {
		return err
	}

	return s.repo.UpdateList(userId, listId, input)
}

func (s *TodoListService) DeleteList(userId, listId int) error {
	return s.repo.DeleteList(userId, listId)
}
