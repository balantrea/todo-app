package handler

import (
	"fmt"
	"net/http/httptest"
	"testing"

	"github.com/balantrea/todo-app/internal/service"
	"github.com/balantrea/todo-app/internal/service/mocks"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/assert/v2"
	"go.uber.org/mock/gomock"
)

func TestHandler_userIdentity(t *testing.T) {
	type mockBehavior func(s *mock_service.MockAuthorization, token string)

	testTable := []struct {
		name                 string
		headerName           string
		headerValue          string
		token                string
		mockBehavior         mockBehavior
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name:        "OK",
			headerName:  "Authorization",
			headerValue: "Bearer token",
			mockBehavior: func(s *mock_service.MockAuthorization, token string) {
				s.EXPECT().ParseToken("token").Return(1, nil)
			},
			expectedStatusCode:   200,
			expectedResponseBody: "1",
		},
		{
			name:                 "invalid header",
			headerName:           "Authorization",
			mockBehavior:         func(s *mock_service.MockAuthorization, token string) {},
			expectedStatusCode:   401,
			expectedResponseBody: `{"message":"empty auth header"}`,
		},
		{
			name:        "failed to user identity",
			headerName:  "Authorization",
			headerValue: "Bearer token",
			mockBehavior: func(s *mock_service.MockAuthorization, token string) {
				s.EXPECT().ParseToken("token").Return(0, errInvalidToken)
			},
			expectedStatusCode:   401,
			expectedResponseBody: `{"message":"invalid token"}`,
		},
	}

	for _, testCase := range testTable {
		// Init Deps

		t.Run(testCase.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			auth := mock_service.NewMockAuthorization(c)
			testCase.mockBehavior(auth, testCase.token)

			services := &service.Service{Authorization: auth}
			handler := NewHandler(services, logger)

			// Test Server
			r := gin.New()
			r.GET("/protected", handler.userIdentity, func(c *gin.Context) {
				id, _ := c.Get(userCtx)
				c.String(200, fmt.Sprintf("%d", id.(int)))
			})

			// Test Request

			w := httptest.NewRecorder()
			req := httptest.NewRequest("GET", "/protected", nil)
			req.Header.Set(testCase.headerName, testCase.headerValue)

			// Perform Request
			r.ServeHTTP(w, req)

			// Assert
			assert.Equal(t, testCase.expectedStatusCode, w.Code)
			assert.Equal(t, testCase.expectedResponseBody, w.Body.String())
		})
	}
}

func TestHandler_getUserID(t *testing.T) {

	testTable := []struct {
		name                 string
		setupContext         func(c *gin.Context)
		expectedID           int
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name: "OK",
			setupContext: func(c *gin.Context) {
				c.Set(userCtx, 1)
			},
			expectedID:           1,
			expectedStatusCode:   200,
			expectedResponseBody: ``,
		},
		{
			name:                 "user ID is not found",
			setupContext:         func(c *gin.Context) {},
			expectedID:           0,
			expectedStatusCode:   400,
			expectedResponseBody: `{"message":"user is not found"}`,
		},
		{
			name: "failed to assert user id",
			setupContext: func(c *gin.Context) {
				c.Set(userCtx, "1")
			},
			expectedID:           0,
			expectedStatusCode:   400,
			expectedResponseBody: `{"message":"user is not found"}`,
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)

			handler := &Handler{
				logger: logger,
			}

			testCase.setupContext(c)

			id, _ := handler.getUserId(c)

			assert.Equal(t, testCase.expectedID, id)
			assert.Equal(t, testCase.expectedStatusCode, w.Code)
			assert.Equal(t, testCase.expectedResponseBody, w.Body.String())
		})
	}
}
