package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/appnity/backend/internal/models"
	"github.com/appnity/backend/internal/service"
	"github.com/appnity/backend/pkg/response"
)

type PublicHandler struct {
	contentService *service.ContentService
}

func NewPublicHandler(contentService *service.ContentService) *PublicHandler {
	return &PublicHandler{contentService: contentService}
}

func (h *PublicHandler) GetTheme(w http.ResponseWriter, r *http.Request) {
	theme, err := h.contentService.GetTheme(r.Context())
	if err != nil {
		response.Error(w, "failed to get theme", http.StatusInternalServerError)
		return
	}
	response.Success(w, "theme retrieved", theme, http.StatusOK)
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
	slug := r.PathValue("slug")
	if slug == "" {
		response.Error(w, "page slug is required", http.StatusBadRequest)
		return
	}

	components, err := h.contentService.GetPageComponentsBySlug(r.Context(), slug)
	if err != nil {
		response.Error(w, "failed to get components", http.StatusInternalServerError)
		return
	}

	response.Success(w, "components retrieved", components, http.StatusOK)
}

func (h *PublicHandler) GetBlogs(w http.ResponseWriter, r *http.Request) {
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	if page < 1 {
		page = 1
	}
	pageSize, _ := strconv.Atoi(r.URL.Query().Get("page_size"))
	if pageSize < 1 {
		pageSize = 10
	}

	blogs, total, err := h.contentService.GetPublishedBlogs(r.Context(), page, pageSize)
	if err != nil {
		response.Error(w, "failed to get blogs", http.StatusInternalServerError)
		return
	}

	response.Success(w, "blogs retrieved", map[string]interface{}{
		"data":  blogs,
		"total": total,
	}, http.StatusOK)
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
	slug := r.PathValue("slug")
	if slug == "" {
		response.Error(w, "blog slug is required", http.StatusBadRequest)
		return
	}

	blogs, err := h.contentService.GetRelatedBlogsBySlug(r.Context(), slug)
	if err != nil {
		response.Error(w, "failed to get related blogs", http.StatusInternalServerError)
		return
	}

	response.Success(w, "related blogs retrieved", blogs, http.StatusOK)
}

func (h *PublicHandler) GetPortfolios(w http.ResponseWriter, r *http.Request) {
	category := r.URL.Query().Get("category")
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	if page < 1 {
		page = 1
	}
	pageSize, _ := strconv.Atoi(r.URL.Query().Get("page_size"))
	if pageSize < 1 {
		pageSize = 10
	}

	portfolios, total, err := h.contentService.GetAllPortfolios(r.Context(), true, category, page, pageSize)
	if err != nil {
		response.Error(w, "failed to get portfolios", http.StatusInternalServerError)
		return
	}

	response.Success(w, "portfolios retrieved", map[string]interface{}{
		"data":  portfolios,
		"total": total,
	}, http.StatusOK)
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

	response.Success(w, "faqs retrieved", faqs, http.StatusOK)
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
		response.Error(w, "failed to get SEO settings", http.StatusInternalServerError)
		return
	}

	response.Success(w, "seo retrieved", seo, http.StatusOK)
}

func (h *PublicHandler) SubmitContact(w http.ResponseWriter, r *http.Request) {
	var contact models.Contact
	if err := json.NewDecoder(r.Body).Decode(&contact); err != nil {
		response.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	if contact.Email == "" {
		response.Error(w, "email is required", http.StatusBadRequest)
		return
	}

	if err := h.contentService.CreateContactSubmission(r.Context(), &contact); err != nil {
		response.Error(w, "failed to submit contact", http.StatusInternalServerError)
		return
	}

	response.Success(w, "contact submitted successfully", nil, http.StatusCreated)
}
