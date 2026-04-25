package handler

import (
	"encoding/json"
	"net/http"

	"github.com/appnity/backend/internal/models"
	"github.com/appnity/backend/internal/service"
	"github.com/appnity/backend/pkg/response"
	"github.com/google/uuid"
)

type AdminHandler struct {
	contentService *service.ContentService
	themeService   *service.ThemeService
	authService    *service.AuthService
	mediaService   *service.MediaService
	configService  *service.ConfigService
}

func NewAdminHandler(
	contentService *service.ContentService,
	themeService *service.ThemeService,
	authService *service.AuthService,
	mediaService *service.MediaService,
	configService *service.ConfigService,
) *AdminHandler {
	return &AdminHandler{
		contentService: contentService,
		themeService:   themeService,
		authService:    authService,
		mediaService:   mediaService,
		configService:  configService,
	}
}

func (h *AdminHandler) GetDashboardStats(w http.ResponseWriter, r *http.Request) {
	stats, err := h.configService.GetDashboardStats(r.Context())
	if err != nil {
		response.Error(w, "failed to get dashboard stats", http.StatusInternalServerError)
		return
	}

	response.Success(w, "dashboard stats retrieved", stats, http.StatusOK)
}

func (h *AdminHandler) GetRecentContacts(w http.ResponseWriter, r *http.Request) {
	stats, err := h.configService.GetDashboardStats(r.Context())
	if err != nil {
		response.Error(w, "failed to get recent contacts", http.StatusInternalServerError)
		return
	}

	response.Success(w, "recent contacts retrieved", stats["recent_contacts"], http.StatusOK)
}

func (h *AdminHandler) GetTheme(w http.ResponseWriter, r *http.Request) {
	theme, err := h.themeService.GetActiveTheme(r.Context())
	if err != nil {
		theme = h.getDefaultTheme()
	}

	response.Success(w, "theme retrieved", theme, http.StatusOK)
}

func (h *AdminHandler) getDefaultTheme() *models.Theme {
	defaultColors := map[string]string{
		"primary":       "#ff6b00",
		"secondary":     "#1a1a1a",
		"accent":        "#ffffff",
		"background":    "#0a0a0a",
		"text":          "#ffffff",
		"surface":       "#1a1a1a",
		"border":        "#2a2a2a",
		"gradientStart": "#ff6b00",
		"gradientEnd":   "#ff8533",
		"error":         "#ef4444",
		"success":       "#22c55e",
		"warning":       "#f59e0b",
		"muted":         "#6b7280",
		"textSecondary": "#a0a0a0",
	}
	defaultFonts := map[string]string{
		"heading": "Poppins",
		"body":    "Inter",
	}

	colorsBytes, _ := json.Marshal(defaultColors)
	fontsBytes, _ := json.Marshal(defaultFonts)

	return &models.Theme{
		Name:     "Dark Orange",
		IsActive: true,
		Colors:   colorsBytes,
		Fonts:    fontsBytes,
	}
}

