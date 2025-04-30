package handler

import (
	"encoding/json"
	"film-library/internal/model"
	"film-library/internal/service"
	"film-library/internal/utils/response"
	"fmt"
	"net/http"
)

type AuthHandler struct {
	service service.Authorization
}

func NewAuthHandler(service service.Authorization) AuthHandler {
	return AuthHandler{service: service}
}

func (h *AuthHandler) HandleAuthPost(w http.ResponseWriter, r *http.Request) {
	switch {
	case r.URL.Path == "/sign_up" && r.Method == http.MethodPost:
		h.CreateUser(w, r)
	case r.URL.Path == "/sign_in" && r.Method == http.MethodPost:
		h.VerifyUser(w, r)
	default:
		response.WriteJSONError(w, "method not allowed", http.StatusMethodNotAllowed)
	}
}

func (h *AuthHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	var req model.SignUpRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.WriteJSONError(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	user := model.User{
		Username: req.Username,
		Password: req.Password,
		Role:     req.Role,
	}

	if user.Username == "" || user.Password == "" || user.Role == 0 {
		response.WriteJSONError(w, "username, password and role are required", http.StatusBadRequest)
		return
	}

	result, err := h.service.CreateUser(r.Context(), user)
	if err != nil {
		response.WriteJSONError(w, fmt.Sprintf("failed to create user: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"token": result})
}

func (h *AuthHandler) VerifyUser(w http.ResponseWriter, r *http.Request) {
	var input model.SignInRequest
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		response.WriteJSONError(w, fmt.Sprintf("Invalid request body: %v", err), http.StatusBadRequest)
		return
	}

	token, user, err := h.service.VerifyUser(r.Context(), input.Username, input.Password)

	if err != nil {
		response.WriteJSONError(w, fmt.Sprintf("failed to verify user: %v", err), http.StatusInternalServerError)
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

// func writeJSONError(w http.ResponseWriter, message string, status int) {
// 	w.Header().Set("Content-Type", "application/json")
// 	w.WriteHeader(status)
// 	json.NewEncoder(w).Encode(map[string]string{"error": message})
// }
