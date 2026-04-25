package service

import (
	"context"
	"encoding/json"

	"github.com/appnity/backend/internal/models"
	"github.com/appnity/backend/internal/repository"
	"github.com/google/uuid"
)

type ContentService struct {
	blogRepo        *repository.BlogRepository
	portfolioRepo   *repository.PortfolioRepository
	teamRepo        *repository.TeamRepository
	serviceRepo     *repository.ServiceRepository
	faqRepo         *repository.FAQRepository
	pricingRepo     *repository.PricingRepository
	testimonialRepo *repository.TestimonialRepository
	awardRepo       *repository.AwardRepository
	valueRepo       *repository.ValueRepository
	jobOpeningRepo  *repository.JobOpeningRepository
	pageRepo        *repository.PageRepository
	componentRepo   *repository.ComponentRepository
	navigationRepo  *repository.NavigationRepository
	socialRepo      *repository.SocialRepository
	contactRepo     *repository.ContactRepository
	seoRepo         *repository.SEORepository
	mediaRepo       *repository.MediaRepository
}

func NewContentService(
	blogRepo *repository.BlogRepository,
	portfolioRepo *repository.PortfolioRepository,
	teamRepo *repository.TeamRepository,
	serviceRepo *repository.ServiceRepository,
	faqRepo *repository.FAQRepository,
	pricingRepo *repository.PricingRepository,
	testimonialRepo *repository.TestimonialRepository,
	awardRepo *repository.AwardRepository,
	valueRepo *repository.ValueRepository,
	jobOpeningRepo *repository.JobOpeningRepository,
	pageRepo *repository.PageRepository,
	componentRepo *repository.ComponentRepository,
	navigationRepo *repository.NavigationRepository,
	socialRepo *repository.SocialRepository,
	contactRepo *repository.ContactRepository,
	seoRepo *repository.SEORepository,
	mediaRepo *repository.MediaRepository,
) *ContentService {
	return &ContentService{
		blogRepo:        blogRepo,
		portfolioRepo:   portfolioRepo,
		teamRepo:        teamRepo,
		serviceRepo:     serviceRepo,
		faqRepo:         faqRepo,
		pricingRepo:     pricingRepo,
		testimonialRepo: testimonialRepo,
		awardRepo:       awardRepo,
		valueRepo:       valueRepo,
		jobOpeningRepo:  jobOpeningRepo,
		pageRepo:        pageRepo,
		componentRepo:   componentRepo,
		navigationRepo:  navigationRepo,
		socialRepo:      socialRepo,
		contactRepo:     contactRepo,
		seoRepo:         seoRepo,
		mediaRepo:       mediaRepo,
	}
}

func (s *ContentService) GetAllBlogs(ctx context.Context, publishedOnly bool, page, pageSize int) ([]models.Blog, int64, error) {
	return s.blogRepo.GetAll(ctx, publishedOnly, page, pageSize)
}

func (s *ContentService) GetBlogByID(ctx context.Context, id uuid.UUID) (*models.Blog, error) {
	return s.blogRepo.GetByID(ctx, id)
}

func (s *ContentService) GetBlogBySlug(ctx context.Context, slug string) (*models.Blog, error) {
	return s.blogRepo.GetBySlug(ctx, slug)
}

func (s *ContentService) CreateBlog(ctx context.Context, blog *models.Blog) error {
	return s.blogRepo.Create(ctx, blog)
}

func (s *ContentService) UpdateBlog(ctx context.Context, blog *models.Blog) error {
	return s.blogRepo.Update(ctx, blog)
}

func (s *ContentService) DeleteBlog(ctx context.Context, id uuid.UUID) error {
	return s.blogRepo.Delete(ctx, id)
}

func (s *ContentService) GetPublishedBlogs(ctx context.Context, page, pageSize int) ([]models.Blog, int64, error) {
	return s.blogRepo.GetAll(ctx, true, page, pageSize)
}

func (s *ContentService) GetRelatedBlogs(ctx context.Context, id uuid.UUID, category string, limit int) ([]models.Blog, error) {
	return s.blogRepo.GetRelated(ctx, id, category, limit)
}

func (s *ContentService) GetRelatedBlogsBySlug(ctx context.Context, slug string) ([]models.Blog, error) {
	return s.blogRepo.GetRelatedBySlug(ctx, slug, 5)
}

