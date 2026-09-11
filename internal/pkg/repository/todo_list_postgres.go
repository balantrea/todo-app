package repository

import (
	"fmt"

	"github.com/balantrea/todo-app"
	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog"
)

type TodoListPostgres struct {
	db     *sqlx.DB
	logger zerolog.Logger
}

func NewTodoListPostgres(db *sqlx.DB, logger zerolog.Logger) *TodoListPostgres {
	return &TodoListPostgres{db: db, logger: logger}
}

func (r *TodoListPostgres) Create(userId int, list todo.TodoList) (int, error) {
	tx, err := r.db.Begin()
	if err != nil {
		return 0, err
	}

	var id int

	createListQuery := fmt.Sprintf("INSERT INTO %s (title, description) VALUES ($1, $2) RETURNING id", todoListTable)
	row := tx.QueryRow(createListQuery, list.Title, list.Description)
	if err := row.Scan(&id); err != nil {
		if rollbackErr := tx.Rollback(); rollbackErr != nil {
			r.logger.Err(rollbackErr).
				Msg("failed to create list")
		}
		return 0, err
	}

	createUsersListQuery := fmt.Sprintf("INSERT INTO %s (user_id, list_id) VALUES ($1, $2)", userListsTable)
	_, err = tx.Exec(createUsersListQuery, userId, id)
	if err != nil {
		if rollbackErr := tx.Rollback(); rollbackErr != nil {
			r.logger.Err(rollbackErr).
				Msg("failed to create user list")
		}

		return 0, err
	}

	return id, tx.Commit()
}
