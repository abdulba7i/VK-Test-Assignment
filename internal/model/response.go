package model

import "time"

type OKResponse struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
	Data    any    `json:"data"`
}

type ErrorResponse struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
	Code    int    `json:"code"`
}

type ActorWithMovies struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	Gender    string    `json:"gender"`
	BirthDate time.Time `json:"birth_date"`
	Movies    []Film    `json:"movies"`
}

type ActorsListResponse struct {
	OKResponse
	Data []ActorWithMovies `json:"data"`
}

type MovieWithActors struct {
	ID          int       `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	ReleaseDate time.Time `json:"release_date"`
	Rating      float64   `json:"rating"`
	Actors      []Actor   `json:"actors"` // Список актёров
}

type MoviesListResponse struct {
	OKResponse
	Data []MovieWithActors `json:"data"`
}
