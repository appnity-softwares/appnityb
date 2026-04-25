package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/appnity/backend/internal/models"
	"github.com/appnity/backend/internal/service"
	"github.com/appnity/backend/pkg/response"
	"github.com/google/uuid"
)

type PublicHandler struct {
	contentService *service.ContentService
	themeService   *service.ThemeService
}

func NewPublicHandler(contentService *service.ContentService, themeService *service.ThemeService) *PublicHandler {
	return &PublicHandler{
		contentService: contentService,
		themeService:   themeService,
	}
}

func (h *PublicHandler) GetTheme(w http.ResponseWriter, r *http.Request) {
	theme, err := h.themeService.GetActiveTheme(r.Context())
	if err != nil {
		theme = h.getDefaultTheme()
	}

	response.Success(w, "theme retrieved", theme, http.StatusOK)
}

func (h *PublicHandler) getDefaultTheme() *models.Theme {
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
		Name:     "Default Dark Orange",
		IsActive: true,
		Colors:   colorsBytes,
		Fonts:    fontsBytes,
	}
}

func (h *PublicHandler) GetNavigation(w http.ResponseWriter, r *http.Request) {
	nav, err := h.contentService.GetAllNavigation(r.Context(), true)
	if err != nil {
		response.Error(w, "failed to get navigation", http.StatusInternalServerError)
		return
	}

	response.Success(w, "navigation retrieved", nav, http.StatusOK)
}

func (h *PublicHandler) GetPage(w http.ResponseWriter, r *http.Request) {
	slug := r.PathValue("slug")
	if slug == "" {
		response.Error(w, "page slug is required", http.StatusBadRequest)
		return
	}

	page, err := h.contentService.GetPageBySlug(r.Context(), slug)
	if err != nil {
		response.Error(w, "page not found", http.StatusNotFound)
		return
	}

	response.Success(w, "page retrieved", page, http.StatusOK)
}

func (h *PublicHandler) GetPageComponents(w http.ResponseWriter, r *http.Request) {
	pageIDStr := r.PathValue("page_id")
	if pageIDStr == "" {
		response.Error(w, "page_id is required", http.StatusBadRequest)
		return
	}

	pageID, err := uuid.Parse(pageIDStr)
	if err != nil {
		response.Error(w, "invalid page ID", http.StatusBadRequest)
		return
	}

	components, err := h.contentService.GetComponentsByPageID(r.Context(), pageID)
	if err != nil {
		response.Error(w, "failed to get components", http.StatusInternalServerError)
		return
	}

	response.Success(w, "components retrieved", components, http.StatusOK)
}

func (h *PublicHandler) GetBlogs(w http.ResponseWriter, r *http.Request) {
	page, pageSize := parsePagination(r)

	blogs, total, err := h.contentService.GetPublishedBlogs(r.Context(), page, pageSize)
	if err != nil {
		response.Error(w, "failed to get blogs", http.StatusInternalServerError)
		return
	}

	response.Paginated(w, "blogs retrieved", blogs, total, int64(page), int64(pageSize))
}

func (h *PublicHandler) GetBlogBySlug(w http.ResponseWriter, r *http.Request) {
	slug := r.PathValue("slug")
	if slug == "" {
		response.Error(w, "blog slug is required", http.StatusBadRequest)
		return
	}

	blog, err := h.contentService.GetBlogBySlug(r.Context(), slug)
	if err != nil {
		response.Error(w, "blog not found", http.StatusNotFound)
		return
	}

	response.Success(w, "blog retrieved", blog, http.StatusOK)
}

func (h *PublicHandler) GetRelatedBlogs(w http.ResponseWriter, r *http.Request) {
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

	category := r.URL.Query().Get("category")
	limitStr := r.URL.Query().Get("limit")
	limit := 5
	if limitStr != "" {
		limit, _ = strconv.Atoi(limitStr)
	}

	blogs, err := h.contentService.GetRelatedBlogs(r.Context(), id, category, limit)
	if err != nil {
		response.Error(w, "failed to get related blogs", http.StatusInternalServerError)
		return
	}

	response.Success(w, "related blogs retrieved", blogs, http.StatusOK)
}

func (h *PublicHandler) GetPortfolios(w http.ResponseWriter, r *http.Request) {
	page, pageSize := parsePagination(r)
	category := r.URL.Query().Get("category")

	portfolios, total, err := h.contentService.GetPortfoliosByCategory(r.Context(), category, page, pageSize)
	if err != nil {
		response.Error(w, "failed to get portfolios", http.StatusInternalServerError)
		return
	}

	response.Paginated(w, "portfolios retrieved", portfolios, total, int64(page), int64(pageSize))
}

func (h *PublicHandler) GetFeaturedPortfolio(w http.ResponseWriter, r *http.Request) {
	portfolios, err := h.contentService.GetFeaturedPortfolios(r.Context())
	if err != nil {
		response.Error(w, "failed to get featured portfolios", http.StatusInternalServerError)
		return
	}

	response.Success(w, "featured portfolios retrieved", portfolios, http.StatusOK)
}

