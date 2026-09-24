package repository

import (
	"github.com/balantrea/todo-app/internal/model"
	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog"
)

type Authorization interface {
	CreateUser(user model.User) (int, error)
	GetUser(username, password string) (model.User, error)
}

type TodoList interface {
	Create(userId int, list model.TodoList) (int, error)
	GetAll(userId int) ([]model.TodoList, error)
	GetById(userId, listId int) (model.TodoList, error)
	UpdateList(userId, listId int, input model.UpdateListInput) error
	DeleteList(userId, listId int) error
}

type TodoItem interface {
	Create(listId int, item model.TodoItem) (int, error)
	GetAll(listId, userId int) ([]model.TodoItem, error)
	GetById(userId, itemId int) (model.TodoItem, error)
	Update(userId, itemId int, input model.UpdateItemInput) error
	Delete(userId, itemId int) error
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
		TodoList:      NewTodoListPostgres(db, logger),
		TodoItem:      NewTodoItemRepository(db, logger),
		logger:        logger,
	}
}