func (s *ContentService) GetAllPortfolios(ctx context.Context, visibleOnly bool, category string, page, pageSize int) ([]models.Portfolio, int64, error) {
	return s.portfolioRepo.GetAll(ctx, visibleOnly, category, page, pageSize)
}

func (s *ContentService) GetPortfolioByID(ctx context.Context, id uuid.UUID) (*models.Portfolio, error) {
	return s.portfolioRepo.GetByID(ctx, id)
}

func (s *ContentService) GetPortfolioBySlug(ctx context.Context, slug string) (*models.Portfolio, error) {
	return s.portfolioRepo.GetBySlug(ctx, slug)
}

func (s *ContentService) CreatePortfolio(ctx context.Context, portfolio *models.Portfolio) error {
	return s.portfolioRepo.Create(ctx, portfolio)
}

func (s *ContentService) UpdatePortfolio(ctx context.Context, portfolio *models.Portfolio) error {
	return s.portfolioRepo.Update(ctx, portfolio)
}

func (s *ContentService) DeletePortfolio(ctx context.Context, id uuid.UUID) error {
	return s.portfolioRepo.Delete(ctx, id)
}

func (s *ContentService) GetVisiblePortfolios(ctx context.Context, page, pageSize int) ([]models.Portfolio, int64, error) {
	return s.portfolioRepo.GetAll(ctx, true, "", page, pageSize)
}

func (s *ContentService) GetFeaturedPortfolios(ctx context.Context) ([]models.Portfolio, error) {
	return s.portfolioRepo.GetFeatured(ctx)
}

func (s *ContentService) GetPortfoliosByCategory(ctx context.Context, category string, page, pageSize int) ([]models.Portfolio, int64, error) {
	return s.portfolioRepo.GetAll(ctx, true, category, page, pageSize)
}

func (s *ContentService) ReorderPortfolios(ctx context.Context, ids []uuid.UUID) error {
	return s.portfolioRepo.Reorder(ctx, ids)
}

func (s *ContentService) GetAllTeam(ctx context.Context, visibleOnly bool) ([]models.Team, error) {
	return s.teamRepo.GetAll(ctx, visibleOnly)
}

func (s *ContentService) GetTeamByID(ctx context.Context, id uuid.UUID) (*models.Team, error) {
	return s.teamRepo.GetByID(ctx, id)
}

func (s *ContentService) CreateTeam(ctx context.Context, team *models.Team) error {
	return s.teamRepo.Create(ctx, team)
}

func (s *ContentService) UpdateTeam(ctx context.Context, team *models.Team) error {
	return s.teamRepo.Update(ctx, team)
}

func (s *ContentService) DeleteTeam(ctx context.Context, id uuid.UUID) error {
	return s.teamRepo.Delete(ctx, id)
}

func (s *ContentService) ReorderTeam(ctx context.Context, ids []uuid.UUID) error {
	return s.teamRepo.Reorder(ctx, ids)
}

func (s *ContentService) GetAllServices(ctx context.Context, visibleOnly bool) ([]models.Service, error) {
	return s.serviceRepo.GetAll(ctx, visibleOnly)
}

func (s *ContentService) GetServiceByID(ctx context.Context, id uuid.UUID) (*models.Service, error) {
	return s.serviceRepo.GetByID(ctx, id)
}

func (s *ContentService) CreateService(ctx context.Context, service *models.Service) error {
	return s.serviceRepo.Create(ctx, service)
}

func (s *ContentService) UpdateService(ctx context.Context, service *models.Service) error {
	return s.serviceRepo.Update(ctx, service)
}

func (s *ContentService) DeleteService(ctx context.Context, id uuid.UUID) error {
	return s.serviceRepo.Delete(ctx, id)
}

func (s *ContentService) ReorderServices(ctx context.Context, ids []uuid.UUID) error {
	return s.serviceRepo.Reorder(ctx, ids)
}

func (s *ContentService) GetAllFAQs(ctx context.Context, visibleOnly bool) ([]models.FAQ, error) {
	return s.faqRepo.GetAll(ctx, visibleOnly)
}

func (s *ContentService) GetFAQByID(ctx context.Context, id uuid.UUID) (*models.FAQ, error) {
	return s.faqRepo.GetByID(ctx, id)
}

