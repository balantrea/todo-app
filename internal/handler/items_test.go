package handler

import (
	"bytes"
	"net/http/httptest"
	"testing"

	"github.com/balantrea/todo-app/internal/model"
	"github.com/balantrea/todo-app/internal/service"
	"github.com/balantrea/todo-app/internal/service/mocks"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/assert/v2"
	"go.uber.org/mock/gomock"
)

func TestItems_Create(t *testing.T) {
	type mockBehavior func(s *mock_service.MockTodoItem, todoItem model.TodoItem)

	testTable := []struct {
		name                 string
		inputURL             string
		inputBody            string
		inputItem            model.TodoItem
		mockBehavior         mockBehavior
		expectedResponseCode int
		expectedResponseBody string
	}{
		{
			name:      "OK",
			inputURL:  "/api/list/1/items/",
			inputBody: `{"title":"Test","description":"Test"}`,
			inputItem: model.TodoItem{Title: "Test", Description: "Test"},
			mockBehavior: func(s *mock_service.MockTodoItem, todoItem model.TodoItem) {
				s.EXPECT().Create(1, 1, todoItem).Return(1, nil)
			},
			expectedResponseCode: 200,
			expectedResponseBody: `{"id":1}`,
		},
		{
			name:                 "invalid input",
			inputURL:             "/api/list/1/items/",
			inputBody:            `{}`,
			mockBehavior:         func(s *mock_service.MockTodoItem, todoItem model.TodoItem) {},
			expectedResponseCode: 400,
			expectedResponseBody: `{"message":"invalid input body"}`,
		},
		{
			name:      "failed to create item",
			inputURL:  "/api/list/1/items/",
			inputBody: `{"title":"Test","description":"Test"}`,
			inputItem: model.TodoItem{Title: "Test", Description: "Test"},
			mockBehavior: func(s *mock_service.MockTodoItem, todoItem model.TodoItem) {
				s.EXPECT().Create(1, 1, todoItem).Return(0, errFailedToCreateItem)
			},
			expectedResponseCode: 500,
			expectedResponseBody: `{"message":"failed to create item"}`,
		},
		{
			name:      "invalid id param",
			inputURL:  "/api/list/invalid/items/",
			inputBody: `{"title":"Test","description":"Test"}`,
			inputItem: model.TodoItem{Title: "Test", Description: "Test"},
			mockBehavior: func(s *mock_service.MockTodoItem, todoItem model.TodoItem) {
			},
			expectedResponseCode: 400,
			expectedResponseBody: `{"message":"invalid id param"}`,
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			todoItems := mock_service.NewMockTodoItem(c)
			testCase.mockBehavior(todoItems, testCase.inputItem)

			services := &service.Service{TodoItem: todoItems}
			handler := NewHandler(services, logger)

			r := gin.New()
			r.POST("/api/list/:id/items/", func(c *gin.Context) {
				c.Set("userId", 1)
				handler.createItems(c)
			})

			w := httptest.NewRecorder()
			req := httptest.NewRequest("POST", testCase.inputURL, bytes.NewBufferString(testCase.inputBody))

			r.ServeHTTP(w, req)

			assert.Equal(t, testCase.expectedResponseCode, w.Code)
			assert.Equal(t, testCase.expectedResponseBody, w.Body.String())
		})
	}
}
