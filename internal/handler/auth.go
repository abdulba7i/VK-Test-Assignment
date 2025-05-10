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

//	func (h *AuthHandler) HandleAuthPost(w http.ResponseWriter, r *http.Request) {
//		switch {
//		case r.URL.Path == "/auth/sign_up" && r.Method == http.MethodPost:
//			h.CreateUser(w, r)
//		case r.URL.Path == "/auth/sign_in" && r.Method == http.MethodPost:
//			h.VerifyUser(w, r)
//		default:
//			response.WriteJSONError(w, "method not allowed", http.StatusMethodNotAllowed)
//		}
//	}
func (h *AuthHandler) HandleAuthPost(w http.ResponseWriter, r *http.Request) {
	switch r.URL.Path {
	case "/auth/sign_up":
		h.CreateUser(w, r)
	case "/auth/sign_in":
		h.VerifyUser(w, r)
	default:
		http.Error(w, "Not found", http.StatusNotFound)
	}
}

// @Summary SignUp
// @Tags auth
// @Description Create a new user account and return JWT token
// @ID create-account
// @Accept  json
// @Produce  json
// @Param req body model.SignUpRequest true "Account info"
// @Success 201 {object} map[string]string
// @Failure 400,405 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Failure default {object} response.ErrorResponse
// @Router /auth/sign_up [post]
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
		response.WriteJSONError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"token": result})
}

// @Summary SignIn
// @Tags auth
// @Description Authenticate user and return JWT token + user info
// @ID login
// @Accept  json
// @Produce  json
// @Param input body model.SignInRequest true "Account info"
// @Success 200 {object} model.AuthResponse
// @Failure 400,405 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Failure default {object} response.ErrorResponse
// @Router /auth/sign_in [post]
func (h *AuthHandler) VerifyUser(w http.ResponseWriter, r *http.Request) {
	var input model.SignInRequest
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		response.WriteJSONError(w, fmt.Sprintf("Invalid request body: %v", err), http.StatusBadRequest)
		return
	}

	if input.Username == "" || input.Password == "" {
		response.WriteJSONError(w, "username or password are required", http.StatusBadRequest)
		return
	}

	token, user, err := h.service.VerifyUser(r.Context(), input.Username, input.Password)

	if err != nil {
		response.WriteJSONError(w, "failed to verify user", http.StatusInternalServerError)
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
