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
	Create(userId int, list todo.TodoList) (int, error)
	GetAll(userId int) ([]todo.TodoList, error)
	GetById(userId, listId int) (todo.TodoList, error)
	UpdateList(userId, listId int, input todo.UpdateListInput) error
	DeleteList(userId, listId int) error
}

type TodoItem interface {
	Create(listId int, item todo.TodoItem) (int, error)
	GetAll(listId, userId int) ([]todo.TodoItem, error)
	GetById(userId, itemId int) (todo.TodoItem, error)
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
