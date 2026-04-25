package repository

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/appnity/backend/internal/models"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type SEORepository struct {
	db *pgxpool.Pool
}

func NewSEORepo(db *pgxpool.Pool) *SEORepository {
	return &SEORepository{db: db}
}

func (r *SEORepository) Get(ctx context.Context) (*models.SEO, error) {
	query := `SELECT id, COALESCE(site_title, ''), COALESCE(site_description, ''), COALESCE(keywords, '{}'), COALESCE(og_image, ''), COALESCE(twitter_handle, ''), COALESCE(canonical_url, ''), COALESCE(robots_txt, ''), COALESCE(google_analytics_id, ''), COALESCE(google_tag_manager, ''), COALESCE(json_ld_schema, '{}'), COALESCE(favicon, ''), created_at, updated_at
		FROM seo_settings LIMIT 1`
	s := &models.SEO{}
	err := r.db.QueryRow(ctx, query).Scan(
		&s.ID, &s.SiteTitle, &s.SiteDescription, &s.Keywords, &s.OgImage, &s.TwitterHandle,
		&s.CanonicalURL, &s.RobotsTxt, &s.GoogleAnalyticsID, &s.GoogleTagManager, &s.JSONLDSchema,
		&s.Favicon, &s.CreatedAt, &s.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("seo settings not found")
		}
		return nil, err
	}
	return s, nil
}

func (r *SEORepository) Create(ctx context.Context, s *models.SEO) error {
	s.ID = uuid.New()
	now := time.Now()
	s.CreatedAt = now
	s.UpdatedAt = now
	query := `INSERT INTO seo_settings (id, site_title, site_description, keywords, og_image, twitter_handle, canonical_url, robots_txt, google_analytics_id, google_tag_manager, json_ld_schema, favicon, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14)`
	_, err := r.db.Exec(ctx, query, s.ID, s.SiteTitle, s.SiteDescription, s.Keywords, s.OgImage,
		s.TwitterHandle, s.CanonicalURL, s.RobotsTxt, s.GoogleAnalyticsID, s.GoogleTagManager,
		s.JSONLDSchema, s.Favicon, s.CreatedAt, s.UpdatedAt)
	return err
}

func (r *SEORepository) Update(ctx context.Context, s *models.SEO) error {
	s.UpdatedAt = time.Now()
	query := `UPDATE seo_settings SET site_title=$1, site_description=$2, keywords=$3, og_image=$4, twitter_handle=$5, canonical_url=$6, robots_txt=$7, google_analytics_id=$8, google_tag_manager=$9, json_ld_schema=$10, favicon=$11, updated_at=$12
		WHERE id=$13`
	result, err := r.db.Exec(ctx, query, s.SiteTitle, s.SiteDescription, s.Keywords, s.OgImage,
		s.TwitterHandle, s.CanonicalURL, s.RobotsTxt, s.GoogleAnalyticsID, s.GoogleTagManager,
		s.JSONLDSchema, s.Favicon, s.UpdatedAt, s.ID)
	if err != nil {
		return err
	}
	if result.RowsAffected() == 0 {
		return fmt.Errorf("seo settings not found")
	}
	return nil
}
