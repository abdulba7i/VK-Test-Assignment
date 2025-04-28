package handler

import (
	"net/http"
	"testing"
	"time"
)

func Test_handler_ActorCreate(t *testing.T) {
	cases := []struct {
		nameCases      string
		name           string
		gender         string
		date_of_birth  time.Time
		wantStatusCode int
		wantErr        bool
	}{
		{
			nameCases:      "Success",
			name:           "John",
			gender:         "male",
			date_of_birth:  parseDate("2025-04-10T21:12:05+03:00"),
			wantStatusCode: http.StatusCreated,
			wantErr:        false,
		},
	}

	for _, tc := range cases {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			// actorCreateMock := mocks.New
		})
	}
}

func parseDate(dateStr string) time.Time {
	t, err := time.Parse(time.RFC3339, dateStr)
	if err != nil {
		panic(err) // Обработка ошибки
	}
	return t
}
