package service

import (
	"context"

	"github.com/appnity/backend/internal/repository"
)

type ConfigService struct {
	contactRepo   *repository.ContactRepository
	blogRepo      *repository.BlogRepository
	portfolioRepo *repository.PortfolioRepository
	teamRepo      *repository.TeamRepository
}

func NewConfigService(
	contactRepo *repository.ContactRepository,
	blogRepo *repository.BlogRepository,
	portfolioRepo *repository.PortfolioRepository,
	teamRepo *repository.TeamRepository,
) *ConfigService {
	return &ConfigService{
		contactRepo:   contactRepo,
		blogRepo:      blogRepo,
		portfolioRepo: portfolioRepo,
		teamRepo:      teamRepo,
	}
}

func (s *ConfigService) GetDashboardStats(ctx context.Context) (map[string]interface{}, error) {
	stats := make(map[string]interface{})

	_, totalContacts, _ := s.contactRepo.GetAll(ctx, "", 1, 1)
	stats["total_contacts"] = totalContacts

	_, totalBlogs, _ := s.blogRepo.GetAll(ctx, false, 1, 1)
	stats["total_blogs"] = totalBlogs

	_, totalPortfolios, _ := s.portfolioRepo.GetAll(ctx, false, "", 1, 1)
	stats["total_portfolios"] = totalPortfolios

	team, _ := s.teamRepo.GetAll(ctx, false)
	stats["total_team_members"] = len(team)

	recentContacts, _, _ := s.contactRepo.GetAll(ctx, "", 1, 5)
	stats["recent_contacts"] = recentContacts

	contactStatusCounts, _ := s.contactRepo.CountByStatus(ctx)
	stats["contact_status_counts"] = contactStatusCounts

	return stats, nil
}
