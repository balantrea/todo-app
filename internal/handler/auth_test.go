package handler

import (
	"bytes"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/balantrea/todo-app/internal/model"
	"github.com/balantrea/todo-app/internal/service"
	"github.com/balantrea/todo-app/internal/service/mocks"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/assert/v2"
	"github.com/rs/zerolog"
	"go.uber.org/mock/gomock"
)

var logger = zerolog.New(os.Stdout).
	With().
	Timestamp().
	Logger()

func TestHandler_signUp(t *testing.T) {
	type mockBehavior func(s *mock_service.MockAuthorization, user model.User)

	testTable := []struct {
		name                 string
		inputBody            string
		inputUser            model.User
		mockBehavior         mockBehavior
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name:      "OK",
			inputBody: `{"name":"Test","username":"Test","password":"Test"}`,
			inputUser: model.User{Name: "Test", Username: "Test", Password: "Test"},
			mockBehavior: func(s *mock_service.MockAuthorization, user model.User) {
				s.EXPECT().CreateUser(user).Return(1, nil)
			},
			expectedStatusCode:   200,
			expectedResponseBody: `{"id":1}`,
		},
		{
			name:                 "invalid input",
			inputBody:            `{"name":"Test","password":"Test"}`,
			mockBehavior:         func(s *mock_service.MockAuthorization, user model.User) {},
			expectedStatusCode:   400,
			expectedResponseBody: `{"message":"invalid input body"}`,
		},
		{
			name:      "failed to create user",
			inputBody: `{"name":"Test","username":"Test","password":"Test"}`,
			inputUser: model.User{Name: "Test", Username: "Test", Password: "Test"},
			mockBehavior: func(s *mock_service.MockAuthorization, user model.User) {
				s.EXPECT().CreateUser(user).Return(0, errFailedToCreateUser)
			},
			expectedStatusCode:   500,
			expectedResponseBody: `{"message":"failed to create user"}`,
		},
	}

	for _, testCase := range testTable {
		// Init Deps

		t.Run(testCase.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			auth := mock_service.NewMockAuthorization(c)
			testCase.mockBehavior(auth, testCase.inputUser)

			services := &service.Service{Authorization: auth}
			handler := NewHandler(services, logger)

			// Test Server
			r := gin.New()
			r.POST("sign-up", handler.singUP)

			// Test Request
			w := httptest.NewRecorder()
			req := httptest.NewRequest("POST", "/sign-up", bytes.NewBufferString(testCase.inputBody))

			// Perform request
			r.ServeHTTP(w, req)

			// Assert
			assert.Equal(t, testCase.expectedStatusCode, w.Code)
			assert.Equal(t, testCase.expectedResponseBody, w.Body.String())
		})
	}
}

func TestHandler_signIp(t *testing.T) {
	type mockBehavior func(s *mock_service.MockAuthorization, signInInput model.SingInInput)

	testTable := []struct {
		name                 string
		inputBody            string
		inputUser            model.SingInInput
		mockBehavior         mockBehavior
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name:      "OK",
			inputBody: `{"username":"Test","password":"Test"}`,
			inputUser: model.SingInInput{Username: "Test", Password: "Test"},
			mockBehavior: func(s *mock_service.MockAuthorization, signInInput model.SingInInput) {
				s.EXPECT().GenerateToken(signInInput.Username, signInInput.Password).Return("generated token", nil)
			},
			expectedStatusCode:   200,
			expectedResponseBody: `{"token":"generated token"}`,
		},
		{
			name:                 "invalid input",
			inputBody:            `{"password":"Test"}`,
			mockBehavior:         func(s *mock_service.MockAuthorization, signInInput model.SingInInput) {},
			expectedStatusCode:   400,
			expectedResponseBody: `{"message":"invalid input body"}`,
		},
		{
			name:      "failed to generate token",
			inputBody: `{"username":"Test","password":"Test"}`,
			inputUser: model.SingInInput{Username: "Test", Password: "Test"},
			mockBehavior: func(s *mock_service.MockAuthorization, signInInput model.SingInInput) {
				s.EXPECT().GenerateToken(signInInput.Username, signInInput.Password).Return("", errFailedToGenerateToken)
			},
			expectedStatusCode:   500,
			expectedResponseBody: `{"message":"failed to create token"}`,
		},
	}

	for _, testCase := range testTable {
		// Init Deps

		t.Run(testCase.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			auth := mock_service.NewMockAuthorization(c)
			testCase.mockBehavior(auth, testCase.inputUser)

			services := &service.Service{Authorization: auth}
			handler := NewHandler(services, logger)

			// Test Server
			r := gin.New()
			r.POST("sign-in", handler.singIn)

			// Test Request
			w := httptest.NewRecorder()
			req := httptest.NewRequest("POST", "/sign-in", bytes.NewBufferString(testCase.inputBody))

			// Perform request
			r.ServeHTTP(w, req)

			// Assert
			assert.Equal(t, testCase.expectedStatusCode, w.Code)
			assert.Equal(t, testCase.expectedResponseBody, w.Body.String())
		})
	}
}
