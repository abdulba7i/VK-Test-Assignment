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
	"time"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
)

func TestHandler_ActorCreate(t *testing.T) {
	type mockBehavior func(r *mock_service.MockActor, actor model.Actor)

	tests := []struct {
		name                 string
		inputBody            string
		inputUser            model.Actor
		mockBehavior         mockBehavior
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name:      "Ok",
			inputBody: `{"name": "name", "gender": "male", "date_of_birth": "2004-04-10T21:12:05+03:00"}`,
			inputUser: model.Actor{
				Name:        "name",
				Gender:      "male",
				DateOfBirth: parseTime("2004-04-10T21:12:05+03:00"),
			},
			mockBehavior: func(r *mock_service.MockActor, actor model.Actor) {
				r.EXPECT().AddActor(gomock.Any(), actor).Return(nil)
			},
			expectedStatusCode:   http.StatusCreated,
			expectedResponseBody: `{"id": 0, "name": "name", "gender": "male", "date_of_birth": "2004-04-10T21:12:05+03:00"}`,
		},

		{
			name:      "Wrong Input NAME",
			inputBody: `{"name": "", "gender": "male", "date_of_birth": "2004-04-10T21:12:05+03:00"}`,
			inputUser: model.Actor{
				Name:        "",
				Gender:      "male",
				DateOfBirth: parseTime("2004-04-10T21:12:05+03:00"),
			},
			mockBehavior:         func(r *mock_service.MockActor, actor model.Actor) {},
			expectedStatusCode:   http.StatusBadRequest,
			expectedResponseBody: `{"message": "имя актера не может быть пустым и не должно превышать 100 символов"}`,
		},

		{
			name:      "Wrong Input GENDER",
			inputBody: `{"name": "name", "gender": "", "date_of_birth": "2004-04-10T21:12:05+03:00"}`,
			inputUser: model.Actor{
				Name:        "name",
				Gender:      "",
				DateOfBirth: parseTime("2004-04-10T21:12:05+03:00"),
			},
			mockBehavior:         func(r *mock_service.MockActor, actor model.Actor) {},
			expectedStatusCode:   http.StatusBadRequest,
			expectedResponseBody: `{"message": "не указан пол"}`,
		},

		{
			name:      "Wrong Input DATE",
			inputBody: `{"name": "name", "gender": "male", "date_of_birth": "2025-04-10T21:12:05+03:00"}`,
			inputUser: model.Actor{
				Name:        "name",
				Gender:      "male",
				DateOfBirth: parseTime("2025-04-10T21:12:05+03:00"),
			},
			mockBehavior:         func(r *mock_service.MockActor, actor model.Actor) {},
			expectedStatusCode:   http.StatusBadRequest,
			expectedResponseBody: `{"message": "актёр должен быть старше 5 лет"}`,
		},

		{
			name:      "Service Error",
			inputBody: `{"name": "name", "gender": "male", "date_of_birth": "2004-04-10T21:12:05+03:00"}`,
			inputUser: model.Actor{
				Name:        "name",
				Gender:      "male",
				DateOfBirth: parseTime("2004-04-10T21:12:05+03:00"),
			},
			mockBehavior: func(r *mock_service.MockActor, actor model.Actor) {
				r.EXPECT().AddActor(gomock.Any(), actor).Return(errors.New("Failed to create actor"))
			},
			expectedStatusCode:   http.StatusInternalServerError,
			expectedResponseBody: `{"message": "Failed to create actor"}`,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			auth := mock_service.NewMockActor(c)
			tc.mockBehavior(auth, tc.inputUser)

			services := &service.Service{Actor: auth}
			handler := NewActorHandler(services)

			req := httptest.NewRequest(http.MethodPost, "/actors", bytes.NewBufferString(tc.inputBody))
			req.Header.Set("Content-Type", "application/json")

			rr := httptest.NewRecorder()
			handler.CreateActor(rr, req)

			require.Equal(t, tc.expectedStatusCode, rr.Code)
			require.JSONEq(t, tc.expectedResponseBody, rr.Body.String())
		})
	}
}

