package handler

import (
	"encoding/json"
	"net/http"

	"github.com/appnity/backend/internal/service"
	"github.com/appnity/backend/pkg/response"
	"github.com/google/uuid"
)

type UserHandler struct {
	authService *service.AuthService
	userService *service.UserService
}

func NewUserHandler(authService *service.AuthService, userService *service.UserService) *UserHandler {
	return &UserHandler{
		authService: authService,
		userService: userService,
	}
}

func (h *UserHandler) GetAllUsers(w http.ResponseWriter, r *http.Request) {
	users, err := h.userService.GetAllUsers(r.Context())
	if err != nil {
		response.Error(w, "failed to get users", http.StatusInternalServerError)
		return
	}

	response.Success(w, "users retrieved", users, http.StatusOK)
}

func (h *UserHandler) GetUser(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		response.Error(w, "invalid user ID", http.StatusBadRequest)
		return
	}

	user, err := h.userService.GetUserByID(r.Context(), id)
	if err != nil {
		response.Error(w, "user not found", http.StatusNotFound)
		return
	}

	response.Success(w, "user retrieved", user, http.StatusOK)
}

func (h *UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Email    string `json:"email"`
		Password string `json:"password"`
		FullName string `json:"full_name"`
		Role     string `json:"role"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	if req.Email == "" || req.Password == "" {
		response.Error(w, "email and password are required", http.StatusBadRequest)
		return
	}

	if req.Role == "" {
		req.Role = "admin"
	}

	user, err := h.userService.CreateUser(r.Context(), req.Email, req.Password, req.FullName, req.Role)
	if err != nil {
		response.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response.Success(w, "user created", user, http.StatusCreated)
}

func (h *UserHandler) UpdateUser(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		response.Error(w, "invalid user ID", http.StatusBadRequest)
		return
	}

	var updates map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&updates); err != nil {
		response.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	user, err := h.userService.UpdateUser(r.Context(), id, updates)
	if err != nil {
		response.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response.Success(w, "user updated", user, http.StatusOK)
}

func (h *UserHandler) DeleteUser(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		response.Error(w, "invalid user ID", http.StatusBadRequest)
		return
	}

	if err := h.userService.DeleteUser(r.Context(), id); err != nil {
		response.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response.Success(w, "user deleted", nil, http.StatusOK)
}

func (h *UserHandler) UpdateUserRole(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		response.Error(w, "invalid user ID", http.StatusBadRequest)
		return
	}

	var req struct {
		Role string `json:"role"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	if req.Role != "super_admin" && req.Role != "admin" && req.Role != "editor" {
		response.Error(w, "invalid role", http.StatusBadRequest)
		return
	}

	user, err := h.userService.UpdateUserRole(r.Context(), id, req.Role)
	if err != nil {
		response.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response.Success(w, "user role updated", user, http.StatusOK)
}
