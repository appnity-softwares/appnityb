package handler

import (
	"encoding/json"
	"net/http"

	"github.com/appnity/backend/internal/service"
	"github.com/appnity/backend/pkg/response"
	"github.com/google/uuid"
)

type MediaHandler struct {
	mediaService *service.MediaService
}

func NewMediaHandler(mediaService *service.MediaService) *MediaHandler {
	return &MediaHandler{mediaService: mediaService}
}

func (h *MediaHandler) GetMedia(w http.ResponseWriter, r *http.Request) {
	page, pageSize := parsePagination(r)
	folder := r.URL.Query().Get("folder")

	media, total, err := h.mediaService.GetMedia(r.Context(), folder, page, pageSize)
	if err != nil {
		response.Error(w, "failed to get media", http.StatusInternalServerError)
		return
	}

	response.Paginated(w, "media retrieved", media, total, int64(page), int64(pageSize))
}

func (h *MediaHandler) UploadMedia(w http.ResponseWriter, r *http.Request) {
	file, header, err := r.FormFile("file")
	if err != nil {
		response.Error(w, "file is required", http.StatusBadRequest)
		return
	}
	defer file.Close()

	media, err := h.mediaService.UploadFile(r.Context(), file, header)
	if err != nil {
		response.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response.Success(w, "media uploaded", media, http.StatusCreated)
}

func (h *MediaHandler) DeleteMedia(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	if idStr == "" {
		response.Error(w, "media ID is required", http.StatusBadRequest)
		return
	}

	id, err := uuid.Parse(idStr)
	if err != nil {
		response.Error(w, "invalid media ID", http.StatusBadRequest)
		return
	}

	if err := h.mediaService.DeleteMedia(r.Context(), id); err != nil {
		response.Error(w, "failed to delete media", http.StatusInternalServerError)
		return
	}

	response.Success(w, "media deleted", nil, http.StatusOK)
}

func (h *MediaHandler) UpdateMedia(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	if idStr == "" {
		response.Error(w, "media ID is required", http.StatusBadRequest)
		return
	}

	id, err := uuid.Parse(idStr)
	if err != nil {
		response.Error(w, "invalid media ID", http.StatusBadRequest)
		return
	}

	var updates map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&updates); err != nil {
		response.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	media, err := h.mediaService.UpdateMedia(r.Context(), id, updates)
	if err != nil {
		response.Error(w, "failed to update media", http.StatusInternalServerError)
		return
	}

	response.Success(w, "media updated", media, http.StatusOK)
}
