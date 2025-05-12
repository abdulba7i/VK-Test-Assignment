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

func TestHandler_CreateMovie(t *testing.T) {
	type mockBehavior func(r *mock_service.MockMovie, actor model.Film)

	tests := []struct {
		name                 string
		inputBody            string
		inputUser            model.Film
		mockBehavior         mockBehavior
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name:      "Ok",
			inputBody: `{"name": "name", "description": "description", "release_date": "2004-04-10T21:12:05+03:00", "rating": 9.3}`,
			inputUser: model.Film{
				Name:        "name",
				Description: "description",
				Releasedate: parseTime("2004-04-10T21:12:05+03:00"),
				Rating:      9.3,
			},
			mockBehavior: func(r *mock_service.MockMovie, actor model.Film) {
				r.EXPECT().AddMovie(gomock.Any(), actor).Return(nil)
			},
			expectedStatusCode:   http.StatusCreated,
			expectedResponseBody: `{"id": 0, "name": "name", "description": "description", "release_date": "2004-04-10T21:12:05+03:00", "rating": 9.3, "list_actors": null}`,
		},
		{
			name:      "Wrong input NAME",
			inputBody: `{"name": "", "description": "description", "release_date": "2004-04-10T21:12:05+03:00", "rating": 9.3}`,
			inputUser: model.Film{
				Name:        "",
				Description: "description",
				Releasedate: parseTime("2004-04-10T21:12:05+03:00"),
				Rating:      9.3,
			},
			mockBehavior:         func(r *mock_service.MockMovie, actor model.Film) {},
			expectedStatusCode:   http.StatusBadRequest,
			expectedResponseBody: `{"message": "название фильма должно быть от 1 до 150 символов"}`,
		},
		{
			name:      "Service Error",
			inputBody: `{"name": "name", "description": "description", "release_date": "2004-04-10T21:12:05+03:00", "rating": 9.3}`,
			inputUser: model.Film{
				Name:        "name",
				Description: "description",
				Releasedate: parseTime("2004-04-10T21:12:05+03:00"),
				Rating:      9.3,
			},
			mockBehavior: func(r *mock_service.MockMovie, actor model.Film) {
				r.EXPECT().AddMovie(gomock.Any(), actor).Return(errors.New("Failed to create film"))
			},
			expectedStatusCode:   http.StatusInternalServerError,
			expectedResponseBody: `{"message": "Failed to create film"}`,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			auth := mock_service.NewMockMovie(c)
			tc.mockBehavior(auth, tc.inputUser)

			services := &service.Service{Movie: auth}
			handler := NewMovieHandler(services)

			req := httptest.NewRequest(http.MethodPost, "/films", bytes.NewBufferString(tc.inputBody))
			req.Header.Set("Content-Type", "application/json")

			rr := httptest.NewRecorder()
			handler.CreateFilm(rr, req)

			require.Equal(t, tc.expectedStatusCode, rr.Code)
			require.JSONEq(t, tc.expectedResponseBody, rr.Body.String())
		})
	}
}

