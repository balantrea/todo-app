package repository

import (
	"fmt"
	"strings"

	"github.com/balantrea/todo-app/internal/model"
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

func (r *TodoListPostgres) Create(userId int, list model.TodoList) (int, error) {
	tx, err := r.db.Begin()
	if err != nil {
		return 0, err
	}

	var id int

	createListQuery := fmt.Sprintf(`
INSERT INTO %s (title, description) 
VALUES ($1, $2) RETURNING id`,
		todoListTable)

	row := tx.QueryRow(createListQuery, list.Title, list.Description)
	if err := row.Scan(&id); err != nil {
		if rollbackErr := tx.Rollback(); rollbackErr != nil {
			r.logger.Err(rollbackErr).
				Msg("failed to create list")
		}
		return 0, err
	}

	createUsersListQuery := fmt.Sprintf(`
INSERT INTO %s (user_id, list_id) 
VALUES ($1, $2)`,
		userListsTable)

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

func (r *TodoListPostgres) GetAll(userId int) ([]model.TodoList, error) {
	var lists []model.TodoList

	query := fmt.Sprintf(`
SELECT tl.id, tl.title, tl.description 
FROM %s tl 
INNER JOIN %s ul on tl.id = ul.list_id 
WHERE ul.user_id = $1`,
		todoListTable, userListsTable)

	if err := r.db.Select(&lists, query, userId); err != nil {
		return nil, err
	}

	return lists, nil
}

func (r *TodoListPostgres) GetById(userId, listId int) (model.TodoList, error) {
	var list model.TodoList

	query := fmt.Sprintf(`
SELECT tl.id, tl.title, tl.description
FROM %s tl
INNER JOIN %s ul ON tl.id = ul.list_id
WHERE ul.user_id = $1 AND ul.list_id = $2`,
		todoListTable, userListsTable)

	if err := r.db.Get(&list, query, userId, listId); err != nil {
		return list, err
	}

	return list, nil
}

func (r *TodoListPostgres) UpdateList(userId, listId int, input model.UpdateListInput) error {
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

	setQuery := strings.Join(setValues, ", ")

	query := fmt.Sprintf(`
UPDATE %s tl 
SET %s 
FROM %s ul 
WHERE tl.id = ul.list_id AND ul.list_id=$%d AND ul.user_id=$%d`,
		todoListTable, setQuery, userListsTable, argsId, argsId+1)

	args = append(args, userId, listId)

	if _, err := r.db.Exec(query, args...); err != nil {
		return err
	}

	return nil
}

func (r *TodoListPostgres) DeleteList(userId, listId int) error {
	query := fmt.Sprintf(`
DELETE 
FROM %s tl 
USING %s ul 
WHERE tl.id = ul.list_id AND ul.id=$1 AND ul.list_id=$2`,
		todoListTable, userListsTable)

	_, err := r.db.Exec(query, userId, listId)
	if err != nil {
		return err
	}

	return nil
}