func (h *AdminHandler) UpdateTheme(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Name   string `json:"name"`
		Colors struct {
			Primary       string `json:"primary"`
			Secondary     string `json:"secondary"`
			Accent        string `json:"accent"`
			Background    string `json:"background"`
			Text          string `json:"text"`
			Surface       string `json:"surface"`
			Border        string `json:"border"`
			GradientStart string `json:"gradientStart"`
			GradientEnd   string `json:"gradientEnd"`
			Error         string `json:"error"`
			Success       string `json:"success"`
			Warning       string `json:"warning"`
			Muted         string `json:"muted"`
			TextSecondary string `json:"textSecondary"`
		} `json:"colors"`
		Fonts struct {
			Heading string `json:"heading"`
			Body    string `json:"body"`
		} `json:"fonts"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	theme, err := h.themeService.GetActiveTheme(r.Context())
	if err != nil {
		theme = &models.Theme{}
	}

	theme.Name = req.Name
	theme.IsActive = true
	colorsBytes, _ := json.Marshal(req.Colors)
	fontsBytes, _ := json.Marshal(req.Fonts)
	theme.Colors = colorsBytes
	theme.Fonts = fontsBytes

	if err := h.themeService.UpdateTheme(r.Context(), theme); err != nil {
		response.Error(w, "failed to update theme", http.StatusInternalServerError)
		return
	}

	response.Success(w, "theme updated", theme, http.StatusOK)
}

func (h *AdminHandler) PreviewTheme(w http.ResponseWriter, r *http.Request) {
	var theme models.Theme
	if err := json.NewDecoder(r.Body).Decode(&theme); err != nil {
		response.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	response.Success(w, "theme preview generated", theme, http.StatusOK)
}

func (h *AdminHandler) ResetTheme(w http.ResponseWriter, r *http.Request) {
	theme, err := h.themeService.ResetToDefaults(r.Context())
	if err != nil {
		response.Error(w, "failed to reset theme", http.StatusInternalServerError)
		return
	}

	response.Success(w, "theme reset to defaults", theme, http.StatusOK)
}

func (h *AdminHandler) GetThemePresets(w http.ResponseWriter, r *http.Request) {
	presets := []map[string]interface{}{
		{
			"name": "Dark Orange",
			"colors": map[string]string{
				"primary":       "#ff6b00",
				"secondary":     "#1a1a1a",
				"accent":        "#ffffff",
				"background":    "#0a0a0a",
				"text":          "#ffffff",
				"surface":       "#1a1a1a",
				"border":        "#2a2a2a",
				"gradientStart": "#ff6b00",
				"gradientEnd":   "#ff8533",
				"error":         "#ef4444",
				"success":       "#22c55e",
				"warning":       "#f59e0b",
				"muted":         "#6b7280",
				"textSecondary": "#a0a0a0",
			},
			"fonts": map[string]string{
				"heading": "Poppins",
				"body":    "Inter",
			},
		},
		{
			"name": "Light Minimal",
			"colors": map[string]string{
				"primary":       "#000000",
				"secondary":     "#f5f5f5",
				"accent":        "#333333",
				"background":    "#ffffff",
				"text":          "#1a1a1a",
				"surface":       "#f9fafb",
				"border":        "#e5e7eb",
				"gradientStart": "#000000",
				"gradientEnd":   "#374151",
				"error":         "#ef4444",
				"success":       "#22c55e",
				"warning":       "#f59e0b",
				"muted":         "#6b7280",
				"textSecondary": "#6b7280",
			},
			"fonts": map[string]string{
				"heading": "Inter",
				"body":    "Inter",
			},
		},
		{
			"name": "Corporate Blue",
			"colors": map[string]string{
				"primary":       "#0066ff",
				"secondary":     "#1e3a5f",
				"accent":        "#ffffff",
				"background":    "#f8fafc",
				"text":          "#1e293b",
				"surface":       "#ffffff",
				"border":        "#e2e8f0",
				"gradientStart": "#0066ff",
				"gradientEnd":   "#0052cc",
				"error":         "#ef4444",
				"success":       "#22c55e",
				"warning":       "#f59e0b",
				"muted":         "#64748b",
				"textSecondary": "#64748b",
			},
			"fonts": map[string]string{
				"heading": "Poppins",
				"body":    "Open Sans",
			},
		},
		{
			"name": "Nature Green",
			"colors": map[string]string{
				"primary":       "#16a34a",
				"secondary":     "#14532d",
				"accent":        "#fef9c3",
				"background":    "#f0fdf4",
				"text":          "#14532d",
				"surface":       "#ffffff",
				"border":        "#bbf7d0",
				"gradientStart": "#16a34a",
				"gradientEnd":   "#15803d",
				"error":         "#ef4444",
				"success":       "#22c55e",
				"warning":       "#f59e0b",
				"muted":         "#6b7280",
				"textSecondary": "#4b5563",
			},
			"fonts": map[string]string{
				"heading": "Montserrat",
				"body":    "Nunito",
			},
		},
		{
			"name": "Royal Purple",
			"colors": map[string]string{
				"primary":       "#7c3aed",
				"secondary":     "#1e1b4b",
				"accent":        "#fef3c7",
				"background":    "#0f0f23",
				"text":          "#e2e8f0",
				"surface":       "#1a1a2e",
				"border":        "#312e81",
				"gradientStart": "#7c3aed",
				"gradientEnd":   "#5b21b6",
				"error":         "#ef4444",
				"success":       "#22c55e",
				"warning":       "#f59e0b",
				"muted":         "#9ca3af",
				"textSecondary": "#9ca3af",
			},
			"fonts": map[string]string{
				"heading": "Playfair Display",
				"body":    "Inter",
			},
		},
		{
			"name": "Sunset Red",
			"colors": map[string]string{
				"primary":       "#dc2626",
				"secondary":     "#450a0a",
				"accent":        "#fef2f2",
				"background":    "#fff5f5",
				"text":          "#1f1f1f",
				"surface":       "#ffffff",
				"border":        "#fecaca",
				"gradientStart": "#dc2626",
				"gradientEnd":   "#b91c1c",
				"error":         "#dc2626",
				"success":       "#22c55e",
				"warning":       "#f59e0b",
				"muted":         "#6b7280",
				"textSecondary": "#6b7280",
			},
			"fonts": map[string]string{
				"heading": "Oswald",
				"body":    "Roboto",
			},
		},
	}
	response.Success(w, "presets retrieved", presets, http.StatusOK)
}

type ApplyPresetRequest struct {
	PresetName string `json:"presetName"`
}

func (h *AdminHandler) ApplyThemePreset(w http.ResponseWriter, r *http.Request) {
	var req ApplyPresetRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	presets := map[string]map[string]interface{}{
		"Dark Orange": {
			"name": "Dark Orange",
			"colors": map[string]string{
				"primary": "#ff6b00", "secondary": "#1a1a1a", "accent": "#ffffff",
				"background": "#0a0a0a", "text": "#ffffff", "surface": "#1a1a1a",
				"border": "#2a2a2a", "gradientStart": "#ff6b00", "gradientEnd": "#ff8533",
			},
			"fonts": map[string]string{"heading": "Poppins", "body": "Inter"},
		},
		"Light Minimal": {
			"name": "Light Minimal",
			"colors": map[string]string{
				"primary": "#000000", "secondary": "#f5f5f5", "accent": "#333333",
				"background": "#ffffff", "text": "#1a1a1a", "surface": "#f9fafb",
				"border": "#e5e7eb", "gradientStart": "#000000", "gradientEnd": "#374151",
			},
			"fonts": map[string]string{"heading": "Inter", "body": "Inter"},
		},
		"Corporate Blue": {
			"name": "Corporate Blue",
			"colors": map[string]string{
				"primary": "#0066ff", "secondary": "#1e3a5f", "accent": "#ffffff",
				"background": "#f8fafc", "text": "#1e293b", "surface": "#ffffff",
				"border": "#e2e8f0", "gradientStart": "#0066ff", "gradientEnd": "#0052cc",
			},
			"fonts": map[string]string{"heading": "Poppins", "body": "Open Sans"},
		},
		"Nature Green": {
			"name": "Nature Green",
			"colors": map[string]string{
				"primary": "#16a34a", "secondary": "#14532d", "accent": "#fef9c3",
				"background": "#f0fdf4", "text": "#14532d", "surface": "#ffffff",
				"border": "#bbf7d0", "gradientStart": "#16a34a", "gradientEnd": "#15803d",
			},
			"fonts": map[string]string{"heading": "Montserrat", "body": "Nunito"},
		},
		"Royal Purple": {
			"name": "Royal Purple",
			"colors": map[string]string{
				"primary": "#7c3aed", "secondary": "#1e1b4b", "accent": "#fef3c7",
				"background": "#0f0f23", "text": "#e2e8f0", "surface": "#1a1a2e",
				"border": "#312e81", "gradientStart": "#7c3aed", "gradientEnd": "#5b21b6",
			},
			"fonts": map[string]string{"heading": "Playfair Display", "body": "Inter"},
		},
		"Sunset Red": {
			"name": "Sunset Red",
			"colors": map[string]string{
				"primary": "#dc2626", "secondary": "#450a0a", "accent": "#fef2f2",
				"background": "#fff5f5", "text": "#1f1f1f", "surface": "#ffffff",
				"border": "#fecaca", "gradientStart": "#dc2626", "gradientEnd": "#b91c1c",
			},
			"fonts": map[string]string{"heading": "Oswald", "body": "Roboto"},
		},
	}

	preset, ok := presets[req.PresetName]
	if !ok {
		response.Error(w, "preset not found", http.StatusNotFound)
		return
	}

	theme, err := h.themeService.GetActiveTheme(r.Context())
	if err != nil {
		theme = &models.Theme{
			ID:       uuid.New(),
			IsActive: true,
		}
	}

	theme.Name = preset["name"].(string)
	colorsBytes, _ := json.Marshal(preset["colors"])
	fontsBytes, _ := json.Marshal(preset["fonts"])
	theme.Colors = colorsBytes
	theme.Fonts = fontsBytes
	theme.IsActive = true

	if err := h.themeService.UpdateTheme(r.Context(), theme); err != nil {
		response.Error(w, "failed to apply preset", http.StatusInternalServerError)
		return
	}

	response.Success(w, "preset applied successfully", theme, http.StatusOK)
}

func (h *AdminHandler) GetAllPages(w http.ResponseWriter, r *http.Request) {
	pages, err := h.contentService.GetAllPages(r.Context())
	if err != nil {
		response.Error(w, "failed to get pages", http.StatusInternalServerError)
		return
	}

	response.Success(w, "pages retrieved", pages, http.StatusOK)
}

func (h *AdminHandler) GetPage(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	if idStr == "" {
		response.Error(w, "page ID is required", http.StatusBadRequest)
		return
	}

	id, err := uuid.Parse(idStr)
	if err != nil {
		response.Error(w, "invalid page ID", http.StatusBadRequest)
		return
	}

	page, err := h.contentService.GetPageByID(r.Context(), id)
	if err != nil {
		response.Error(w, "page not found", http.StatusNotFound)
		return
	}

	response.Success(w, "page retrieved", page, http.StatusOK)
}

func (h *AdminHandler) UpdatePage(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	if idStr == "" {
		response.Error(w, "page ID is required", http.StatusBadRequest)
		return
	}

	id, err := uuid.Parse(idStr)
	if err != nil {
		response.Error(w, "invalid page ID", http.StatusBadRequest)
		return
	}

	var page models.Page
	if err := json.NewDecoder(r.Body).Decode(&page); err != nil {
		response.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	page.ID = id

	if err := h.contentService.UpdatePage(r.Context(), &page); err != nil {
		response.Error(w, "failed to update page", http.StatusInternalServerError)
		return
	}

	response.Success(w, "page updated", page, http.StatusOK)
}

func (h *AdminHandler) UpdatePageSections(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	if idStr == "" {
		response.Error(w, "page ID is required", http.StatusBadRequest)
		return
	}

	id, err := uuid.Parse(idStr)
	if err != nil {
		response.Error(w, "invalid page ID", http.StatusBadRequest)
		return
	}

	var req struct {
		Sections json.RawMessage `json:"sections"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	if err := h.contentService.UpdatePageSections(r.Context(), id, req.Sections); err != nil {
		response.Error(w, "failed to update page sections", http.StatusInternalServerError)
		return
	}

	response.Success(w, "page sections updated", nil, http.StatusOK)
}

