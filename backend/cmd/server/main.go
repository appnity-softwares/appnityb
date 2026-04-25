package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/appnity/backend/internal/config"
	"github.com/appnity/backend/internal/database"
	"github.com/appnity/backend/internal/handler"
	"github.com/appnity/backend/internal/middleware"
	"github.com/appnity/backend/internal/repository"
	"github.com/appnity/backend/internal/service"
	jwtpkg "github.com/appnity/backend/pkg/jwt"

	"github.com/go-chi/chi/v5"
	chimiddleware "github.com/go-chi/chi/v5/middleware"
)

func main() {
	cfg := config.Load()

	db, err := database.NewDatabase(cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	jwtService := jwtpkg.NewJWTService(
		cfg.JWTSecret,
		time.Duration(cfg.JWTExpiryHours)*time.Hour,
		time.Duration(cfg.RefreshExpiryHours)*time.Hour,
	)

	userRepo := repository.NewUserRepo(db.Pool)
	themeRepo := repository.NewThemeRepo(db.Pool)
	pageRepo := repository.NewPageRepo(db.Pool)
	componentRepo := repository.NewComponentRepo(db.Pool)
	navRepo := repository.NewNavigationRepo(db.Pool)
	blogRepo := repository.NewBlogRepo(db.Pool)
	portfolioRepo := repository.NewPortfolioRepo(db.Pool)
	teamRepo := repository.NewTeamRepo(db.Pool)
	awardRepo := repository.NewAwardRepo(db.Pool)
	faqRepo := repository.NewFAQRepo(db.Pool)
	pricingRepo := repository.NewPricingRepo(db.Pool)
	testimonialRepo := repository.NewTestimonialRepo(db.Pool)
	serviceRepo := repository.NewServiceRepo(db.Pool)
	contactRepo := repository.NewContactRepo(db.Pool)
	seoRepo := repository.NewSEORepo(db.Pool)
	socialRepo := repository.NewSocialRepo(db.Pool)
	valueRepo := repository.NewValueRepo(db.Pool)
	jobOpeningRepo := repository.NewJobOpeningRepo(db.Pool)
	mediaRepo := repository.NewMediaRepo(db.Pool)

	authService := service.NewAuthService(userRepo, jwtService, cfg)
	userService := service.NewUserService(userRepo)
	contentService := service.NewContentService(
		blogRepo, portfolioRepo, teamRepo, serviceRepo, faqRepo,
		pricingRepo, testimonialRepo, awardRepo, valueRepo, jobOpeningRepo,
		pageRepo, componentRepo, navRepo, socialRepo, contactRepo, seoRepo, mediaRepo,
	)
	themeService := service.NewThemeService(themeRepo)
	mediaService := service.NewMediaService(mediaRepo, cfg.UploadDir, cfg.MaxUploadSize, cfg.BaseURL)
	configService := service.NewConfigService(contactRepo, blogRepo, portfolioRepo, teamRepo)

	authHandler := handler.NewAuthHandler(authService)
	userHandler := handler.NewUserHandler(authService, userService)
	publicHandler := handler.NewPublicHandler(contentService, themeService)
	adminHandler := handler.NewAdminHandler(contentService, themeService, authService, mediaService, configService)
	mediaHandler := handler.NewMediaHandler(mediaService)

	authMW := middleware.NewAuthMiddleware(jwtService)

	r := chi.NewRouter()

	r.Use(chimiddleware.Logger)
	r.Use(chimiddleware.Recoverer)
	r.Use(middleware.NewCORSMiddleware(cfg).Handler)
	r.Use(middleware.NewRateLimiter().Handler)

	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	r.Handle("/uploads/*", http.StripPrefix("/uploads/", http.FileServer(http.Dir(cfg.UploadDir))))

	r.Route("/api/v1", func(r chi.Router) {

		r.Get("/theme", publicHandler.GetTheme)
		r.Get("/navigation", publicHandler.GetNavigation)
		r.Get("/pages/{slug}", publicHandler.GetPage)
		r.Get("/pages/{slug}/components", publicHandler.GetPageComponents)
		r.Get("/blogs", publicHandler.GetBlogs)
		r.Get("/blogs/{slug}", publicHandler.GetBlogBySlug)
		r.Get("/blogs/{slug}/related", publicHandler.GetRelatedBlogs)
		r.Get("/portfolios", publicHandler.GetPortfolios)
		r.Get("/portfolios/featured", publicHandler.GetFeaturedPortfolio)
		r.Get("/portfolios/{slug}", publicHandler.GetPortfolioBySlug)
		r.Get("/team", publicHandler.GetTeam)
		r.Get("/services", publicHandler.GetServices)
		r.Get("/faqs", publicHandler.GetFAQs)
		r.Get("/pricing", publicHandler.GetPricing)
		r.Get("/testimonials", publicHandler.GetTestimonials)
		r.Get("/awards", publicHandler.GetAwards)
		r.Get("/values", publicHandler.GetValues)
		r.Get("/job-openings", publicHandler.GetJobOpenings)
		r.Get("/social-links", publicHandler.GetSocialLinks)
		r.Get("/seo", publicHandler.GetSEO)
		r.Post("/contact", publicHandler.SubmitContact)

		r.Route("/admin", func(r chi.Router) {
			// Public admin routes
			r.Post("/auth/login", authHandler.Login)
			r.Post("/auth/refresh", authHandler.RefreshToken)
			r.Post("/auth/forgot-password", authHandler.ForgotPassword)
			r.Post("/auth/reset-password", authHandler.ResetPassword)

			// Protected admin routes
			r.Group(func(r chi.Router) {
				r.Use(authMW.Auth)

				r.Get("/auth/me", authHandler.GetMe)
				r.Put("/auth/profile", authHandler.UpdateProfile)
				r.Put("/auth/password", authHandler.ChangePassword)

				r.Get("/users", userHandler.GetAllUsers)
				r.Post("/users", userHandler.CreateUser)
				r.Get("/users/{id}", userHandler.GetUser)
				r.Put("/users/{id}", userHandler.UpdateUser)
				r.Delete("/users/{id}", userHandler.DeleteUser)
				r.Put("/users/{id}/role", userHandler.UpdateUserRole)

				r.Get("/dashboard/stats", adminHandler.GetDashboardStats)
				r.Get("/dashboard/recent-contacts", adminHandler.GetRecentContacts)

				r.Get("/theme", adminHandler.GetTheme)
				r.Put("/theme", adminHandler.UpdateTheme)
				r.Post("/theme/preview", adminHandler.PreviewTheme)
				r.Post("/theme/reset", adminHandler.ResetTheme)
				r.Get("/theme/presets", adminHandler.GetThemePresets)
				r.Post("/theme/preset", adminHandler.ApplyThemePreset)

				r.Get("/pages", adminHandler.GetAllPages)
				r.Get("/pages/{id}", adminHandler.GetPage)
				r.Put("/pages/{id}", adminHandler.UpdatePage)
				r.Put("/pages/{id}/sections", adminHandler.UpdatePageSections)

				r.Get("/components", adminHandler.GetAllComponents)
				r.Get("/components/{id}", adminHandler.GetComponent)
				r.Put("/components/{id}", adminHandler.UpdateComponent)
				r.Put("/components/{id}/visibility", adminHandler.ToggleComponentVisibility)
				r.Put("/components/reorder", adminHandler.ReorderComponents)

				r.Get("/navigation", adminHandler.GetAllNav)
				r.Post("/navigation", adminHandler.CreateNav)
				r.Put("/navigation/{id}", adminHandler.UpdateNav)
				r.Delete("/navigation/{id}", adminHandler.DeleteNav)
				r.Put("/navigation/reorder", adminHandler.ReorderNav)

				r.Get("/blogs", adminHandler.GetAllBlogs)
				r.Post("/blogs", adminHandler.CreateBlog)
				r.Get("/blogs/{id}", adminHandler.GetBlog)
				r.Put("/blogs/{id}", adminHandler.UpdateBlog)
				r.Delete("/blogs/{id}", adminHandler.DeleteBlog)

				r.Get("/portfolios", adminHandler.GetAllPortfolios)
				r.Post("/portfolios", adminHandler.CreatePortfolio)
				r.Get("/portfolios/{id}", adminHandler.GetPortfolio)
				r.Put("/portfolios/{id}", adminHandler.UpdatePortfolio)
				r.Delete("/portfolios/{id}", adminHandler.DeletePortfolio)

				r.Get("/team", adminHandler.GetAllTeam)
				r.Post("/team", adminHandler.CreateTeam)
				r.Get("/team/{id}", adminHandler.GetTeam)
				r.Put("/team/{id}", adminHandler.UpdateTeam)
				r.Delete("/team/{id}", adminHandler.DeleteTeam)

				r.Get("/awards", adminHandler.GetAllAwards)
				r.Post("/awards", adminHandler.CreateAward)
				r.Get("/awards/{id}", adminHandler.GetAward)
				r.Put("/awards/{id}", adminHandler.UpdateAward)
				r.Delete("/awards/{id}", adminHandler.DeleteAward)

				r.Get("/faqs", adminHandler.GetAllFAQs)
				r.Post("/faqs", adminHandler.CreateFAQ)
				r.Get("/faqs/{id}", adminHandler.GetFAQ)
				r.Put("/faqs/{id}", adminHandler.UpdateFAQ)
				r.Delete("/faqs/{id}", adminHandler.DeleteFAQ)

				r.Get("/pricing", adminHandler.GetAllPricing)
				r.Post("/pricing", adminHandler.CreatePricing)
				r.Get("/pricing/{id}", adminHandler.GetPricing)
				r.Put("/pricing/{id}", adminHandler.UpdatePricing)
				r.Delete("/pricing/{id}", adminHandler.DeletePricing)

				r.Get("/testimonials", adminHandler.GetAllTestimonials)
				r.Post("/testimonials", adminHandler.CreateTestimonial)
				r.Get("/testimonials/{id}", adminHandler.GetTestimonial)
				r.Put("/testimonials/{id}", adminHandler.UpdateTestimonial)
				r.Delete("/testimonials/{id}", adminHandler.DeleteTestimonial)

				r.Get("/services", adminHandler.GetAllServices)
				r.Post("/services", adminHandler.CreateService)
				r.Get("/services/{id}", adminHandler.GetService)
				r.Put("/services/{id}", adminHandler.UpdateService)
				r.Delete("/services/{id}", adminHandler.DeleteService)

				r.Get("/contacts", adminHandler.GetAllContacts)
				r.Get("/contacts/{id}", adminHandler.GetContact)
				r.Put("/contacts/{id}", adminHandler.UpdateContact)
				r.Delete("/contacts/{id}", adminHandler.DeleteContact)

				r.Get("/seo", adminHandler.GetSEO)
				r.Put("/seo", adminHandler.UpdateSEO)

				r.Get("/social-links", adminHandler.GetAllSocial)
				r.Post("/social-links", adminHandler.CreateSocial)
				r.Put("/social-links/{id}", adminHandler.UpdateSocial)
				r.Delete("/social-links/{id}", adminHandler.DeleteSocial)

				r.Get("/values", adminHandler.GetAllValues)
				r.Post("/values", adminHandler.CreateValue)
				r.Get("/values/{id}", adminHandler.GetValue)
				r.Put("/values/{id}", adminHandler.UpdateValue)
				r.Delete("/values/{id}", adminHandler.DeleteValue)

				r.Get("/job-openings", adminHandler.GetAllJobOpenings)
				r.Post("/job-openings", adminHandler.CreateJobOpening)
				r.Get("/job-openings/{id}", adminHandler.GetJobOpening)
				r.Put("/job-openings/{id}", adminHandler.UpdateJobOpening)
				r.Delete("/job-openings/{id}", adminHandler.DeleteJobOpening)

				r.Get("/media", mediaHandler.GetMedia)
				r.Post("/media/upload", mediaHandler.UploadMedia)
				r.Put("/media/{id}", mediaHandler.UpdateMedia)
				r.Delete("/media/{id}", mediaHandler.DeleteMedia)
			})
		})
	})

	log.Printf("Server starting on port %s", cfg.Port)

	go func() {
		if err := http.ListenAndServe(":"+cfg.Port, r); err != nil {
			log.Fatalf("Server failed: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down server...")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_ = ctx
	_ = cancel

	log.Println("Server exited")
}
