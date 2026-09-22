package repository

import (
	"fmt"

	"github.com/balantrea/todo-app"
	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog"
)

type TodoItemRepository struct {
	db     *sqlx.DB
	logger zerolog.Logger
}

func NewTodoItemRepository(db *sqlx.DB, logger zerolog.Logger) *TodoItemRepository {
	return &TodoItemRepository{db: db, logger: logger}
}

func (r *TodoItemRepository) Create(listId int, item todo.TodoItem) (int, error) {
	tx, err := r.db.Begin()
	if err != nil {
		return 0, err
	}

	var itemId int

	createItemQuery := fmt.Sprintf("INSERT INTO %s (title, description) values ($1, $2) RETURNING id",
		todoItemsTable)

	row := tx.QueryRow(createItemQuery, item.Title, item.Description)

	if err := row.Scan(&itemId); err != nil {
		if rollbackErr := tx.Rollback(); rollbackErr != nil {
			r.logger.Err(rollbackErr).
				Msg("failed to create item")
		}
		return 0, err
	}

	createListItemQuery := fmt.Sprintf("INSERT INTO %s (list_id, item_id) values ($1, $2)",
		listItemTable)

	_, err = tx.Exec(createListItemQuery, listId, itemId)
	if err != nil {
		if rollbackErr := tx.Rollback(); rollbackErr != nil {
			r.logger.Err(rollbackErr).
				Msg("failed to create item list")
		}
		return 0, err
	}

	if err = tx.Commit(); err != nil {
		return 0, err
	}

	return itemId, nil
}
