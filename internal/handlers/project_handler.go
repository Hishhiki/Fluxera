package handlers

import (
	"encoding/json"
	"fluxera/internal/middleware"
	"fluxera/internal/service"
	"net/http"
)

type ProjectHandler struct {
	project *service.ProjectService
}

func NewProjectHandler(project *service.ProjectService) *ProjectHandler {
	return &ProjectHandler{project: project}
}

type createProjectRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

func (h *ProjectHandler) CreateProject(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(middleware.UserIDKey).(int64)
	if !ok {
		writeError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	var req createProjectRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid json body")
		return
	}

	newProject, err := h.project.Create(r.Context(), userID, req.Name, req.Description)
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	writeJSON(w, http.StatusCreated, newProject)
}
