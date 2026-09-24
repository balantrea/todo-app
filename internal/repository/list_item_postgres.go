package repository

import (
	"fmt"
	"strings"

	"github.com/balantrea/todo-app/internal/model"
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

func (r *TodoItemRepository) Create(listId int, item model.TodoItem) (int, error) {
	tx, err := r.db.Begin()
	if err != nil {
		return 0, err
	}

	var itemId int

	createItemQuery := fmt.Sprintf(`
INSERT INTO %s (title, description) 
VALUES ($1, $2) RETURNING id`,
		todoItemsTable)

	row := tx.QueryRow(createItemQuery, item.Title, item.Description)

	if err = row.Scan(&itemId); err != nil {
		if rollbackErr := tx.Rollback(); rollbackErr != nil {
			r.logger.Err(rollbackErr).
				Msg("failed to create item")
		}
		return 0, err
	}

	createListItemQuery := fmt.Sprintf(`
INSERT INTO %s (list_id, item_id) 
VALUES ($1, $2)`,
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

func (r *TodoItemRepository) GetAll(listId, userId int) ([]model.TodoItem, error) {
	var lists []model.TodoItem

	query := fmt.Sprintf(`
SELECT ti.id, ti.title, ti.description, ti.done 
FROM %s ti 
INNER JOIN %s li on li.item_id = ti.id
INNER JOIN %s ul on ul.list_id = li.list_id 
WHERE li.list_id = $1 AND ul.user_id = $2`,
		todoItemsTable, listItemTable, userListsTable)

	if err := r.db.Select(&lists, query, listId, userId); err != nil {
		return nil, err
	}

	return lists, nil
}

func (r *TodoItemRepository) GetById(userId, itemId int) (model.TodoItem, error) {
	var list model.TodoItem

	query := fmt.Sprintf(`
SELECT
    ti.id,
    ti.title,
    ti.description,
    ti.done
FROM %s ti
INNER JOIN %s li ON li.item_id = ti.id
INNER JOIN %s ul ON ul.list_id = li.list_id
WHERE ti.id = $1 AND ul.user_id = $2`,
		todoItemsTable, listItemTable, userListsTable)

	if err := r.db.Get(&list, query, itemId, userId); err != nil {
		return list, err
	}

	return list, nil
}

func (r *TodoItemRepository) Update(userId, itemId int, input model.UpdateItemInput) error {
	setValues := make([]string, 0)
	args := make([]interface{}, 0)
	argsId := 1

	if input.Title != nil {
		setValues = append(setValues, fmt.Sprintf("title=$%d", argsId))
		args = append(args, *input.Title)
		argsId++
	}

	if input.Description != nil {
		setValues = append(setValues, fmt.Sprintf("description=$%d", argsId))
		args = append(args, *input.Description)
		argsId++
	}

	if input.Done != nil {
		setValues = append(setValues, fmt.Sprintf("done=$%d", argsId))
		args = append(args, *input.Done)
		argsId++
	}

	setQuery := strings.Join(setValues, ", ")

	query := fmt.Sprintf(`
UPDATE %s ti SET %s FROM %s li, %s ul
WHERE ti.id = li.item_id AND li.list_id = ul.list_id AND ul.user_id = $%d AND ti.id = $%d`,
		todoItemsTable, setQuery, listItemTable, userListsTable, argsId, argsId+1)

	args = append(args, userId, itemId)

	if _, err := r.db.Exec(query, args...); err != nil {
		return err
	}

	return nil
}

func (r *TodoItemRepository) Delete(userId, itemId int) error {
	query := fmt.Sprintf(`
DELETE FROM %s ti USING %s li, %s ul 
WHERE ti.id = li.item_id AND li.list_id = ul.list_id AND ul.user_id = $1 AND ti.id = $2`,
		todoItemsTable, listItemTable, userListsTable)

	_, err := r.db.Exec(query, userId, itemId)
	if err != nil {
		return err
	}

	return nil
}
