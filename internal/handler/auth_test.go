package handler

import (
	"bytes"
	"errors"
	"film-library/internal/model"
	"film-library/internal/service"
	mock_service "film-library/internal/service/mocks"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
)

func TestHandler_CreateUser(t *testing.T) {
	type mockBehavior func(r *mock_service.MockAuthorization, user model.User)

	tests := []struct {
		name                 string
		inputBody            string
		inputUser            model.User
		mockBehavior         mockBehavior
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name:      "Ok",
			inputBody: `{"username": "username", "password": "qwerty", "role": 1}`,
			inputUser: model.User{
				Username: "username",
				Password: "qwerty",
				Role:     1,
			},
			mockBehavior: func(r *mock_service.MockAuthorization, user model.User) {
				r.EXPECT().CreateUser(gomock.Any(), user).Return("1", nil)
			},
			expectedStatusCode:   http.StatusCreated,
			expectedResponseBody: `{"token":"1"}`,
		},
		{
			name:      "Wrong Input",
			inputBody: `{"username": "", "password": "qwerty", "role": 1}`,
			inputUser: model.User{
				Username: "",
				Password: "qwerty",
				Role:     1,
			},
			mockBehavior:         func(r *mock_service.MockAuthorization, user model.User) {},
			expectedStatusCode:   http.StatusBadRequest,
			expectedResponseBody: `{"message":"username, password and role are required"}`,
		},
		{
			name:      "Service Error",
			inputBody: `{"username": "username", "password": "qwerty", "role": 1}`,
			inputUser: model.User{
				Username: "username",
				Password: "qwerty",
				Role:     1,
			},
			mockBehavior: func(r *mock_service.MockAuthorization, user model.User) {
				r.EXPECT().CreateUser(gomock.Any(), user).Return("", errors.New("failed to create user"))
			},
			expectedStatusCode:   http.StatusInternalServerError,
			expectedResponseBody: `{"message":"failed to create user"}`,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			auth := mock_service.NewMockAuthorization(c)
			tc.mockBehavior(auth, tc.inputUser)

			services := &service.Service{Authorization: auth}
			handler := NewAuthHandler(services)

			req := httptest.NewRequest(http.MethodPost, "/sign_up", bytes.NewBufferString(tc.inputBody))
			req.Header.Set("Content-Type", "application/json")

			rr := httptest.NewRecorder()
			handler.CreateUser(rr, req)

			require.Equal(t, tc.expectedStatusCode, rr.Code)
			require.JSONEq(t, tc.expectedResponseBody, rr.Body.String())
		})
	}
}

//

//

func TestHandler_VerifyUser(t *testing.T) {
	type mockBehavior func(r *mock_service.MockAuthorization, user *model.User)

	tests := []struct {
		name                 string
		inputBody            string
		inputUser            *model.User
		mockBehavior         mockBehavior
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name:      "Ok",
			inputBody: `{"username": "username", "password": "qwerty"}`,
			inputUser: &model.User{
				Username: "username",
				Password: "qwerty",
			},
			mockBehavior: func(r *mock_service.MockAuthorization, user *model.User) {
				r.EXPECT().VerifyUser(gomock.Any(), user.Username, user.Password).Return("", user, nil)
			},
			expectedStatusCode:   http.StatusOK,
			expectedResponseBody: `{"access_token": "", "user": {"id": 0, "username": "username", "role": 0}}`,
		},
		{
			name:      "Wrong Input",
			inputBody: `{"username": "", "password": "qwerty"}`,
			inputUser: &model.User{
				Username: "",
				Password: "qwerty",
			},
			mockBehavior:         func(r *mock_service.MockAuthorization, user *model.User) {},
			expectedStatusCode:   http.StatusBadRequest,
			expectedResponseBody: `{"message": "username or password are required"}`,
		},
		{
			name:      "Service Error",
			inputBody: `{"username": "username", "password": "qwerty"}`,
			inputUser: &model.User{
				Username: "username",
				Password: "qwerty",
			},
			mockBehavior: func(r *mock_service.MockAuthorization, user *model.User) {
				r.EXPECT().VerifyUser(gomock.Any(), user.Username, user.Password).Return("", user, errors.New("failed to verify user"))
			},
			expectedStatusCode:   http.StatusInternalServerError,
			expectedResponseBody: `{"message": "failed to verify user"}`,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			auth := mock_service.NewMockAuthorization(c)
			tc.mockBehavior(auth, tc.inputUser)

			services := &service.Service{Authorization: auth}
			handler := NewAuthHandler(services)

			req := httptest.NewRequest(http.MethodPost, "/sign_in", bytes.NewBufferString(tc.inputBody))
			req.Header.Set("Content-Type", "application/json")

			rr := httptest.NewRecorder()
			handler.VerifyUser(rr, req)

			require.Equal(t, tc.expectedStatusCode, rr.Code)
			require.JSONEq(t, tc.expectedResponseBody, rr.Body.String())
		})
	}
}
