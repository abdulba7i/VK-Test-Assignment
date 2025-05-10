package handler

import (
	"encoding/json"
	"film-library/internal/model"
	authmid "film-library/internal/utils/auth_mid"
	"film-library/internal/utils/response"
	"fmt"
	"strconv"
	"strings"

	"film-library/internal/service"
	"net/http"
)

type ActorHandler struct {
	service service.Actor
}

func NewActorHandler(service service.Actor) ActorHandler {
	return ActorHandler{service: service}
}

//

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
	default:
		response.WriteJSONError(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func (h *ActorHandler) HandleActorDelete(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodDelete:
		h.DeleteActor(w, r)
	default:
		response.WriteJSONError(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

// @Summary Create Actor
// @Security ApiKeyAuth
// @Tags actor
// @Description Create Actor
// @ID create-actor
// @Accept  json
// @Produce  json
// @Param actor body model.Actor true "Create Actor"
// @Success 201 {object} model.Actor
// @Failure 400,403 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Failure default {object} response.ErrorResponse
// @Router /actor_create [post]
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

// @Summary Update Actor
// @Security ApiKeyAuth
// @Tags actor
// @Description Update Actor
// @ID update-actor
// @Accept  json
// @Produce  json
// @Param actor body model.Actor true "Update Actor"
// @Success 201 {object} model.Actor
// @Failure 400,403 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Failure default {object} response.ErrorResponse
// @Router /actor_update [put]
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
		response.WriteJSONError(w, "Failed to update actor", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(actor)
}

// @Summary Delete Actor
// @Security ApiKeyAuth
// @Tags actor
// @Description Delete Actor
// @ID delete-actor
// @Accept  json
// @Produce  json
// @Param id query int true "Actor ID"
// @Success 200 {object} map[string]string
// @Failure 400,403 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Failure default {object} response.ErrorResponse
// @Router /actor_delete/{id} [delete]
func (h *ActorHandler) DeleteActor(w http.ResponseWriter, r *http.Request) {
	if authmid.IsAdmin(r) {
		response.WriteJSONError(w, "Forbidden: admin access required", http.StatusForbidden)
		return
	}

	parts := strings.Split(strings.Trim(r.URL.Path, "/"), "/")
	if len(parts) < 2 {
		response.WriteJSONError(w, "Missing actor ID", http.StatusBadRequest)
		return
	}

	idStr := parts[len(parts)-1]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		response.WriteJSONError(w, "Invalid actor ID", http.StatusBadRequest)
		return
	}

	err = h.service.DeleteActor(r.Context(), id)
	if err != nil {
		response.WriteJSONError(w, "Failed to delete actor", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	response := map[string]string{
		"message": "actor deleted successfully",
	}

	json.NewEncoder(w).Encode(response)
}
