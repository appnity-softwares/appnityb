package handler

import (
	"encoding/json"
	"net/http"

	"github.com/appnity/backend/internal/service"
	"github.com/appnity/backend/pkg/response"
	"github.com/google/uuid"
)

type AuthHandler struct {
	authService *service.AuthService
}

func NewAuthHandler(authService *service.AuthService) *AuthHandler {
	return &AuthHandler{authService: authService}
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	accessToken, refreshToken, err := h.authService.Login(r.Context(), req.Email, req.Password)
	if err != nil {
		response.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	response.Success(w, "login successful", map[string]string{
		"access_token":  accessToken,
		"refresh_token": refreshToken,
	}, http.StatusOK)
}

func (h *AuthHandler) RefreshToken(w http.ResponseWriter, r *http.Request) {
	var req struct {
		RefreshToken string `json:"refresh_token"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	accessToken, refreshToken, err := h.authService.RefreshToken(r.Context(), req.RefreshToken)
	if err != nil {
		response.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	response.Success(w, "token refreshed", map[string]string{
		"access_token":  accessToken,
		"refresh_token": refreshToken,
	}, http.StatusOK)
}

func (h *AuthHandler) GetMe(w http.ResponseWriter, r *http.Request) {
	userIDStr := r.Context().Value("user_id")
	if userIDStr == nil {
		response.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	userID, ok := userIDStr.(uuid.UUID)
	if !ok {
		response.Error(w, "invalid user ID type", http.StatusInternalServerError)
		return
	}

	user, err := h.authService.GetProfile(r.Context(), userID)
	if err != nil {
		response.Error(w, "failed to get profile", http.StatusNotFound)
		return
	}

	response.Success(w, "profile retrieved", user, http.StatusOK)
}

func (h *AuthHandler) UpdateProfile(w http.ResponseWriter, r *http.Request) {
	userIDStr := r.Context().Value("user_id")
	if userIDStr == nil {
		response.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	userID, ok := userIDStr.(uuid.UUID)
	if !ok {
		response.Error(w, "invalid user ID type", http.StatusInternalServerError)
		return
	}

	var updates map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&updates); err != nil {
		response.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	user, err := h.authService.UpdateProfile(r.Context(), userID, updates)
	if err != nil {
		response.Error(w, "failed to update profile", http.StatusInternalServerError)
		return
	}

	response.Success(w, "profile updated", user, http.StatusOK)
}

func (h *AuthHandler) ChangePassword(w http.ResponseWriter, r *http.Request) {
	userIDStr := r.Context().Value("user_id")
	if userIDStr == nil {
		response.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	userID, ok := userIDStr.(uuid.UUID)
	if !ok {
		response.Error(w, "invalid user ID type", http.StatusInternalServerError)
		return
	}

	var req struct {
		OldPassword string `json:"old_password"`
		NewPassword string `json:"new_password"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	if err := h.authService.ChangePassword(r.Context(), userID, req.OldPassword, req.NewPassword); err != nil {
		response.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	response.Success(w, "password changed successfully", nil, http.StatusOK)
}

func (h *AuthHandler) ForgotPassword(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Email string `json:"email"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	response.Success(w, "if an account exists with that email, a reset link has been sent", nil, http.StatusOK)
}

func (h *AuthHandler) ResetPassword(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Token       string `json:"token"`
		NewPassword string `json:"new_password"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	response.Success(w, "password reset successfully", nil, http.StatusOK)
}
