package handlers

import (
	"encoding/json"
	"fluxera/internal/service"
	"net/http"
)

type AuthHandler struct {
	auth *service.AuthService
}

func NewAuthHandler(auth *service.AuthService) *AuthHandler {
	return &AuthHandler{auth: auth}
}

type registerRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}
type loginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	var req registerRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid json body")
		return
	}

	newUser, err := h.auth.Register(r.Context(), req.Email, req.Password)
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	writeJSON(w, http.StatusCreated, newUser)
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var login loginRequest

	if err := json.NewDecoder(r.Body).Decode(&login); err != nil {
		writeError(w, http.StatusBadRequest, "invalid json body")
		return
	}

	token, err := h.auth.Login(r.Context(), login.Email, login.Password)
	if err != nil {
		writeError(w, http.StatusUnauthorized, err.Error())
		return
	}

	writeJSON(w, http.StatusOK, map[string]string{"token": token})
}
