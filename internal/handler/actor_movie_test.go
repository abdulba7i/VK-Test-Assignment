package handler

import (
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

func TestHandler_GetActorMovies(t *testing.T) {
	type mockBehavior func(r *mock_service.MockActorMovie)

	tests := []struct {
		name                 string
		mockBehavior         mockBehavior
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name: "OK",
			mockBehavior: func(r *mock_service.MockActorMovie) {
				r.EXPECT().GetAllActorWithFilms(gomock.Any()).Return(map[int]model.ActorWithFilms{}, nil)
			},
			expectedStatusCode:   http.StatusOK,
			expectedResponseBody: `{}`,
		},
		{
			name: "Service Error",
			mockBehavior: func(r *mock_service.MockActorMovie) {
				r.EXPECT().GetAllActorWithFilms(gomock.Any()).Return(map[int]model.ActorWithFilms{}, errors.New("error"))
			},
			expectedStatusCode:   http.StatusInternalServerError,
			expectedResponseBody: `{"message": "error"}`,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			auth := mock_service.NewMockActorMovie(c)
			tc.mockBehavior(auth)

			services := &service.Service{ActorMovie: auth}
			handler := NewActorMovieHandler(services)

			req := httptest.NewRequest(http.MethodGet, "/actors_films", nil)
			req.Header.Set("Content-Type", "application/json")

			rr := httptest.NewRecorder()
			handler.GetActorMovies(rr, req)

			require.Equal(t, tc.expectedStatusCode, rr.Code)
			require.JSONEq(t, tc.expectedResponseBody, rr.Body.String())
		})
	}
}