func (s *ContentService) CreateFAQ(ctx context.Context, faq *models.FAQ) error {
	return s.faqRepo.Create(ctx, faq)
}

func (s *ContentService) UpdateFAQ(ctx context.Context, faq *models.FAQ) error {
	return s.faqRepo.Update(ctx, faq)
}

func (s *ContentService) DeleteFAQ(ctx context.Context, id uuid.UUID) error {
	return s.faqRepo.Delete(ctx, id)
}

func (s *ContentService) ReorderFAQs(ctx context.Context, ids []uuid.UUID) error {
	return s.faqRepo.Reorder(ctx, ids)
}

func (s *ContentService) GetAllPricing(ctx context.Context, visibleOnly bool) ([]models.Pricing, error) {
	return s.pricingRepo.GetAll(ctx, visibleOnly)
}

func (s *ContentService) GetPricingByID(ctx context.Context, id uuid.UUID) (*models.Pricing, error) {
	return s.pricingRepo.GetByID(ctx, id)
}

func (s *ContentService) CreatePricing(ctx context.Context, pricing *models.Pricing) error {
	return s.pricingRepo.Create(ctx, pricing)
}

func (s *ContentService) UpdatePricing(ctx context.Context, pricing *models.Pricing) error {
	return s.pricingRepo.Update(ctx, pricing)
}

func (s *ContentService) DeletePricing(ctx context.Context, id uuid.UUID) error {
	return s.pricingRepo.Delete(ctx, id)
}

func (s *ContentService) ReorderPricing(ctx context.Context, ids []uuid.UUID) error {
	return s.pricingRepo.Reorder(ctx, ids)
}

func (s *ContentService) GetAllTestimonials(ctx context.Context, visibleOnly bool) ([]models.Testimonial, error) {
	return s.testimonialRepo.GetAll(ctx, visibleOnly)
}

func (s *ContentService) GetTestimonialByID(ctx context.Context, id uuid.UUID) (*models.Testimonial, error) {
	return s.testimonialRepo.GetByID(ctx, id)
}

func (s *ContentService) CreateTestimonial(ctx context.Context, testimonial *models.Testimonial) error {
	return s.testimonialRepo.Create(ctx, testimonial)
}

func (s *ContentService) UpdateTestimonial(ctx context.Context, testimonial *models.Testimonial) error {
	return s.testimonialRepo.Update(ctx, testimonial)
}

func (s *ContentService) DeleteTestimonial(ctx context.Context, id uuid.UUID) error {
	return s.testimonialRepo.Delete(ctx, id)
}

func (s *ContentService) ReorderTestimonials(ctx context.Context, ids []uuid.UUID) error {
	return s.testimonialRepo.Reorder(ctx, ids)
}

func (s *ContentService) GetAllAwards(ctx context.Context, visibleOnly bool) ([]models.Award, error) {
	return s.awardRepo.GetAll(ctx, visibleOnly)
}

func (s *ContentService) GetAwardByID(ctx context.Context, id uuid.UUID) (*models.Award, error) {
	return s.awardRepo.GetByID(ctx, id)
}

func (s *ContentService) CreateAward(ctx context.Context, award *models.Award) error {
	return s.awardRepo.Create(ctx, award)
}

func (s *ContentService) UpdateAward(ctx context.Context, award *models.Award) error {
	return s.awardRepo.Update(ctx, award)
}

func (s *ContentService) DeleteAward(ctx context.Context, id uuid.UUID) error {
	return s.awardRepo.Delete(ctx, id)
}

func (s *ContentService) ReorderAwards(ctx context.Context, ids []uuid.UUID) error {
	return s.awardRepo.Reorder(ctx, ids)
}

func (s *ContentService) GetAllValues(ctx context.Context, visibleOnly bool) ([]models.Value, error) {
	return s.valueRepo.GetAll(ctx, visibleOnly)
}

func (s *ContentService) GetValueByID(ctx context.Context, id uuid.UUID) (*models.Value, error) {
	return s.valueRepo.GetByID(ctx, id)
}

func (s *ContentService) CreateValue(ctx context.Context, value *models.Value) error {
	return s.valueRepo.Create(ctx, value)
}

func (s *ContentService) UpdateValue(ctx context.Context, value *models.Value) error {
	return s.valueRepo.Update(ctx, value)
}

func (s *ContentService) DeleteValue(ctx context.Context, id uuid.UUID) error {
	return s.valueRepo.Delete(ctx, id)
}