func (h *PublicHandler) GetPortfolioBySlug(w http.ResponseWriter, r *http.Request) {
	slug := r.PathValue("slug")
	if slug == "" {
		response.Error(w, "portfolio slug is required", http.StatusBadRequest)
		return
	}

	portfolio, err := h.contentService.GetPortfolioBySlug(r.Context(), slug)
	if err != nil {
		response.Error(w, "portfolio not found", http.StatusNotFound)
		return
	}

	response.Success(w, "portfolio retrieved", portfolio, http.StatusOK)
}

func (h *PublicHandler) GetTeam(w http.ResponseWriter, r *http.Request) {
	team, err := h.contentService.GetAllTeam(r.Context(), true)
	if err != nil {
		response.Error(w, "failed to get team", http.StatusInternalServerError)
		return
	}

	response.Success(w, "team retrieved", team, http.StatusOK)
}

func (h *PublicHandler) GetServices(w http.ResponseWriter, r *http.Request) {
	services, err := h.contentService.GetAllServices(r.Context(), true)
	if err != nil {
		response.Error(w, "failed to get services", http.StatusInternalServerError)
		return
	}

	response.Success(w, "services retrieved", services, http.StatusOK)
}

func (h *PublicHandler) GetFAQs(w http.ResponseWriter, r *http.Request) {
	faqs, err := h.contentService.GetAllFAQs(r.Context(), true)
	if err != nil {
		response.Error(w, "failed to get FAQs", http.StatusInternalServerError)
		return
	}

	response.Success(w, "FAQs retrieved", faqs, http.StatusOK)
}

func (h *PublicHandler) GetPricing(w http.ResponseWriter, r *http.Request) {
	pricing, err := h.contentService.GetAllPricing(r.Context(), true)
	if err != nil {
		response.Error(w, "failed to get pricing", http.StatusInternalServerError)
		return
	}

	response.Success(w, "pricing retrieved", pricing, http.StatusOK)
}

func (h *PublicHandler) GetTestimonials(w http.ResponseWriter, r *http.Request) {
	testimonials, err := h.contentService.GetAllTestimonials(r.Context(), true)
	if err != nil {
		response.Error(w, "failed to get testimonials", http.StatusInternalServerError)
		return
	}

	response.Success(w, "testimonials retrieved", testimonials, http.StatusOK)
}

func (h *PublicHandler) GetAwards(w http.ResponseWriter, r *http.Request) {
	awards, err := h.contentService.GetAllAwards(r.Context(), true)
	if err != nil {
		response.Error(w, "failed to get awards", http.StatusInternalServerError)
		return
	}

	response.Success(w, "awards retrieved", awards, http.StatusOK)
}

func (h *PublicHandler) GetValues(w http.ResponseWriter, r *http.Request) {
	values, err := h.contentService.GetAllValues(r.Context(), true)
	if err != nil {
		response.Error(w, "failed to get values", http.StatusInternalServerError)
		return
	}

	response.Success(w, "values retrieved", values, http.StatusOK)
}

func (h *PublicHandler) GetJobOpenings(w http.ResponseWriter, r *http.Request) {
	jobs, err := h.contentService.GetAllJobOpenings(r.Context(), true)
	if err != nil {
		response.Error(w, "failed to get job openings", http.StatusInternalServerError)
		return
	}

	response.Success(w, "job openings retrieved", jobs, http.StatusOK)
}

func (h *PublicHandler) GetSocialLinks(w http.ResponseWriter, r *http.Request) {
	socials, err := h.contentService.GetAllSocials(r.Context(), true)
	if err != nil {
		response.Error(w, "failed to get social links", http.StatusInternalServerError)
		return
	}

	response.Success(w, "social links retrieved", socials, http.StatusOK)
}

func (h *PublicHandler) GetSEO(w http.ResponseWriter, r *http.Request) {
	seo, err := h.contentService.GetSEOSettings(r.Context())
	if err != nil {
		response.Error(w, "failed to get SEO settings", http.StatusNotFound)
		return
	}

	response.Success(w, "SEO settings retrieved", seo, http.StatusOK)
}

func (h *PublicHandler) SubmitContact(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Name    string `json:"name"`
		Email   string `json:"email"`
		Service string `json:"service"`
		Message string `json:"message"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	contact := models.Contact{
		Name:    req.Name,
		Email:   req.Email,
		Service: req.Service,
		Message: req.Message,
		Status:  "new",
	}

	if err := h.contentService.CreateContactSubmission(r.Context(), &contact); err != nil {
		response.Error(w, "failed to submit contact", http.StatusInternalServerError)
		return
	}

	response.Success(w, "contact form submitted", nil, http.StatusCreated)
}

func parsePagination(r *http.Request) (int, int) {
	page := 1
	pageSize := 10

	if p := r.URL.Query().Get("page"); p != "" {
		page, _ = strconv.Atoi(p)
		if page < 1 {
			page = 1
		}
	}

	if ps := r.URL.Query().Get("page_size"); ps != "" {
		pageSize, _ = strconv.Atoi(ps)
		if pageSize < 1 {
			pageSize = 10
		}
		if pageSize > 100 {
			pageSize = 100
		}
	}

	return page, pageSize
}