func (h *AdminHandler) GetAllComponents(w http.ResponseWriter, r *http.Request) {
	components, err := h.contentService.GetAllComponents(r.Context())
	if err != nil {
		response.Error(w, "failed to get components", http.StatusInternalServerError)
		return
	}

	response.Success(w, "components retrieved", components, http.StatusOK)
}

func (h *AdminHandler) GetComponent(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	if idStr == "" {
		response.Error(w, "component ID is required", http.StatusBadRequest)
		return
	}

	id, err := uuid.Parse(idStr)
	if err != nil {
		response.Error(w, "invalid component ID", http.StatusBadRequest)
		return
	}

	component, err := h.contentService.GetComponentByID(r.Context(), id)
	if err != nil {
		response.Error(w, "component not found", http.StatusNotFound)
		return
	}

	response.Success(w, "component retrieved", component, http.StatusOK)
}

func (h *AdminHandler) UpdateComponent(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	if idStr == "" {
		response.Error(w, "component ID is required", http.StatusBadRequest)
		return
	}

	id, err := uuid.Parse(idStr)
	if err != nil {
		response.Error(w, "invalid component ID", http.StatusBadRequest)
		return
	}

	var component models.Component
	if err := json.NewDecoder(r.Body).Decode(&component); err != nil {
		response.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	component.ID = id

	if err := h.contentService.UpdateComponent(r.Context(), &component); err != nil {
		response.Error(w, "failed to update component", http.StatusInternalServerError)
		return
	}

	response.Success(w, "component updated", component, http.StatusOK)
}

func (h *AdminHandler) ToggleComponentVisibility(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	if idStr == "" {
		response.Error(w, "component ID is required", http.StatusBadRequest)
		return
	}

	id, err := uuid.Parse(idStr)
	if err != nil {
		response.Error(w, "invalid component ID", http.StatusBadRequest)
		return
	}

	if err := h.contentService.ToggleComponentVisibility(r.Context(), id); err != nil {
		response.Error(w, "failed to toggle component visibility", http.StatusInternalServerError)
		return
	}

	response.Success(w, "component visibility toggled", nil, http.StatusOK)
}

func (h *AdminHandler) ReorderComponents(w http.ResponseWriter, r *http.Request) {
	pageIDStr := r.PathValue("page_id")
	if pageIDStr == "" {
		response.Error(w, "page ID is required", http.StatusBadRequest)
		return
	}

	pageID, err := uuid.Parse(pageIDStr)
	if err != nil {
		response.Error(w, "invalid page ID", http.StatusBadRequest)
		return
	}

	var req struct {
		IDs []uuid.UUID `json:"ids"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	if err := h.contentService.ReorderComponents(r.Context(), pageID, req.IDs); err != nil {
		response.Error(w, "failed to reorder components", http.StatusInternalServerError)
		return
	}

	response.Success(w, "components reordered", nil, http.StatusOK)
}

func (h *AdminHandler) GetAllNav(w http.ResponseWriter, r *http.Request) {
	nav, err := h.contentService.GetAllNavigation(r.Context(), false)
	if err != nil {
		response.Error(w, "failed to get navigation", http.StatusInternalServerError)
		return
	}

	response.Success(w, "navigation retrieved", nav, http.StatusOK)
}

func (h *AdminHandler) CreateNav(w http.ResponseWriter, r *http.Request) {
	var nav models.Navigation
	if err := json.NewDecoder(r.Body).Decode(&nav); err != nil {
		response.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	if err := h.contentService.CreateNavigation(r.Context(), &nav); err != nil {
		response.Error(w, "failed to create navigation", http.StatusInternalServerError)
		return
	}

	response.Success(w, "navigation created", nav, http.StatusCreated)
}

func (h *AdminHandler) UpdateNav(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	if idStr == "" {
		response.Error(w, "navigation ID is required", http.StatusBadRequest)
		return
	}

	id, err := uuid.Parse(idStr)
	if err != nil {
		response.Error(w, "invalid navigation ID", http.StatusBadRequest)
		return
	}

	var nav models.Navigation
	if err := json.NewDecoder(r.Body).Decode(&nav); err != nil {
		response.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	nav.ID = id

	if err := h.contentService.UpdateNavigation(r.Context(), &nav); err != nil {
		response.Error(w, "failed to update navigation", http.StatusInternalServerError)
		return
	}

	response.Success(w, "navigation updated", nav, http.StatusOK)
}

func (h *AdminHandler) DeleteNav(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	if idStr == "" {
		response.Error(w, "navigation ID is required", http.StatusBadRequest)
		return
	}

	id, err := uuid.Parse(idStr)
	if err != nil {
		response.Error(w, "invalid navigation ID", http.StatusBadRequest)
		return
	}

	if err := h.contentService.DeleteNavigation(r.Context(), id); err != nil {
		response.Error(w, "failed to delete navigation", http.StatusInternalServerError)
		return
	}

	response.Success(w, "navigation deleted", nil, http.StatusOK)
}

func (h *AdminHandler) ReorderNav(w http.ResponseWriter, r *http.Request) {
	var req struct {
		IDs []uuid.UUID `json:"ids"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	if err := h.contentService.ReorderNavigation(r.Context(), req.IDs); err != nil {
		response.Error(w, "failed to reorder navigation", http.StatusInternalServerError)
		return
	}

	response.Success(w, "navigation reordered", nil, http.StatusOK)
}

func (h *AdminHandler) GetAllBlogs(w http.ResponseWriter, r *http.Request) {
	page, pageSize := parsePagination(r)
	publishedOnly := r.URL.Query().Get("published_only") == "true"

	blogs, total, err := h.contentService.GetAllBlogs(r.Context(), publishedOnly, page, pageSize)
	if err != nil {
		response.Error(w, "failed to get blogs", http.StatusInternalServerError)
		return
	}

	response.Paginated(w, "blogs retrieved", blogs, total, int64(page), int64(pageSize))
}

func (h *AdminHandler) GetBlog(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	if idStr == "" {
		response.Error(w, "blog ID is required", http.StatusBadRequest)
		return
	}

	id, err := uuid.Parse(idStr)
	if err != nil {
		response.Error(w, "invalid blog ID", http.StatusBadRequest)
		return
	}

	blog, err := h.contentService.GetBlogByID(r.Context(), id)
	if err != nil {
		response.Error(w, "blog not found", http.StatusNotFound)
		return
	}

	response.Success(w, "blog retrieved", blog, http.StatusOK)
}

func (h *AdminHandler) CreateBlog(w http.ResponseWriter, r *http.Request) {
	var blog models.Blog
	if err := json.NewDecoder(r.Body).Decode(&blog); err != nil {
		response.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	if err := h.contentService.CreateBlog(r.Context(), &blog); err != nil {
		response.Error(w, "failed to create blog", http.StatusInternalServerError)
		return
	}

	response.Success(w, "blog created", blog, http.StatusCreated)
}

func (h *AdminHandler) UpdateBlog(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	if idStr == "" {
		response.Error(w, "blog ID is required", http.StatusBadRequest)
		return
	}

	id, err := uuid.Parse(idStr)
	if err != nil {
		response.Error(w, "invalid blog ID", http.StatusBadRequest)
		return
	}

	var blog models.Blog
	if err := json.NewDecoder(r.Body).Decode(&blog); err != nil {
		response.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	blog.ID = id

	if err := h.contentService.UpdateBlog(r.Context(), &blog); err != nil {
		response.Error(w, "failed to update blog", http.StatusInternalServerError)
		return
	}

	response.Success(w, "blog updated", blog, http.StatusOK)
}

func (h *AdminHandler) DeleteBlog(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	if idStr == "" {
		response.Error(w, "blog ID is required", http.StatusBadRequest)
		return
	}

	id, err := uuid.Parse(idStr)
	if err != nil {
		response.Error(w, "invalid blog ID", http.StatusBadRequest)
		return
	}

	if err := h.contentService.DeleteBlog(r.Context(), id); err != nil {
		response.Error(w, "failed to delete blog", http.StatusInternalServerError)
		return
	}

	response.Success(w, "blog deleted", nil, http.StatusOK)
}

func (h *AdminHandler) GetAllPortfolios(w http.ResponseWriter, r *http.Request) {
	page, pageSize := parsePagination(r)
	visibleOnly := r.URL.Query().Get("visible_only") == "true"
	category := r.URL.Query().Get("category")

	portfolios, total, err := h.contentService.GetAllPortfolios(r.Context(), visibleOnly, category, page, pageSize)
	if err != nil {
		response.Error(w, "failed to get portfolios", http.StatusInternalServerError)
		return
	}

	response.Paginated(w, "portfolios retrieved", portfolios, total, int64(page), int64(pageSize))
}

func (h *AdminHandler) GetPortfolio(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	if idStr == "" {
		response.Error(w, "portfolio ID is required", http.StatusBadRequest)
		return
	}

	id, err := uuid.Parse(idStr)
	if err != nil {
		response.Error(w, "invalid portfolio ID", http.StatusBadRequest)
		return
	}

	portfolio, err := h.contentService.GetPortfolioByID(r.Context(), id)
	if err != nil {
		response.Error(w, "portfolio not found", http.StatusNotFound)
		return
	}

	response.Success(w, "portfolio retrieved", portfolio, http.StatusOK)
}

func (h *AdminHandler) CreatePortfolio(w http.ResponseWriter, r *http.Request) {
	var portfolio models.Portfolio
	if err := json.NewDecoder(r.Body).Decode(&portfolio); err != nil {
		response.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	if err := h.contentService.CreatePortfolio(r.Context(), &portfolio); err != nil {
		response.Error(w, "failed to create portfolio", http.StatusInternalServerError)
		return
	}

	response.Success(w, "portfolio created", portfolio, http.StatusCreated)
}

func (h *AdminHandler) UpdatePortfolio(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	if idStr == "" {
		response.Error(w, "portfolio ID is required", http.StatusBadRequest)
		return
	}

	id, err := uuid.Parse(idStr)
	if err != nil {
		response.Error(w, "invalid portfolio ID", http.StatusBadRequest)
		return
	}

	var portfolio models.Portfolio
	if err := json.NewDecoder(r.Body).Decode(&portfolio); err != nil {
		response.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	portfolio.ID = id

	if err := h.contentService.UpdatePortfolio(r.Context(), &portfolio); err != nil {
		response.Error(w, "failed to update portfolio", http.StatusInternalServerError)
		return
	}

	response.Success(w, "portfolio updated", portfolio, http.StatusOK)
}

func (h *AdminHandler) DeletePortfolio(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	if idStr == "" {
		response.Error(w, "portfolio ID is required", http.StatusBadRequest)
		return
	}

	id, err := uuid.Parse(idStr)
	if err != nil {
		response.Error(w, "invalid portfolio ID", http.StatusBadRequest)
		return
	}

	if err := h.contentService.DeletePortfolio(r.Context(), id); err != nil {
		response.Error(w, "failed to delete portfolio", http.StatusInternalServerError)
		return
	}

	response.Success(w, "portfolio deleted", nil, http.StatusOK)
}

func (h *AdminHandler) GetAllTeam(w http.ResponseWriter, r *http.Request) {
	visibleOnly := r.URL.Query().Get("visible_only") == "true"

	team, err := h.contentService.GetAllTeam(r.Context(), visibleOnly)
	if err != nil {
		response.Error(w, "failed to get team", http.StatusInternalServerError)
		return
	}

	response.Success(w, "team retrieved", team, http.StatusOK)
}

func (h *AdminHandler) GetTeam(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	if idStr == "" {
		response.Error(w, "team ID is required", http.StatusBadRequest)
		return
	}

	id, err := uuid.Parse(idStr)
	if err != nil {
		response.Error(w, "invalid team ID", http.StatusBadRequest)
		return
	}

	team, err := h.contentService.GetTeamByID(r.Context(), id)
	if err != nil {
		response.Error(w, "team member not found", http.StatusNotFound)
		return
	}

	response.Success(w, "team member retrieved", team, http.StatusOK)
}

func (h *AdminHandler) CreateTeam(w http.ResponseWriter, r *http.Request) {
	var team models.Team
	if err := json.NewDecoder(r.Body).Decode(&team); err != nil {
		response.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	if err := h.contentService.CreateTeam(r.Context(), &team); err != nil {
		response.Error(w, "failed to create team member", http.StatusInternalServerError)
		return
	}

	response.Success(w, "team member created", team, http.StatusCreated)
}

func (h *AdminHandler) UpdateTeam(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	if idStr == "" {
		response.Error(w, "team ID is required", http.StatusBadRequest)
		return
	}

	id, err := uuid.Parse(idStr)
	if err != nil {
		response.Error(w, "invalid team ID", http.StatusBadRequest)
		return
	}

	var team models.Team
	if err := json.NewDecoder(r.Body).Decode(&team); err != nil {
		response.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	team.ID = id

	if err := h.contentService.UpdateTeam(r.Context(), &team); err != nil {
		response.Error(w, "failed to update team member", http.StatusInternalServerError)
		return
	}

	response.Success(w, "team member updated", team, http.StatusOK)
}

func (h *AdminHandler) DeleteTeam(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	if idStr == "" {
		response.Error(w, "team ID is required", http.StatusBadRequest)
		return
	}

	id, err := uuid.Parse(idStr)
	if err != nil {
		response.Error(w, "invalid team ID", http.StatusBadRequest)
		return
	}

	if err := h.contentService.DeleteTeam(r.Context(), id); err != nil {
		response.Error(w, "failed to delete team member", http.StatusInternalServerError)
		return
	}

	response.Success(w, "team member deleted", nil, http.StatusOK)
}

func (h *AdminHandler) GetAllAwards(w http.ResponseWriter, r *http.Request) {
	visibleOnly := r.URL.Query().Get("visible_only") == "true"

	awards, err := h.contentService.GetAllAwards(r.Context(), visibleOnly)
	if err != nil {
		response.Error(w, "failed to get awards", http.StatusInternalServerError)
		return
	}

	response.Success(w, "awards retrieved", awards, http.StatusOK)
}

func (h *AdminHandler) GetAward(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	if idStr == "" {
		response.Error(w, "award ID is required", http.StatusBadRequest)
		return
	}

	id, err := uuid.Parse(idStr)
	if err != nil {
		response.Error(w, "invalid award ID", http.StatusBadRequest)
		return
	}

	award, err := h.contentService.GetAwardByID(r.Context(), id)
	if err != nil {
		response.Error(w, "award not found", http.StatusNotFound)
		return
	}

	response.Success(w, "award retrieved", award, http.StatusOK)
}

func (h *AdminHandler) CreateAward(w http.ResponseWriter, r *http.Request) {
	var award models.Award
	if err := json.NewDecoder(r.Body).Decode(&award); err != nil {
		response.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	if err := h.contentService.CreateAward(r.Context(), &award); err != nil {
		response.Error(w, "failed to create award", http.StatusInternalServerError)
		return
	}

	response.Success(w, "award created", award, http.StatusCreated)
}

func (h *AdminHandler) UpdateAward(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	if idStr == "" {
		response.Error(w, "award ID is required", http.StatusBadRequest)
		return
	}

	id, err := uuid.Parse(idStr)
	if err != nil {
		response.Error(w, "invalid award ID", http.StatusBadRequest)
		return
	}

	var award models.Award
	if err := json.NewDecoder(r.Body).Decode(&award); err != nil {
		response.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	award.ID = id

	if err := h.contentService.UpdateAward(r.Context(), &award); err != nil {
		response.Error(w, "failed to update award", http.StatusInternalServerError)
		return
	}

	response.Success(w, "award updated", award, http.StatusOK)
}

func (h *AdminHandler) DeleteAward(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	if idStr == "" {
		response.Error(w, "award ID is required", http.StatusBadRequest)
		return
	}

	id, err := uuid.Parse(idStr)
	if err != nil {
		response.Error(w, "invalid award ID", http.StatusBadRequest)
		return
	}

	if err := h.contentService.DeleteAward(r.Context(), id); err != nil {
		response.Error(w, "failed to delete award", http.StatusInternalServerError)
		return
	}

	response.Success(w, "award deleted", nil, http.StatusOK)
}

func (h *AdminHandler) GetAllFAQs(w http.ResponseWriter, r *http.Request) {
	visibleOnly := r.URL.Query().Get("visible_only") == "true"

	faqs, err := h.contentService.GetAllFAQs(r.Context(), visibleOnly)
	if err != nil {
		response.Error(w, "failed to get FAQs", http.StatusInternalServerError)
		return
	}

	response.Success(w, "FAQs retrieved", faqs, http.StatusOK)
}

func (h *AdminHandler) GetFAQ(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	if idStr == "" {
		response.Error(w, "FAQ ID is required", http.StatusBadRequest)
		return
	}

	id, err := uuid.Parse(idStr)
	if err != nil {
		response.Error(w, "invalid FAQ ID", http.StatusBadRequest)
		return
	}

	faq, err := h.contentService.GetFAQByID(r.Context(), id)
	if err != nil {
		response.Error(w, "FAQ not found", http.StatusNotFound)
		return
	}

	response.Success(w, "FAQ retrieved", faq, http.StatusOK)
}

func (h *AdminHandler) CreateFAQ(w http.ResponseWriter, r *http.Request) {
	var faq models.FAQ
	if err := json.NewDecoder(r.Body).Decode(&faq); err != nil {
		response.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	if err := h.contentService.CreateFAQ(r.Context(), &faq); err != nil {
		response.Error(w, "failed to create FAQ", http.StatusInternalServerError)
		return
	}

	response.Success(w, "FAQ created", faq, http.StatusCreated)
}

func (h *AdminHandler) UpdateFAQ(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	if idStr == "" {
		response.Error(w, "FAQ ID is required", http.StatusBadRequest)
		return
	}

	id, err := uuid.Parse(idStr)
	if err != nil {
		response.Error(w, "invalid FAQ ID", http.StatusBadRequest)
		return
	}

	var faq models.FAQ
	if err := json.NewDecoder(r.Body).Decode(&faq); err != nil {
		response.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	faq.ID = id

	if err := h.contentService.UpdateFAQ(r.Context(), &faq); err != nil {
		response.Error(w, "failed to update FAQ", http.StatusInternalServerError)
		return
	}

	response.Success(w, "FAQ updated", faq, http.StatusOK)
}

func (h *AdminHandler) DeleteFAQ(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	if idStr == "" {
		response.Error(w, "FAQ ID is required", http.StatusBadRequest)
		return
	}

	id, err := uuid.Parse(idStr)
	if err != nil {
		response.Error(w, "invalid FAQ ID", http.StatusBadRequest)
		return
	}

	if err := h.contentService.DeleteFAQ(r.Context(), id); err != nil {
		response.Error(w, "failed to delete FAQ", http.StatusInternalServerError)
		return
	}

	response.Success(w, "FAQ deleted", nil, http.StatusOK)
}

func (h *AdminHandler) GetAllPricing(w http.ResponseWriter, r *http.Request) {
	visibleOnly := r.URL.Query().Get("visible_only") == "true"

	pricing, err := h.contentService.GetAllPricing(r.Context(), visibleOnly)
	if err != nil {
		response.Error(w, "failed to get pricing", http.StatusInternalServerError)
		return
	}

	response.Success(w, "pricing retrieved", pricing, http.StatusOK)
}

func (h *AdminHandler) GetPricing(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	if idStr == "" {
		response.Error(w, "pricing ID is required", http.StatusBadRequest)
		return
	}

	id, err := uuid.Parse(idStr)
	if err != nil {
		response.Error(w, "invalid pricing ID", http.StatusBadRequest)
		return
	}

	pricing, err := h.contentService.GetPricingByID(r.Context(), id)
	if err != nil {
		response.Error(w, "pricing not found", http.StatusNotFound)
		return
	}

	response.Success(w, "pricing retrieved", pricing, http.StatusOK)
}

func (h *AdminHandler) CreatePricing(w http.ResponseWriter, r *http.Request) {
	var pricing models.Pricing
	if err := json.NewDecoder(r.Body).Decode(&pricing); err != nil {
		response.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	if err := h.contentService.CreatePricing(r.Context(), &pricing); err != nil {
		response.Error(w, "failed to create pricing", http.StatusInternalServerError)
		return
	}

	response.Success(w, "pricing created", pricing, http.StatusCreated)
}

func (h *AdminHandler) UpdatePricing(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	if idStr == "" {
		response.Error(w, "pricing ID is required", http.StatusBadRequest)
		return
	}

	id, err := uuid.Parse(idStr)
	if err != nil {
		response.Error(w, "invalid pricing ID", http.StatusBadRequest)
		return
	}

	var pricing models.Pricing
	if err := json.NewDecoder(r.Body).Decode(&pricing); err != nil {
		response.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	pricing.ID = id

	if err := h.contentService.UpdatePricing(r.Context(), &pricing); err != nil {
		response.Error(w, "failed to update pricing", http.StatusInternalServerError)
		return
	}

	response.Success(w, "pricing updated", pricing, http.StatusOK)
}

func (h *AdminHandler) DeletePricing(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	if idStr == "" {
		response.Error(w, "pricing ID is required", http.StatusBadRequest)
		return
	}

	id, err := uuid.Parse(idStr)
	if err != nil {
		response.Error(w, "invalid pricing ID", http.StatusBadRequest)
		return
	}

	if err := h.contentService.DeletePricing(r.Context(), id); err != nil {
		response.Error(w, "failed to delete pricing", http.StatusInternalServerError)
		return
	}

	response.Success(w, "pricing deleted", nil, http.StatusOK)
}

func (h *AdminHandler) GetAllTestimonials(w http.ResponseWriter, r *http.Request) {
	visibleOnly := r.URL.Query().Get("visible_only") == "true"

	testimonials, err := h.contentService.GetAllTestimonials(r.Context(), visibleOnly)
	if err != nil {
		response.Error(w, "failed to get testimonials", http.StatusInternalServerError)
		return
	}

	response.Success(w, "testimonials retrieved", testimonials, http.StatusOK)
}

func (h *AdminHandler) GetTestimonial(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	if idStr == "" {
		response.Error(w, "testimonial ID is required", http.StatusBadRequest)
		return
	}

	id, err := uuid.Parse(idStr)
	if err != nil {
		response.Error(w, "invalid testimonial ID", http.StatusBadRequest)
		return
	}

	testimonial, err := h.contentService.GetTestimonialByID(r.Context(), id)
	if err != nil {
		response.Error(w, "testimonial not found", http.StatusNotFound)
		return
	}

	response.Success(w, "testimonial retrieved", testimonial, http.StatusOK)
}

func (h *AdminHandler) CreateTestimonial(w http.ResponseWriter, r *http.Request) {
	var testimonial models.Testimonial
	if err := json.NewDecoder(r.Body).Decode(&testimonial); err != nil {
		response.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	if err := h.contentService.CreateTestimonial(r.Context(), &testimonial); err != nil {
		response.Error(w, "failed to create testimonial", http.StatusInternalServerError)
		return
	}

	response.Success(w, "testimonial created", testimonial, http.StatusCreated)
}

func (h *AdminHandler) UpdateTestimonial(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	if idStr == "" {
		response.Error(w, "testimonial ID is required", http.StatusBadRequest)
		return
	}

	id, err := uuid.Parse(idStr)
	if err != nil {
		response.Error(w, "invalid testimonial ID", http.StatusBadRequest)
		return
	}

	var testimonial models.Testimonial
	if err := json.NewDecoder(r.Body).Decode(&testimonial); err != nil {
		response.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	testimonial.ID = id

	if err := h.contentService.UpdateTestimonial(r.Context(), &testimonial); err != nil {
		response.Error(w, "failed to update testimonial", http.StatusInternalServerError)
		return
	}

	response.Success(w, "testimonial updated", testimonial, http.StatusOK)
}

func (h *AdminHandler) DeleteTestimonial(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	if idStr == "" {
		response.Error(w, "testimonial ID is required", http.StatusBadRequest)
		return
	}

	id, err := uuid.Parse(idStr)
	if err != nil {
		response.Error(w, "invalid testimonial ID", http.StatusBadRequest)
		return
	}

	if err := h.contentService.DeleteTestimonial(r.Context(), id); err != nil {
		response.Error(w, "failed to delete testimonial", http.StatusInternalServerError)
		return
	}

	response.Success(w, "testimonial deleted", nil, http.StatusOK)
}

func (h *AdminHandler) GetAllServices(w http.ResponseWriter, r *http.Request) {
	visibleOnly := r.URL.Query().Get("visible_only") == "true"

	services, err := h.contentService.GetAllServices(r.Context(), visibleOnly)
	if err != nil {
		response.Error(w, "failed to get services", http.StatusInternalServerError)
		return
	}

	response.Success(w, "services retrieved", services, http.StatusOK)
}

func (h *AdminHandler) GetService(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	if idStr == "" {
		response.Error(w, "service ID is required", http.StatusBadRequest)
		return
	}

	id, err := uuid.Parse(idStr)
	if err != nil {
		response.Error(w, "invalid service ID", http.StatusBadRequest)
		return
	}

	serviceItem, err := h.contentService.GetServiceByID(r.Context(), id)
	if err != nil {
		response.Error(w, "service not found", http.StatusNotFound)
		return
	}

	response.Success(w, "service retrieved", serviceItem, http.StatusOK)
}

func (h *AdminHandler) CreateService(w http.ResponseWriter, r *http.Request) {
	var serviceItem models.Service
	if err := json.NewDecoder(r.Body).Decode(&serviceItem); err != nil {
		response.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	if err := h.contentService.CreateService(r.Context(), &serviceItem); err != nil {
		response.Error(w, "failed to create service", http.StatusInternalServerError)
		return
	}

	response.Success(w, "service created", serviceItem, http.StatusCreated)
}

func (h *AdminHandler) UpdateService(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	if idStr == "" {
		response.Error(w, "service ID is required", http.StatusBadRequest)
		return
	}

	id, err := uuid.Parse(idStr)
	if err != nil {
		response.Error(w, "invalid service ID", http.StatusBadRequest)
		return
	}

	var serviceItem models.Service
	if err := json.NewDecoder(r.Body).Decode(&serviceItem); err != nil {
		response.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	serviceItem.ID = id

	if err := h.contentService.UpdateService(r.Context(), &serviceItem); err != nil {
		response.Error(w, "failed to update service", http.StatusInternalServerError)
		return
	}

	response.Success(w, "service updated", serviceItem, http.StatusOK)
}

func (h *AdminHandler) DeleteService(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	if idStr == "" {
		response.Error(w, "service ID is required", http.StatusBadRequest)
		return
	}

	id, err := uuid.Parse(idStr)
	if err != nil {
		response.Error(w, "invalid service ID", http.StatusBadRequest)
		return
	}

	if err := h.contentService.DeleteService(r.Context(), id); err != nil {
		response.Error(w, "failed to delete service", http.StatusInternalServerError)
		return
	}

	response.Success(w, "service deleted", nil, http.StatusOK)
}

func (h *AdminHandler) GetAllContacts(w http.ResponseWriter, r *http.Request) {
	page, pageSize := parsePagination(r)
	statusFilter := r.URL.Query().Get("status")

	contacts, total, err := h.contentService.GetAllContacts(r.Context(), statusFilter, page, pageSize)
	if err != nil {
		response.Error(w, "failed to get contacts", http.StatusInternalServerError)
		return
	}

	response.Paginated(w, "contacts retrieved", contacts, total, int64(page), int64(pageSize))
}

func (h *AdminHandler) GetContact(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	if idStr == "" {
		response.Error(w, "contact ID is required", http.StatusBadRequest)
		return
	}

	id, err := uuid.Parse(idStr)
	if err != nil {
		response.Error(w, "invalid contact ID", http.StatusBadRequest)
		return
	}

	contact, err := h.contentService.GetContactByID(r.Context(), id)
	if err != nil {
		response.Error(w, "contact not found", http.StatusNotFound)
		return
	}

	response.Success(w, "contact retrieved", contact, http.StatusOK)
}

func (h *AdminHandler) UpdateContact(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	if idStr == "" {
		response.Error(w, "contact ID is required", http.StatusBadRequest)
		return
	}

	id, err := uuid.Parse(idStr)
	if err != nil {
		response.Error(w, "invalid contact ID", http.StatusBadRequest)
		return
	}

	var req struct {
		Status     string  `json:"status"`
		AdminNotes *string `json:"admin_notes"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	if req.Status != "" {
		if err := h.contentService.UpdateContactStatus(r.Context(), id, req.Status); err != nil {
			response.Error(w, "failed to update contact status", http.StatusInternalServerError)
			return
		}
	}

	if req.AdminNotes != nil {
		if err := h.contentService.UpdateContactNotes(r.Context(), id, *req.AdminNotes); err != nil {
			response.Error(w, "failed to update contact notes", http.StatusInternalServerError)
			return
		}
	}

	response.Success(w, "contact updated", nil, http.StatusOK)
}

func (h *AdminHandler) DeleteContact(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	if idStr == "" {
		response.Error(w, "contact ID is required", http.StatusBadRequest)
		return
	}

	id, err := uuid.Parse(idStr)
	if err != nil {
		response.Error(w, "invalid contact ID", http.StatusBadRequest)
		return
	}

	if err := h.contentService.DeleteContact(r.Context(), id); err != nil {
		response.Error(w, "failed to delete contact", http.StatusInternalServerError)
		return
	}

	response.Success(w, "contact deleted", nil, http.StatusOK)
}

func (h *AdminHandler) GetSEO(w http.ResponseWriter, r *http.Request) {
	seo, err := h.contentService.GetSEOSettings(r.Context())
	if err != nil {
		response.Error(w, "failed to get SEO settings", http.StatusNotFound)
		return
	}

	response.Success(w, "SEO settings retrieved", seo, http.StatusOK)
}

func (h *AdminHandler) UpdateSEO(w http.ResponseWriter, r *http.Request) {
	var seo models.SEO
	if err := json.NewDecoder(r.Body).Decode(&seo); err != nil {
		response.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	if err := h.contentService.UpdateSEOSettings(r.Context(), &seo); err != nil {
		response.Error(w, "failed to update SEO settings", http.StatusInternalServerError)
		return
	}

	response.Success(w, "SEO settings updated", seo, http.StatusOK)
}

func (h *AdminHandler) GetAllSocial(w http.ResponseWriter, r *http.Request) {
	visibleOnly := r.URL.Query().Get("visible_only") == "true"

	socials, err := h.contentService.GetAllSocials(r.Context(), visibleOnly)
	if err != nil {
		response.Error(w, "failed to get social links", http.StatusInternalServerError)
		return
	}

	response.Success(w, "social links retrieved", socials, http.StatusOK)
}

func (h *AdminHandler) CreateSocial(w http.ResponseWriter, r *http.Request) {
	var social models.Social
	if err := json.NewDecoder(r.Body).Decode(&social); err != nil {
		response.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	if err := h.contentService.CreateSocial(r.Context(), &social); err != nil {
		response.Error(w, "failed to create social link", http.StatusInternalServerError)
		return
	}

	response.Success(w, "social link created", social, http.StatusCreated)
}

func (h *AdminHandler) UpdateSocial(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	if idStr == "" {
		response.Error(w, "social link ID is required", http.StatusBadRequest)
		return
	}

	id, err := uuid.Parse(idStr)
	if err != nil {
		response.Error(w, "invalid social link ID", http.StatusBadRequest)
		return
	}

	var social models.Social
	if err := json.NewDecoder(r.Body).Decode(&social); err != nil {
		response.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	social.ID = id

	if err := h.contentService.UpdateSocial(r.Context(), &social); err != nil {
		response.Error(w, "failed to update social link", http.StatusInternalServerError)
		return
	}

	response.Success(w, "social link updated", social, http.StatusOK)
}

func (h *AdminHandler) DeleteSocial(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	if idStr == "" {
		response.Error(w, "social link ID is required", http.StatusBadRequest)
		return
	}

	id, err := uuid.Parse(idStr)
	if err != nil {
		response.Error(w, "invalid social link ID", http.StatusBadRequest)
		return
	}

	if err := h.contentService.DeleteSocial(r.Context(), id); err != nil {
		response.Error(w, "failed to delete social link", http.StatusInternalServerError)
		return
	}

	response.Success(w, "social link deleted", nil, http.StatusOK)
}

func (h *AdminHandler) GetAllValues(w http.ResponseWriter, r *http.Request) {
	visibleOnly := r.URL.Query().Get("visible_only") == "true"

	values, err := h.contentService.GetAllValues(r.Context(), visibleOnly)
	if err != nil {
		response.Error(w, "failed to get values", http.StatusInternalServerError)
		return
	}

	response.Success(w, "values retrieved", values, http.StatusOK)
}

func (h *AdminHandler) GetValue(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	if idStr == "" {
		response.Error(w, "value ID is required", http.StatusBadRequest)
		return
	}

	id, err := uuid.Parse(idStr)
	if err != nil {
		response.Error(w, "invalid value ID", http.StatusBadRequest)
		return
	}

	value, err := h.contentService.GetValueByID(r.Context(), id)
	if err != nil {
		response.Error(w, "value not found", http.StatusNotFound)
		return
	}

	response.Success(w, "value retrieved", value, http.StatusOK)
}

func (h *AdminHandler) CreateValue(w http.ResponseWriter, r *http.Request) {
	var value models.Value
	if err := json.NewDecoder(r.Body).Decode(&value); err != nil {
		response.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	if err := h.contentService.CreateValue(r.Context(), &value); err != nil {
		response.Error(w, "failed to create value", http.StatusInternalServerError)
		return
	}

	response.Success(w, "value created", value, http.StatusCreated)
}

func (h *AdminHandler) UpdateValue(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	if idStr == "" {
		response.Error(w, "value ID is required", http.StatusBadRequest)
		return
	}

	id, err := uuid.Parse(idStr)
	if err != nil {
		response.Error(w, "invalid value ID", http.StatusBadRequest)
		return
	}

	var value models.Value
	if err := json.NewDecoder(r.Body).Decode(&value); err != nil {
		response.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	value.ID = id

	if err := h.contentService.UpdateValue(r.Context(), &value); err != nil {
		response.Error(w, "failed to update value", http.StatusInternalServerError)
		return
	}

	response.Success(w, "value updated", value, http.StatusOK)
}

func (h *AdminHandler) DeleteValue(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	if idStr == "" {
		response.Error(w, "value ID is required", http.StatusBadRequest)
		return
	}

	id, err := uuid.Parse(idStr)
	if err != nil {
		response.Error(w, "invalid value ID", http.StatusBadRequest)
		return
	}

	if err := h.contentService.DeleteValue(r.Context(), id); err != nil {
		response.Error(w, "failed to delete value", http.StatusInternalServerError)
		return
	}

	response.Success(w, "value deleted", nil, http.StatusOK)
}

func (h *AdminHandler) GetAllJobOpenings(w http.ResponseWriter, r *http.Request) {
	activeOnly := r.URL.Query().Get("active_only") == "true"

	jobs, err := h.contentService.GetAllJobOpenings(r.Context(), activeOnly)
	if err != nil {
		response.Error(w, "failed to get job openings", http.StatusInternalServerError)
		return
	}

	response.Success(w, "job openings retrieved", jobs, http.StatusOK)
}

func (h *AdminHandler) GetJobOpening(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	if idStr == "" {
		response.Error(w, "job opening ID is required", http.StatusBadRequest)
		return
	}

	id, err := uuid.Parse(idStr)
	if err != nil {
		response.Error(w, "invalid job opening ID", http.StatusBadRequest)
		return
	}

	job, err := h.contentService.GetJobOpeningByID(r.Context(), id)
	if err != nil {
		response.Error(w, "job opening not found", http.StatusNotFound)
		return
	}

	response.Success(w, "job opening retrieved", job, http.StatusOK)
}

func (h *AdminHandler) CreateJobOpening(w http.ResponseWriter, r *http.Request) {
	var job models.JobOpening
	if err := json.NewDecoder(r.Body).Decode(&job); err != nil {
		response.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	if err := h.contentService.CreateJobOpening(r.Context(), &job); err != nil {
		response.Error(w, "failed to create job opening", http.StatusInternalServerError)
		return
	}

	response.Success(w, "job opening created", job, http.StatusCreated)
}

func (h *AdminHandler) UpdateJobOpening(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	if idStr == "" {
		response.Error(w, "job opening ID is required", http.StatusBadRequest)
		return
	}

	id, err := uuid.Parse(idStr)
	if err != nil {
		response.Error(w, "invalid job opening ID", http.StatusBadRequest)
		return
	}

	var job models.JobOpening
	if err := json.NewDecoder(r.Body).Decode(&job); err != nil {
		response.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	job.ID = id

	if err := h.contentService.UpdateJobOpening(r.Context(), &job); err != nil {
		response.Error(w, "failed to update job opening", http.StatusInternalServerError)
		return
	}

	response.Success(w, "job opening updated", job, http.StatusOK)
}

func (h *AdminHandler) DeleteJobOpening(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	if idStr == "" {
		response.Error(w, "job opening ID is required", http.StatusBadRequest)
		return
	}

	id, err := uuid.Parse(idStr)
	if err != nil {
		response.Error(w, "invalid job opening ID", http.StatusBadRequest)
		return
	}

	if err := h.contentService.DeleteJobOpening(r.Context(), id); err != nil {
		response.Error(w, "failed to delete job opening", http.StatusInternalServerError)
		return
	}

	response.Success(w, "job opening deleted", nil, http.StatusOK)
}