func (s *ContentService) ReorderValues(ctx context.Context, ids []uuid.UUID) error {
	return s.valueRepo.Reorder(ctx, ids)
}

func (s *ContentService) GetAllJobOpenings(ctx context.Context, activeOnly bool) ([]models.JobOpening, error) {
	return s.jobOpeningRepo.GetAll(ctx, activeOnly)
}

func (s *ContentService) GetJobOpeningByID(ctx context.Context, id uuid.UUID) (*models.JobOpening, error) {
	return s.jobOpeningRepo.GetByID(ctx, id)
}

func (s *ContentService) CreateJobOpening(ctx context.Context, job *models.JobOpening) error {
	return s.jobOpeningRepo.Create(ctx, job)
}

func (s *ContentService) UpdateJobOpening(ctx context.Context, job *models.JobOpening) error {
	return s.jobOpeningRepo.Update(ctx, job)
}

func (s *ContentService) DeleteJobOpening(ctx context.Context, id uuid.UUID) error {
	return s.jobOpeningRepo.Delete(ctx, id)
}

func (s *ContentService) GetAllPages(ctx context.Context) ([]models.Page, error) {
	return s.pageRepo.GetAll(ctx)
}

func (s *ContentService) GetPageByID(ctx context.Context, id uuid.UUID) (*models.Page, error) {
	return s.pageRepo.GetByID(ctx, id)
}

func (s *ContentService) GetPageBySlug(ctx context.Context, slug string) (*models.Page, error) {
	return s.pageRepo.GetBySlug(ctx, slug)
}

func (s *ContentService) CreatePage(ctx context.Context, page *models.Page) error {
	return s.pageRepo.Create(ctx, page)
}

func (s *ContentService) UpdatePage(ctx context.Context, page *models.Page) error {
	return s.pageRepo.Update(ctx, page)
}

func (s *ContentService) DeletePage(ctx context.Context, id uuid.UUID) error {
	return s.pageRepo.Delete(ctx, id)
}

func (s *ContentService) GetAllComponents(ctx context.Context) ([]models.Component, error) {
	return s.componentRepo.GetAll(ctx)
}

func (s *ContentService) GetComponentsByPageID(ctx context.Context, pageID uuid.UUID) ([]models.Component, error) {
	return s.componentRepo.GetByPageID(ctx, pageID)
}

func (s *ContentService) GetComponentByID(ctx context.Context, id uuid.UUID) (*models.Component, error) {
	return s.componentRepo.GetByID(ctx, id)
}

func (s *ContentService) CreateComponent(ctx context.Context, component *models.Component) error {
	return s.componentRepo.Create(ctx, component)
}

func (s *ContentService) UpdateComponent(ctx context.Context, component *models.Component) error {
	return s.componentRepo.Update(ctx, component)
}

func (s *ContentService) DeleteComponent(ctx context.Context, id uuid.UUID) error {
	return s.componentRepo.Delete(ctx, id)
}

func (s *ContentService) ReorderComponents(ctx context.Context, pageID uuid.UUID, ids []uuid.UUID) error {
	return s.componentRepo.Reorder(ctx, pageID, ids)
}

func (s *ContentService) GetPageComponentsBySlug(ctx context.Context, slug string) ([]models.Component, error) {
	return s.componentRepo.GetByPageSlug(ctx, slug)
}

func (s *ContentService) ToggleComponentVisibility(ctx context.Context, id uuid.UUID) error {
	return s.componentRepo.ToggleVisibility(ctx, id)
}

func (s *ContentService) GetAllNavigation(ctx context.Context, visibleOnly bool) ([]models.Navigation, error) {
	return s.navigationRepo.GetAll(ctx, visibleOnly)
}

func (s *ContentService) CreateNavigation(ctx context.Context, nav *models.Navigation) error {
	return s.navigationRepo.Create(ctx, nav)
}

func (s *ContentService) UpdateNavigation(ctx context.Context, nav *models.Navigation) error {
	return s.navigationRepo.Update(ctx, nav)
}

func (s *ContentService) DeleteNavigation(ctx context.Context, id uuid.UUID) error {
	return s.navigationRepo.Delete(ctx, id)
}

func (s *ContentService) ReorderNavigation(ctx context.Context, ids []uuid.UUID) error {
	return s.navigationRepo.Reorder(ctx, ids)
}

