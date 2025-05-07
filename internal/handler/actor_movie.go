package handler

import (
	"encoding/json"
	"film-library/internal/model"
	"film-library/internal/service"
	"film-library/internal/utils/response"
	"fmt"
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
		response.WriteJSONError(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

// @Summary Get Actors with Their Films
// @Security ApiKeyAuth
// @Tags actor_movie
// @Description Get list of all actors with their films
// @ID get-actors-with-films
// @Accept  json
// @Produce  json
// @Success 200 {array} model.ActorWithFilms
// @Failure 400 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Failure default {object} response.ErrorResponse
// @Router /get_list_actors_films [get]
func (h *ActorMovieHandler) GetActorMovies(w http.ResponseWriter, r *http.Request) {
	if err := model.ValidateGetActors(); err != nil {
		response.WriteJSONError(w, fmt.Sprintf("%v", err.Error()), http.StatusBadRequest)
		return
	}

	actorsMap, err := h.service.GetAllActorWithFilms(r.Context())
	if err != nil {
		response.WriteJSONError(w, fmt.Sprintf("%v", err.Error()), http.StatusInternalServerError)
		return
	}

	// Преобразование map[int]model.ActorWithFilms → []model.ActorWithFilms
	actors := make([]model.ActorWithFilms, 0, len(actorsMap))
	for _, v := range actorsMap {
		actors = append(actors, v)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(actors)

}