func TestHandler_UpdateActor(t *testing.T) {
	type mockBehavior func(r *mock_service.MockActor, actor model.Actor)

	tests := []struct {
		name                 string
		inputBody            string
		inputUser            model.Actor
		mockBehavior         mockBehavior
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name:      "Ok",
			inputBody: `{"name": "name", "gender": "male", "date_of_birth": "2004-04-10T21:12:05+03:00"}`,
			inputUser: model.Actor{
				Name:        "name",
				Gender:      "male",
				DateOfBirth: parseTime("2004-04-10T21:12:05+03:00"),
			},
			mockBehavior: func(r *mock_service.MockActor, actor model.Actor) {
				r.EXPECT().UpdateActor(gomock.Any(), actor).Return(nil)
			},
			expectedStatusCode:   http.StatusCreated,
			expectedResponseBody: `{"id": 0, "name": "name", "gender": "male", "date_of_birth": "2004-04-10T21:12:05+03:00"}`,
		},

		{
			name:      "Wrong Input NAME",
			inputBody: `{"name": "", "gender": "male", "date_of_birth": "2004-04-10T21:12:05+03:00"}`,
			inputUser: model.Actor{
				Name:        "",
				Gender:      "male",
				DateOfBirth: parseTime("2004-04-10T21:12:05+03:00"),
			},
			mockBehavior:         func(r *mock_service.MockActor, actor model.Actor) {},
			expectedStatusCode:   http.StatusBadRequest,
			expectedResponseBody: `{"message": "имя актера не может быть пустым и не должно превышать 100 символов"}`,
		},

		{
			name:      "Wrong Input GENDER",
			inputBody: `{"name": "name", "gender": "", "date_of_birth": "2004-04-10T21:12:05+03:00"}`,
			inputUser: model.Actor{
				Name:        "name",
				Gender:      "",
				DateOfBirth: parseTime("2004-04-10T21:12:05+03:00"),
			},
			mockBehavior:         func(r *mock_service.MockActor, actor model.Actor) {},
			expectedStatusCode:   http.StatusBadRequest,
			expectedResponseBody: `{"message": "не указан пол"}`,
		},

		{
			name:      "Wrong Input DATE",
			inputBody: `{"name": "name", "gender": "male", "date_of_birth": "2025-04-10T21:12:05+03:00"}`,
			inputUser: model.Actor{
				Name:        "name",
				Gender:      "male",
				DateOfBirth: parseTime("2025-04-10T21:12:05+03:00"),
			},
			mockBehavior:         func(r *mock_service.MockActor, actor model.Actor) {},
			expectedStatusCode:   http.StatusBadRequest,
			expectedResponseBody: `{"message": "актёр должен быть старше 5 лет"}`,
		},

		{
			name:      "Service Error",
			inputBody: `{"name": "name", "gender": "male", "date_of_birth": "2004-04-10T21:12:05+03:00"}`,
			inputUser: model.Actor{
				Name:        "name",
				Gender:      "male",
				DateOfBirth: parseTime("2004-04-10T21:12:05+03:00"),
			},
			mockBehavior: func(r *mock_service.MockActor, actor model.Actor) {
				r.EXPECT().UpdateActor(gomock.Any(), actor).Return(errors.New("Failed to update actor"))
			},
			expectedStatusCode:   http.StatusInternalServerError,
			expectedResponseBody: `{"message": "Failed to update actor"}`,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			auth := mock_service.NewMockActor(c)
			tc.mockBehavior(auth, tc.inputUser)

			services := &service.Service{Actor: auth}
			handler := NewActorHandler(services)

			req := httptest.NewRequest(http.MethodPut, "/actor/", bytes.NewBufferString(tc.inputBody))
			req.Header.Set("Content-Type", "application/json")

			rr := httptest.NewRecorder()
			handler.UpdateActor(rr, req)

			require.Equal(t, tc.expectedStatusCode, rr.Code)
			require.JSONEq(t, tc.expectedResponseBody, rr.Body.String())
		})
	}
}

func TestHandler_DeleteActor(t *testing.T) {
	type mockBehavior func(r *mock_service.MockActor, id int)

	tests := []struct {
		name                 string
		queryParam           string
		mockBehavior         mockBehavior
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name:       "Ok",
			queryParam: `id=0`,
			mockBehavior: func(r *mock_service.MockActor, id int) {
				r.EXPECT().DeleteActor(gomock.Any(), id).Return(nil)
			},
			expectedStatusCode:   http.StatusOK,
			expectedResponseBody: `{"message": "actor deleted successfully"}`,
		},
		{
			name:                 "Wrong input ID",
			queryParam:           `id=first`,
			mockBehavior:         func(r *mock_service.MockActor, id int) {},
			expectedStatusCode:   http.StatusBadRequest,
			expectedResponseBody: `{"message": "Invalid actor ID"}`,
		},
		{
			name:       "Error Service",
			queryParam: `id=0`,
			mockBehavior: func(r *mock_service.MockActor, id int) {
				r.EXPECT().DeleteActor(gomock.Any(), id).Return(errors.New("Failed to delete actor"))
			},
			expectedStatusCode:   http.StatusInternalServerError,
			expectedResponseBody: `{"message": "Failed to delete actor"}`,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			auth := mock_service.NewMockActor(c)
			tc.mockBehavior(auth, 0)

			services := &service.Service{Actor: auth}
			handler := NewActorHandler(services)

			req := httptest.NewRequest(http.MethodPut, "/actor_del/?"+tc.queryParam, nil)
			req.Header.Set("Content-Type", "application/json")

			rr := httptest.NewRecorder()
			handler.DeleteActor(rr, req)

			require.Equal(t, tc.expectedStatusCode, rr.Code)
			require.JSONEq(t, tc.expectedResponseBody, rr.Body.String())
		})
	}
}

func parseTime(timeStr string) time.Time {
	t, err := time.Parse(time.RFC3339, timeStr)
	if err != nil {
		panic(err)
	}
	return t
}