func (s *ContentService) GetAllSocials(ctx context.Context, visibleOnly bool) ([]models.Social, error) {
	return s.socialRepo.GetAll(ctx, visibleOnly)
}

func (s *ContentService) GetSocialByID(ctx context.Context, id uuid.UUID) (*models.Social, error) {
	return s.socialRepo.GetByID(ctx, id)
}

func (s *ContentService) CreateSocial(ctx context.Context, social *models.Social) error {
	return s.socialRepo.Create(ctx, social)
}

func (s *ContentService) UpdateSocial(ctx context.Context, social *models.Social) error {
	return s.socialRepo.Update(ctx, social)
}

func (s *ContentService) DeleteSocial(ctx context.Context, id uuid.UUID) error {
	return s.socialRepo.Delete(ctx, id)
}

func (s *ContentService) ReorderSocials(ctx context.Context, ids []uuid.UUID) error {
	return s.socialRepo.Reorder(ctx, ids)
}

func (s *ContentService) GetAllContacts(ctx context.Context, statusFilter string, page, pageSize int) ([]models.Contact, int64, error) {
	return s.contactRepo.GetAll(ctx, statusFilter, page, pageSize)
}

func (s *ContentService) GetContactByID(ctx context.Context, id uuid.UUID) (*models.Contact, error) {
	return s.contactRepo.GetByID(ctx, id)
}

func (s *ContentService) UpdateContact(ctx context.Context, contact *models.Contact) error {
	return s.contactRepo.Update(ctx, contact)
}

func (s *ContentService) DeleteContact(ctx context.Context, id uuid.UUID) error {
	return s.contactRepo.Delete(ctx, id)
}

func (s *ContentService) GetContactSubmissions(ctx context.Context, statusFilter string, page, pageSize int) ([]models.Contact, int64, error) {
	return s.contactRepo.GetAll(ctx, statusFilter, page, pageSize)
}

func (s *ContentService) CreateContactSubmission(ctx context.Context, contact *models.Contact) error {
	return s.contactRepo.Create(ctx, contact)
}

func (s *ContentService) GetSEOSettings(ctx context.Context) (*models.SEO, error) {
	return s.seoRepo.Get(ctx)
}

func (s *ContentService) UpdateSEOSettings(ctx context.Context, seo *models.SEO) error {
	existing, err := s.seoRepo.Get(ctx)
	if err != nil {
		return s.seoRepo.Create(ctx, seo)
	}
	seo.ID = existing.ID
	return s.seoRepo.Update(ctx, seo)
}

func (s *ContentService) GetAllMedia(ctx context.Context, folder string, page, pageSize int) ([]models.Media, int64, error) {
	return s.mediaRepo.GetAll(ctx, folder, page, pageSize)
}

func (s *ContentService) GetMediaByID(ctx context.Context, id uuid.UUID) (*models.Media, error) {
	return s.mediaRepo.GetByID(ctx, id)
}

func (s *ContentService) UpdateMedia(ctx context.Context, media *models.Media) error {
	return s.mediaRepo.Update(ctx, media)
}

func (s *ContentService) DeleteMedia(ctx context.Context, id uuid.UUID) error {
	return s.mediaRepo.Delete(ctx, id)
}

func (s *ContentService) UpdateContactStatus(ctx context.Context, id uuid.UUID, status string) error {
	contact, err := s.contactRepo.GetByID(ctx, id)
	if err != nil {
		return err
	}
	contact.Status = status
	return s.contactRepo.Update(ctx, contact)
}

func (s *ContentService) UpdateContactNotes(ctx context.Context, id uuid.UUID, notes string) error {
	contact, err := s.contactRepo.GetByID(ctx, id)
	if err != nil {
		return err
	}
	contact.AdminNotes = &notes
	return s.contactRepo.Update(ctx, contact)
}

func (s *ContentService) GetContactStatusCounts(ctx context.Context) (map[string]int64, error) {
	return s.contactRepo.CountByStatus(ctx)
}

func (s *ContentService) UpdatePageSections(ctx context.Context, id uuid.UUID, sections json.RawMessage) error {
	page, err := s.pageRepo.GetByID(ctx, id)
	if err != nil {
		return err
	}
	page.Sections = sections
	return s.pageRepo.Update(ctx, page)
}
