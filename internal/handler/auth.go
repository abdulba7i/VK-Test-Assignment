package handler

import (
	"encoding/json"
	"film-library/internal/model"
	"film-library/internal/service"
	"fmt"
	"net/http"
)

type AuthHandler struct {
	service service.AuthService
}

func NewAuthHandler(service service.AuthService) AuthHandler {
	return AuthHandler{service: service}
}

func (h *AuthHandler) HandleAuthPost(w http.ResponseWriter, r *http.Request) {
	switch {
	case r.URL.Path == "/sign_up" && r.Method == http.MethodPost:
		h.CreateUser(w, r)
	case r.URL.Path == "/sign_in" && r.Method == http.MethodPost:
		h.VerifyUser(w, r)
	default:
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	}
}
func (h *AuthHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	var user model.User

	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// ВАЖНО: проверяем обязательные поля
	if user.Username == "" || user.Password == "" || user.Role == 0 {
		http.Error(w, "username, password and role are required", http.StatusBadRequest)
		return
	}

	result, err := h.service.CreateUser(r.Context(), user)
	if err != nil {
		http.Error(w, fmt.Sprintf("failed to create user: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(result)
}

func (h *AuthHandler) VerifyUser(w http.ResponseWriter, r *http.Request) {
	var input model.SignInRequest
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	token, user, err := h.service.VerifyUser(r.Context(), input.Username, input.Password)

	if err != nil {
		http.Error(w, fmt.Sprintf("failed to verify user: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(model.AuthResponse{
		AccessToken: token,
		User: model.UserResponse{
			ID:       user.ID,
			Username: user.Username,
			Role:     user.Role,
		},
	})
}
