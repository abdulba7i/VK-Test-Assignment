package handler

import (
	"encoding/json"
	"film-library/internal/model"
	"film-library/internal/service"
	authmid "film-library/internal/utils/auth_mid"
	"film-library/internal/utils/response"
	"fmt"
	"net/http"
	"strconv"
	"strings"
)

type MovieHandler struct {
	service service.Movie
}

func NewMovieHandler(service service.Movie) MovieHandler {
	return MovieHandler{service: service}
}

func (h *MovieHandler) HandleMovieGet(w http.ResponseWriter, r *http.Request) {
	switch {
	case r.URL.Path == "/films" && r.Method == http.MethodGet:
		h.GetAllFilms(w, r)
	case r.URL.Path == "/films/search" && r.Method == http.MethodGet:
		h.SearchFilm(w, r)
	default:
		response.WriteJSONError(w, "Method not allowed", http.StatusNotFound)
	}
}
func (h *MovieHandler) HandleMoviePost(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		h.CreateFilm(w, r)
	default:
		response.WriteJSONError(w, "Method not allowed", http.StatusNotFound)
	}
}

func (h *MovieHandler) HandleMoviePut(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPut:
		h.UpdateFilm(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func (h *MovieHandler) HandleMovieDelete(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodDelete:
		h.DeleteFilm(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

// @Summary Create Film
// @Security ApiKeyAuth
// @Tags film
// @Description Create Film
// @ID create-film
// @Accept  json
// @Produce  json
// @Param film body model.Film true "Create Film"
// @Success 201 {object} model.Film
// @Failure 400,403 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Failure default {object} response.ErrorResponse
// @Router /film_create [post]
func (h *MovieHandler) CreateFilm(w http.ResponseWriter, r *http.Request) {
	if authmid.IsAdmin(r) {
		response.WriteJSONError(w, "Forbidden: admin access required", http.StatusForbidden)
		return
	}

	var film model.Film
	if err := json.NewDecoder(r.Body).Decode(&film); err != nil {
		response.WriteJSONError(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if err := film.Validate(); err != nil {
		response.WriteJSONError(w, fmt.Sprintf("%v", err.Error()), http.StatusBadRequest)
		return
	}

	err := h.service.AddMovie(r.Context(), film)
	if err != nil {
		response.WriteJSONError(w, "Failed to create film", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(film)
}

// @Summary Update Film
// @Security ApiKeyAuth
// @Tags film
// @Description Update Film
// @ID update-film
// @Accept  json
// @Produce  json
// @Param film body model.Film true "Update Film"
// @Success 201 {object} model.Film
// @Failure 400,403 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Failure default {object} response.ErrorResponse
// @Router /film_update [put]
func (h *MovieHandler) UpdateFilm(w http.ResponseWriter, r *http.Request) {
	if authmid.IsAdmin(r) {
		response.WriteJSONError(w, "Forbidden: admin access required", http.StatusForbidden)
		return
	}

	var film model.Film
	if err := json.NewDecoder(r.Body).Decode(&film); err != nil {
		response.WriteJSONError(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if err := film.Validate(); err != nil {
		response.WriteJSONError(w, fmt.Sprintf("%v", err.Error()), http.StatusBadRequest)
		return
	}

	err := h.service.UpdateMovie(r.Context(), film)
	if err != nil {
		response.WriteJSONError(w, "Failed to update film", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(film)
}

// @Summary Delete Film
// @Security ApiKeyAuth
// @Tags film
// @Description Delete Film
// @ID delete-film
// @Accept  json
// @Produce  json
// @Param id query int true "Film ID"
// @Success 200 {object} map[string]string
// @Failure 400,403 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Failure default {object} response.ErrorResponse
// @Router /film_delete/{id} [delete]
func (h *MovieHandler) DeleteFilm(w http.ResponseWriter, r *http.Request) {
	if authmid.IsAdmin(r) {
		response.WriteJSONError(w, "Forbidden: admin access required", http.StatusForbidden)
		return
	}

	parts := strings.Split(strings.Trim(r.URL.Path, "/"), "/")
	if len(parts) < 2 {
		response.WriteJSONError(w, "Missing film ID", http.StatusBadRequest)
		return
	}

	idStr := parts[len(parts)-1]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		response.WriteJSONError(w, "Invalid film ID", http.StatusBadRequest)
		return
	}

	err = h.service.DeleteMovie(r.Context(), id)
	if err != nil {
		response.WriteJSONError(w, "Failed to delete film", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	response := map[string]string{
		"message": "movie deleted successfully",
	}

	json.NewEncoder(w).Encode(response)
}

// @Summary Get All Films
// @Security ApiKeyAuth
// @Tags film
// @Description Get all films with optional sorting
// @ID get-all-films
// @Accept  json
// @Produce  json
// @Param sort_by query string false "Sort films by field (optional)"
// @Success 200 {array} model.Film
// @Failure 400 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Failure default {object} response.ErrorResponse
// @Router /films_get_list [get]
func (h *MovieHandler) GetAllFilms(w http.ResponseWriter, r *http.Request) {
	var listFilms []model.Film
	sortBy := r.URL.Query().Get("sort_by")

	if err := model.ValidateSortFilm(sortBy); err != nil {
		response.WriteJSONError(w, fmt.Sprintf("%v", err.Error()), http.StatusBadRequest)
		return
	}

	listFilms, err := h.service.GetFilms(r.Context(), sortBy)
	if err != nil {
		response.WriteJSONError(w, "Failed to get films", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(listFilms)
}

// @Summary Search Film
// @Security ApiKeyAuth
// @Tags film
// @Description Search films by actor and/or movie name
// @ID search-film
// @Accept  json
// @Produce  json
// @Param actor query string false "Actor name to search for"
// @Param movie query string false "Movie title to search for"
// @Success 200 {array} model.Film
// @Failure 400 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Failure default {object} response.ErrorResponse
// @Router /films/search [get]
func (h *MovieHandler) SearchFilm(w http.ResponseWriter, r *http.Request) {
	var filmSearch model.Film

	actor, movie := strings.TrimSpace(r.URL.Query().Get("actor")), strings.TrimSpace(r.URL.Query().Get("movie"))

	if err := filmSearch.ValidateFilmSearchParams(movie, actor); err != nil {
		response.WriteJSONError(w, fmt.Sprintf("%v", err.Error()), http.StatusBadRequest)
		return
	}

	films, err := h.service.SearchFilm(r.Context(), actor, movie)
	if err != nil {
		response.WriteJSONError(w, fmt.Sprintf("Search failed: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(films)
}
