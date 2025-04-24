package handler

import (
	"encoding/json"
	"film-library/internal/model"
	"film-library/internal/service"
	"film-library/internal/utils"
	"fmt"
	"net/http"
	"strconv"
	"strings"
)

type MovieHandler struct {
	service service.MovieService
}

func NewMovieHandler(service service.MovieService) MovieHandler {
	return MovieHandler{service: service}
}

func (h *MovieHandler) HandleMovieGet(w http.ResponseWriter, r *http.Request) {
	switch {
	case r.URL.Path == "/films_get_list" && r.Method == http.MethodGet:
		h.GetAllFilms(w, r)
	case r.URL.Path == "/films/search" && r.Method == http.MethodGet:
		h.SearchFilm(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusNotFound)
	}
}
func (h *MovieHandler) HandleMoviePost(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		h.CreateFilm(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func (h *MovieHandler) HandleMoviePut(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPut:
		h.UpdateFilm(w, r)
	case http.MethodDelete:
		h.DeleteFilm(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func (h *MovieHandler) CreateFilm(w http.ResponseWriter, r *http.Request) {
	if utils.IsAdmin(r) {
		http.Error(w, "Forbidden: admin access required", http.StatusForbidden)
		return
	}

	var film model.Film
	if err := json.NewDecoder(r.Body).Decode(&film); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if err := film.Validate(); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err := h.service.AddMovie(r.Context(), film)
	if err != nil {
		http.Error(w, "Failed to create film", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(film)
}

func (h *MovieHandler) UpdateFilm(w http.ResponseWriter, r *http.Request) {
	if utils.IsAdmin(r) {
		http.Error(w, "Forbidden: admin access required", http.StatusForbidden)
		return
	}

	var film model.Film
	if err := json.NewDecoder(r.Body).Decode(&film); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if err := film.Validate(); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err := h.service.UpdateMovie(r.Context(), film)
	if err != nil {
		http.Error(w, "Failed to update film", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(film)
}

func (h *MovieHandler) DeleteFilm(w http.ResponseWriter, r *http.Request) {
	if utils.IsAdmin(r) {
		http.Error(w, "Forbidden: admin access required", http.StatusForbidden)
		return
	}

	id := r.URL.Query().Get("id")
	idInt, err := strconv.Atoi(id)
	if err != nil {
		http.Error(w, "Invalid film ID", http.StatusBadRequest)
		return
	}
	err = h.service.DeleteMovie(r.Context(), idInt)
	if err != nil {
		http.Error(w, "Failed to delete film", http.StatusInternalServerError)
		return
	}
}

func (h *MovieHandler) GetAllFilms(w http.ResponseWriter, r *http.Request) {
	var listFilms []model.Film
	sortBy := r.URL.Query().Get("sort_by")

	if err := model.ValidateSortFilm(sortBy); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	listFilms, err := h.service.GetFilms(r.Context(), sortBy)
	if err != nil {
		http.Error(w, "Failed to get films", http.StatusInternalServerError)
		fmt.Println(err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(listFilms)
}

func (h *MovieHandler) SearchFilm(w http.ResponseWriter, r *http.Request) {
	var filmSearch model.Film

	actor, movie := strings.TrimSpace(r.URL.Query().Get("actor")), strings.TrimSpace(r.URL.Query().Get("movie"))

	if err := filmSearch.ValidateFilmSearchParams(movie, actor); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	films, err := h.service.SearchFilm(r.Context(), actor, movie)
	if err != nil {
		http.Error(w, fmt.Sprintf("Search failed: %v", err), http.StatusInternalServerError)
		return
	}

	// if (films == model.Film{}) {
	// 	w.WriteHeader(http.StatusNotFound)
	// 	return
	// }

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(films)
}
