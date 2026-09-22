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
	Create(userId int, list todo.TodoList) (int, error)
	GetAll(userId int) ([]todo.TodoList, error)
	GetById(userId, listId int) (todo.TodoList, error)
	UpdateList(userId, listId int, input todo.UpdateListInput) error
	DeleteList(userId, listId int) error
}

type TodoItem interface {
	Create(userId, listId int, item todo.TodoItem) (int, error)
	GetAll(userId, listId int) ([]todo.TodoItem, error)
}

type Service struct {
	Authorization
	TodoList
	TodoItem
	logger zerolog.Logger
}

func NewService(repos *repository.Repository, logger zerolog.Logger) *Service {
	return &Service{
		Authorization: NewAuthService(repos.Authorization),
		TodoList:      NewTodoListService(repos.TodoList, logger),
		TodoItem:      NewTodoItemService(repos.TodoItem, repos.TodoList),
		logger:        logger,
	}
}
