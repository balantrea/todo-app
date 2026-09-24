package service

import (
	"github.com/balantrea/todo-app/internal/model"
	"github.com/balantrea/todo-app/internal/repository"
	"github.com/rs/zerolog"
)

type Authorization interface {
	CreateUser(user model.User) (int, error)
	GenerateToken(username, password string) (string, error)
	ParseToken(token string) (int, error)
}

type TodoList interface {
	Create(userId int, list model.TodoList) (int, error)
	GetAll(userId int) ([]model.TodoList, error)
	GetById(userId, listId int) (model.TodoList, error)
	UpdateList(userId, listId int, input model.UpdateListInput) error
	DeleteList(userId, listId int) error
}

type TodoItem interface {
	Create(userId, listId int, item model.TodoItem) (int, error)
	GetAll(userId, listId int) ([]model.TodoItem, error)
	GetById(userId, itemId int) (model.TodoItem, error)
	Update(userId, itemId int, input model.UpdateItemInput) error
	Delete(userId, itemId int) error
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
