package handlers

import (
	"encoding/json"
	"fluxera/internal/middleware"
	"fluxera/internal/service"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

type CommentHandler struct {
	comments *service.CommentService
}

func NewCommentHandler(comments *service.CommentService) *CommentHandler {
	return &CommentHandler{comments: comments}
}

type createCommentRequest struct {
	Content string `json:"content"`
}

type updateCommentRequest struct {
	Content string `json:"content"`
}

func (h *CommentHandler) CreateComment(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(middleware.UserIDKey).(int64)
	if !ok {
		writeError(w, http.StatusUnauthorized, "unauthorized")
		return
	}
	var req createCommentRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid json body")
		return
	}

	taskID, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid task id")
		return
	}

	comment, err := h.comments.Create(r.Context(), taskID, userID, req.Content)
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	writeJSON(w, http.StatusCreated, comment)
}

func (h *CommentHandler) GetCommentsByTask(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(middleware.UserIDKey).(int64)
	if !ok {
		writeError(w, http.StatusUnauthorized, "unauthorized")
		return
	}
	taskID, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid task id")
		return
	}

	comments, err := h.comments.GetByTask(r.Context(), taskID, userID)
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	writeJSON(w, http.StatusOK, comments)
}

func (h *CommentHandler) UpdateComment(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(middleware.UserIDKey).(int64)
	if !ok {
		writeError(w, http.StatusUnauthorized, "unauthorized")
		return
	}
	commentID, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid comment id")
		return
	}

	var req updateCommentRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid json body")
		return
	}

	updateComment, err := h.comments.Update(r.Context(), commentID, userID, req.Content)
	if err != nil {
		writeError(w, http.StatusForbidden, err.Error())
		return
	}

	writeJSON(w, http.StatusOK, updateComment)
}

func (h *CommentHandler) DeleteComment(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(middleware.UserIDKey).(int64)
	if !ok {
		writeError(w, http.StatusUnauthorized, "unauthorized")
		return
	}
	commentID, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid comment id")
		return
	}

	if err := h.comments.Delete(r.Context(), commentID, userID); err != nil {
		writeError(w, http.StatusNotFound, "comment not found")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
