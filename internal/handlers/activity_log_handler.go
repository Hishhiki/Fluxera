package handlers

import (
	"fluxera/internal/middleware"
	"fluxera/internal/service"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

type ActivityLogHandler struct {
	logs *service.ActivityLogService
}

func NewActivityLogHandler(logs *service.ActivityLogService) *ActivityLogHandler {
	return &ActivityLogHandler{logs: logs}
}

func (h *ActivityLogHandler) GetProjectActivity(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(middleware.UserIDKey).(int64)
	if !ok {
		writeError(w, http.StatusUnauthorized, "unauthorized")
		return
	}
	projectID, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid project id")
		return
	}
	logs, err := h.logs.GetByProject(r.Context(), projectID, userID)
	if err != nil {
		writeError(w, http.StatusNotFound, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, logs)
}