func TestHandler_UpdateMovie(t *testing.T) {
	type mockBehavior func(r *mock_service.MockMovie, actor model.Film)

	tests := []struct {
		name                 string
		inputBody            string
		inputUser            model.Film
		mockBehavior         mockBehavior
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name:      "Ok",
			inputBody: `{"name": "name", "description": "description", "release_date": "2004-04-10T21:12:05+03:00", "rating": 9.3}`,
			inputUser: model.Film{
				Name:        "name",
				Description: "description",
				Releasedate: parseTime("2004-04-10T21:12:05+03:00"),
				Rating:      9.3,
			},
			mockBehavior: func(r *mock_service.MockMovie, actor model.Film) {
				r.EXPECT().UpdateMovie(gomock.Any(), actor).Return(nil)
			},
			expectedStatusCode:   http.StatusCreated,
			expectedResponseBody: `{"id": 0, "name": "name", "description": "description", "release_date": "2004-04-10T21:12:05+03:00", "rating": 9.3, "list_actors": null}`,
		},
		{
			name:      "Wrong input NAME",
			inputBody: `{"name": "", "description": "description", "release_date": "2004-04-10T21:12:05+03:00", "rating": 9.3}`,
			inputUser: model.Film{
				Name:        "",
				Description: "description",
				Releasedate: parseTime("2004-04-10T21:12:05+03:00"),
				Rating:      9.3,
			},
			mockBehavior:         func(r *mock_service.MockMovie, actor model.Film) {},
			expectedStatusCode:   http.StatusBadRequest,
			expectedResponseBody: `{"message": "название фильма должно быть от 1 до 150 символов"}`,
		},
		{
			name:      "Service Error",
			inputBody: `{"name": "name", "description": "description", "release_date": "2004-04-10T21:12:05+03:00", "rating": 9.3}`,
			inputUser: model.Film{
				Name:        "name",
				Description: "description",
				Releasedate: parseTime("2004-04-10T21:12:05+03:00"),
				Rating:      9.3,
			},
			mockBehavior: func(r *mock_service.MockMovie, actor model.Film) {
				r.EXPECT().UpdateMovie(gomock.Any(), actor).Return(errors.New("Failed to update film"))
			},
			expectedStatusCode:   http.StatusInternalServerError,
			expectedResponseBody: `{"message": "Failed to update film"}`,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			auth := mock_service.NewMockMovie(c)
			tc.mockBehavior(auth, tc.inputUser)

			services := &service.Service{Movie: auth}
			handler := NewMovieHandler(services)

			req := httptest.NewRequest(http.MethodPost, "/films", bytes.NewBufferString(tc.inputBody))
			req.Header.Set("Content-Type", "application/json")

			rr := httptest.NewRecorder()
			handler.UpdateFilm(rr, req)

			require.Equal(t, tc.expectedStatusCode, rr.Code)
			require.JSONEq(t, tc.expectedResponseBody, rr.Body.String())
		})
	}
}

func TestHandler_DeleteMovie(t *testing.T) {
	type mockBehavior func(r *mock_service.MockMovie, id int)

	tests := []struct {
		name                 string
		queryParam           string
		mockBehavior         mockBehavior
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name:       "Ok",
			queryParam: `0`,
			mockBehavior: func(r *mock_service.MockMovie, id int) {
				r.EXPECT().DeleteMovie(gomock.Any(), id).Return(nil)
			},
			expectedStatusCode:   http.StatusOK,
			expectedResponseBody: `{"message": "movie deleted successfully"}`,
		},
		{
			name:                 "Wrong input ID",
			queryParam:           `first`,
			mockBehavior:         func(r *mock_service.MockMovie, id int) {},
			expectedStatusCode:   http.StatusBadRequest,
			expectedResponseBody: `{"message": "Invalid film ID"}`,
		},
		{
			name:       "Error Service",
			queryParam: `0`,
			mockBehavior: func(r *mock_service.MockMovie, id int) {
				r.EXPECT().DeleteMovie(gomock.Any(), id).Return(errors.New("Failed to delete film"))
			},
			expectedStatusCode:   http.StatusInternalServerError,
			expectedResponseBody: `{"message": "Failed to delete film"}`,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			auth := mock_service.NewMockMovie(c)
			tc.mockBehavior(auth, 0)

			services := &service.Service{Movie: auth}
			handler := NewMovieHandler(services)

			req := httptest.NewRequest(http.MethodDelete, "/film_del/"+tc.queryParam, nil)
			req.Header.Set("Content-Type", "application/json")

			rr := httptest.NewRecorder()
			handler.DeleteFilm(rr, req)

			require.Equal(t, tc.expectedStatusCode, rr.Code)
			require.JSONEq(t, tc.expectedResponseBody, rr.Body.String())
		})
	}
}

