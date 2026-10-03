package repository

import (
	"os"
	"testing"

	"github.com/balantrea/todo-app/internal/model"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"github.com/zhashkevych/go-sqlxmock"
)

var logger = zerolog.New(os.Stdout).
	With().
	Timestamp().
	Logger()

func TestTodoItemPostgres_Create(t *testing.T) {
	db, mock, err := sqlmock.Newx()
	if err != nil {
		logger.Err(err)
		return
	}

	defer func() {
		if err = db.Close(); err != nil {
			logger.Err(err)
		}
	}()

	r := NewTodoItemRepository(db, logger)

	type args struct {
		listId int
		item   model.TodoItem
	}

	type mockBehavior func(args args, id int)

	testTable := []struct {
		name         string
		mockBehavior mockBehavior
		args         args
		id           int
		wantErr      bool
	}{
		{
			name: "OK",
			args: args{
				listId: 1,
				item: model.TodoItem{
					Title:       "Test",
					Description: "Test",
				},
			},
			id: 1,
			mockBehavior: func(args args, id int) {
				mock.ExpectBegin()

				rows := sqlmock.NewRows([]string{"id"}).AddRow(id)
				mock.ExpectQuery("INSERT INTO todo_item").
					WithArgs(args.item.Title, args.item.Description).WillReturnRows(rows)

				mock.ExpectExec("INSERT INTO list_item").WithArgs(args.listId, id).
					WillReturnResult(sqlmock.NewResult(1, 1))

				mock.ExpectCommit()
			},
		},
		{
			name: "invalid input",
			args: args{
				listId: 1,
				item:   model.TodoItem{},
			},
			wantErr:      true,
			mockBehavior: func(args args, id int) {},
		},
		{
			name: "empty field",
			args: args{
				listId: 1,
				item: model.TodoItem{
					Title:       "",
					Description: "Test",
				},
			},
			wantErr: true,
			id:      1,
			mockBehavior: func(args args, id int) {
				mock.ExpectBegin()

				rows := sqlmock.NewRows([]string{"id"}).AddRow(id).RowError(1, errTitleCannotBeEmpty)
				mock.ExpectQuery("INSERT INTO todo_item").
					WithArgs(args.item.Title, args.item.Description).WillReturnRows(rows)

				mock.ExpectRollback()
			},
		},
		{
			name: "2nd insert error",
			args: args{
				listId: 1,
				item: model.TodoItem{
					Title:       "Test",
					Description: "Test",
				},
			},
			wantErr: true,
			id:      1,
			mockBehavior: func(args args, id int) {
				mock.ExpectBegin()

				rows := sqlmock.NewRows([]string{"id"}).AddRow(id)
				mock.ExpectQuery("INSERT INTO todo_item").
					WithArgs(args.item.Title, args.item.Description).WillReturnRows(rows)

				mock.ExpectExec("INSERT INTO list_item").WithArgs(args.listId, id).
					WillReturnError(errFailedToInsertListItem)

				mock.ExpectRollback()
			},
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			testCase.mockBehavior(testCase.args, testCase.id)

			got, err := r.Create(testCase.args.listId, testCase.args.item)
			if testCase.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, testCase.id, got)
			}
		})
	}
}
