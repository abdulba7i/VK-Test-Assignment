package handler

import (
	"encoding/json"
	"film-library/internal/model"
	authmid "film-library/internal/utils/auth_mid"
	"film-library/internal/utils/response"
	"fmt"
	"strconv"

	"film-library/internal/service"
	"net/http"
)

type ActorHandler struct {
	service service.Actor
}

func NewActorHandler(service service.Actor) ActorHandler {
	return ActorHandler{service: service}
}

func (h *ActorHandler) HandleActorPost(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		h.CreateActor(w, r)
	default:
		response.WriteJSONError(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func (h *ActorHandler) HandleActorPut(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPut:
		h.UpdateActor(w, r)
	case http.MethodDelete:
		h.DeleteActor(w, r)
	default:
		response.WriteJSONError(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func (h *ActorHandler) CreateActor(w http.ResponseWriter, r *http.Request) {
	if authmid.IsAdmin(r) {
		response.WriteJSONError(w, "Forbidden: admin access required", http.StatusForbidden)
		return
	}

	var actor model.Actor
	if err := json.NewDecoder(r.Body).Decode(&actor); err != nil {
		response.WriteJSONError(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if err := actor.Validate(); err != nil {
		response.WriteJSONError(w, fmt.Sprintf("%v", err.Error()), http.StatusBadRequest)
		return
	}

	err := h.service.AddActor(r.Context(), actor)
	if err != nil {
		response.WriteJSONError(w, "Failed to create actor", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(actor)
}

func (h *ActorHandler) UpdateActor(w http.ResponseWriter, r *http.Request) {
	if authmid.IsAdmin(r) {
		response.WriteJSONError(w, "Forbidden: admin access required", http.StatusForbidden)
		return
	}

	var actor model.Actor
	if err := json.NewDecoder(r.Body).Decode(&actor); err != nil {
		response.WriteJSONError(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if err := actor.Validate(); err != nil {
		response.WriteJSONError(w, fmt.Sprintf("%v", err.Error()), http.StatusBadRequest)
		return
	}

	err := h.service.UpdateActor(r.Context(), actor)
	if err != nil {
		response.WriteJSONError(w, fmt.Sprintf("Failed to update actor: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(actor)

}

func (h *ActorHandler) DeleteActor(w http.ResponseWriter, r *http.Request) {
	if authmid.IsAdmin(r) {
		response.WriteJSONError(w, "Forbidden: admin access required", http.StatusForbidden)
		return
	}

	id := r.URL.Query().Get("id")
	idInt, err := strconv.Atoi(id)

	if err != nil {
		response.WriteJSONError(w, "Invalid actor ID", http.StatusBadRequest)
		return
	}

	err = h.service.DeleteActor(r.Context(), idInt)
	if err != nil {
		response.WriteJSONError(w, "Failed to delete actor", http.StatusInternalServerError)
		return
	}
}