func TestHandler_GetAllFilms(t *testing.T) {
	type mockBehavior func(r *mock_service.MockMovie, sortBy string)

	tests := []struct {
		name                 string
		queryParam           string
		mockBehavior         mockBehavior
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name:       "Ok",
			queryParam: "name",
			mockBehavior: func(r *mock_service.MockMovie, sortBy string) {
				r.EXPECT().GetFilms(gomock.Any(), sortBy).Return([]model.Film{}, nil)
			},
			expectedStatusCode:   http.StatusOK,
			expectedResponseBody: `[]`,
		},
		{
			name:                 "Wrong Input",
			queryParam:           "actor",
			mockBehavior:         func(r *mock_service.MockMovie, sortBy string) {},
			expectedStatusCode:   http.StatusBadRequest,
			expectedResponseBody: `{"message": "Некорректная сортировка"}`,
		},
		{
			name:       "Server Error",
			queryParam: "name",
			mockBehavior: func(r *mock_service.MockMovie, sortBy string) {
				r.EXPECT().GetFilms(gomock.Any(), sortBy).Return([]model.Film{}, errors.New("Failed to get films"))
			},
			expectedStatusCode:   http.StatusInternalServerError,
			expectedResponseBody: `{"message": "Failed to get films"}`,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			auth := mock_service.NewMockMovie(c)
			tc.mockBehavior(auth, tc.queryParam)

			services := &service.Service{Movie: auth}
			handler := NewMovieHandler(services)

			req := httptest.NewRequest(http.MethodGet, "/films_get_list/?sort_by="+tc.queryParam, nil)
			req.Header.Set("Content-Type", "application/json")

			rr := httptest.NewRecorder()
			handler.GetAllFilms(rr, req)

			require.Equal(t, tc.expectedStatusCode, rr.Code)
			require.JSONEq(t, tc.expectedResponseBody, rr.Body.String())
		})
	}
}

func TestHandler_SearchFilm(t *testing.T) {
	type mockBehavior func(r *mock_service.MockMovie, actor, movie string)

	tests := []struct {
		name                 string
		queryParamActor      string
		queryParamMovie      string
		mockBehavior         mockBehavior
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name:            "OK",
			queryParamActor: "Дизель",
			queryParamMovie: "Форсаж",
			mockBehavior: func(r *mock_service.MockMovie, actor, movie string) {
				r.EXPECT().SearchFilm(gomock.Any(), actor, movie).Return(model.Film{}, nil)
			},
			expectedStatusCode:   http.StatusOK,
			expectedResponseBody: `{"id": 0, "name": "", "description": "", "rating": 0, "release_date": "0001-01-01T00:00:00Z", "list_actors": null}`,
		},
		{
			name:                 "OK",
			queryParamActor:      "",
			queryParamMovie:      "",
			mockBehavior:         func(r *mock_service.MockMovie, actor, movie string) {},
			expectedStatusCode:   http.StatusBadRequest,
			expectedResponseBody: `{"message": "необходимо указать либо название фильма, либо имя актёра"}`,
		},

		{
			name:            "Service Error",
			queryParamActor: "Дизель",
			queryParamMovie: "ТрансфФорсажормеры",
			mockBehavior: func(r *mock_service.MockMovie, actor, movie string) {
				r.EXPECT().SearchFilm(gomock.Any(), actor, movie).Return(model.Film{}, errors.New("Search failed: ошибка поиска: sql: no rows in result set"))
			},
			expectedStatusCode:   http.StatusInternalServerError,
			expectedResponseBody: `{"message":"Search failed: Search failed: ошибка поиска: sql: no rows in result set"}`,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			auth := mock_service.NewMockMovie(c)
			tc.mockBehavior(auth, tc.queryParamActor, tc.queryParamMovie)

			services := &service.Service{Movie: auth}
			handler := NewMovieHandler(services)

			req := httptest.NewRequest(http.MethodGet, "/films/search?actor="+tc.queryParamActor+"&movie="+tc.queryParamMovie, nil)
			req.Header.Set("Content-Type", "application/json")

			rr := httptest.NewRecorder()
			handler.SearchFilm(rr, req)

			require.Equal(t, tc.expectedStatusCode, rr.Code)
			require.JSONEq(t, tc.expectedResponseBody, rr.Body.String())
		})
	}
}
