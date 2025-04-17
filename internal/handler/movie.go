package handler

import (
	"encoding/json"
	"film-library/internal/model"
	"film-library/internal/service"
	"net/http"
	"strconv"
)

type MovieHandler struct {
	service service.MovieService
}

func NewMovieHandler(service service.MovieService) MovieHandler {
	return MovieHandler{service: service}
}

func (h *MovieHandler) HandleMovies(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		h.CreateFilm(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func (h *MovieHandler) HandleMovie(w http.ResponseWriter, r *http.Request) {
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
	// var listFilms []model.Film
	// sortBy := r.URL.Query().Get("sort_by")

	// if err := listFilms.ValidateSortFilm(sortBy); err != nil {
	// 	http.Error(w, err.Error(), http.StatusBadRequest)
	// 	return
	// }

}
