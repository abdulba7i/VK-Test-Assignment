package handler

import (
	"film-library/internal/service"
	"net/http"
)

type ActorMovieHandler struct {
	service service.ActorMovieService
}

func NewActorMovieHandler(service service.ActorMovieService) ActorMovieHandler {
	return ActorMovieHandler{service: service}
}

func (h *ActorMovieHandler) GetAllFilms(w http.ResponseWriter, r *http.Request) {
	// var actor model.ActorWithFilms
	_ = r.URL.Query().Get("sort_by")

	// if actor.V
}
