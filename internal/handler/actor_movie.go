package handler

import (
	"encoding/json"
	"film-library/internal/model"
	"film-library/internal/service"
	"net/http"
)

type ActorMovieHandler struct {
	service service.ActorMovie
}

func NewActorMovieHandler(service service.ActorMovie) ActorMovieHandler {
	return ActorMovieHandler{service: service}
}

func (h *ActorMovieHandler) HandleActorMovieGet(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		h.GetActorMovies(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func (h *ActorMovieHandler) GetActorMovies(w http.ResponseWriter, r *http.Request) {
	if err := model.ValidateGetActors(); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	actors, err := h.service.GetAllActorWithFilms(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(actors)
}
