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

type PageRepository struct {
	db *pgxpool.Pool
}

func NewPageRepo(db *pgxpool.Pool) *PageRepository {
	return &PageRepository{db: db}
}

func (r *PageRepository) GetAll(ctx context.Context) ([]models.Page, error) {
	query := `SELECT id, slug, title, meta_description, meta_keywords, og_image, canonical_url, background_color, custom_css, is_published, sections, created_at, updated_at
		FROM pages ORDER BY created_at DESC`
	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var pages []models.Page
	for rows.Next() {
		var p models.Page
		if err := rows.Scan(
			&p.ID, &p.Slug, &p.Title, &p.MetaDescription, &p.MetaKeywords, &p.OgImage,
			&p.CanonicalURL, &p.BackgroundColor, &p.CustomCSS, &p.IsPublished, &p.Sections,
			&p.CreatedAt, &p.UpdatedAt,
		); err != nil {
			return nil, err
		}
		pages = append(pages, p)
	}
	return pages, nil
}

func (r *PageRepository) GetBySlug(ctx context.Context, slug string) (*models.Page, error) {
	query := `SELECT id, slug, title, meta_description, meta_keywords, og_image, canonical_url, background_color, custom_css, is_published, sections, created_at, updated_at
		FROM pages WHERE slug = $1`
	p := &models.Page{}
	err := r.db.QueryRow(ctx, query, slug).Scan(
		&p.ID, &p.Slug, &p.Title, &p.MetaDescription, &p.MetaKeywords, &p.OgImage,
		&p.CanonicalURL, &p.BackgroundColor, &p.CustomCSS, &p.IsPublished, &p.Sections,
		&p.CreatedAt, &p.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("page not found")
		}
		return nil, err
	}
	return p, nil
}

func (r *PageRepository) GetByID(ctx context.Context, id uuid.UUID) (*models.Page, error) {
	query := `SELECT id, slug, title, meta_description, meta_keywords, og_image, canonical_url, background_color, custom_css, is_published, sections, created_at, updated_at
		FROM pages WHERE id = $1`
	p := &models.Page{}
	err := r.db.QueryRow(ctx, query, id).Scan(
		&p.ID, &p.Slug, &p.Title, &p.MetaDescription, &p.MetaKeywords, &p.OgImage,
		&p.CanonicalURL, &p.BackgroundColor, &p.CustomCSS, &p.IsPublished, &p.Sections,
		&p.CreatedAt, &p.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("page not found")
		}
		return nil, err
	}
	return p, nil
}

func (r *PageRepository) Create(ctx context.Context, p *models.Page) error {
	p.ID = uuid.New()
	now := time.Now()
	p.CreatedAt = now
	p.UpdatedAt = now
	query := `INSERT INTO pages (id, slug, title, meta_description, meta_keywords, og_image, canonical_url, background_color, custom_css, is_published, sections, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13)`
	_, err := r.db.Exec(ctx, query, p.ID, p.Slug, p.Title, p.MetaDescription, p.MetaKeywords,
		p.OgImage, p.CanonicalURL, p.BackgroundColor, p.CustomCSS, p.IsPublished, p.Sections, p.CreatedAt, p.UpdatedAt)
	return err
}

func (r *PageRepository) Update(ctx context.Context, p *models.Page) error {
	p.UpdatedAt = time.Now()
	query := `UPDATE pages SET slug=$1, title=$2, meta_description=$3, meta_keywords=$4, og_image=$5, canonical_url=$6, background_color=$7, custom_css=$8, is_published=$9, sections=$10, updated_at=$11
		WHERE id=$12`
	result, err := r.db.Exec(ctx, query, p.Slug, p.Title, p.MetaDescription, p.MetaKeywords,
		p.OgImage, p.CanonicalURL, p.BackgroundColor, p.CustomCSS, p.IsPublished, p.Sections, p.UpdatedAt, p.ID)
	if err != nil {
		return err
	}
	if result.RowsAffected() == 0 {
		return fmt.Errorf("page not found")
	}
	return nil
}

func (r *PageRepository) Delete(ctx context.Context, id uuid.UUID) error {
	result, err := r.db.Exec(ctx, `DELETE FROM pages WHERE id = $1`, id)
	if err != nil {
		return err
	}
	if result.RowsAffected() == 0 {
		return fmt.Errorf("page not found")
	}
	return nil
}
