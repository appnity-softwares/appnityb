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

type PortfolioRepository struct {
	db *pgxpool.Pool
}

func NewPortfolioRepo(db *pgxpool.Pool) *PortfolioRepository {
	return &PortfolioRepository{db: db}
}

func (r *PortfolioRepository) GetAll(ctx context.Context, visibleOnly bool, category string, page, pageSize int) ([]models.Portfolio, int64, error) {
	var countQuery string
	var dataQuery string
	var countArgs []interface{}
	var dataArgs []interface{}

	fields := `id, title, slug, COALESCE(description, ''), COALESCE(cover_image, ''), COALESCE(gallery_images, '[]'), category, client_name, project_url, COALESCE(metrics, '{}'), is_featured, is_visible, order_index, created_at, updated_at, COALESCE(industry, ''), COALESCE(year, 0), COALESCE(stack, ''), COALESCE(problem, ''), COALESCE(approach, '[]')`

	if visibleOnly && category != "" {
		countQuery = `SELECT COUNT(*) FROM portfolios WHERE is_visible = true AND category = $1`
		dataQuery = fmt.Sprintf(`SELECT %s FROM portfolios WHERE is_visible = true AND category = $3 ORDER BY order_index LIMIT $1 OFFSET $2`, fields)
		countArgs = append(countArgs, category)
		dataArgs = append(dataArgs, pageSize, (page-1)*pageSize, category)
	} else if visibleOnly {
		countQuery = `SELECT COUNT(*) FROM portfolios WHERE is_visible = true`
		dataQuery = fmt.Sprintf(`SELECT %s FROM portfolios WHERE is_visible = true ORDER BY order_index LIMIT $1 OFFSET $2`, fields)
		dataArgs = append(dataArgs, pageSize, (page-1)*pageSize)
	} else if category != "" {
		countQuery = `SELECT COUNT(*) FROM portfolios WHERE category = $1`
		dataQuery = fmt.Sprintf(`SELECT %s FROM portfolios WHERE category = $3 ORDER BY order_index LIMIT $1 OFFSET $2`, fields)
		countArgs = append(countArgs, category)
		dataArgs = append(dataArgs, pageSize, (page-1)*pageSize, category)
	} else {
		countQuery = `SELECT COUNT(*) FROM portfolios`
		dataQuery = fmt.Sprintf(`SELECT %s FROM portfolios ORDER BY order_index LIMIT $1 OFFSET $2`, fields)
		dataArgs = append(dataArgs, pageSize, (page-1)*pageSize)
	}

	var total int64
	if err := r.db.QueryRow(ctx, countQuery, countArgs...).Scan(&total); err != nil {
		return nil, 0, err
	}

	rows, err := r.db.Query(ctx, dataQuery, dataArgs...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var portfolios []models.Portfolio
	for rows.Next() {
		var p models.Portfolio
		if err := rows.Scan(&p.ID, &p.Title, &p.Slug, &p.Description, &p.CoverImage, &p.GalleryImages,
			&p.Category, &p.ClientName, &p.ProjectURL, &p.Metrics, &p.IsFeatured, &p.IsVisible,
			&p.OrderIndex, &p.CreatedAt, &p.UpdatedAt, &p.Industry, &p.Year, &p.Stack, &p.Problem, &p.Approach); err != nil {
			return nil, 0, err
		}
		p.Tagline = p.Description
		portfolios = append(portfolios, p)
	}
	return portfolios, total, nil
}

func (r *PortfolioRepository) GetBySlug(ctx context.Context, slug string) (*models.Portfolio, error) {
	fields := `id, title, slug, COALESCE(description, ''), COALESCE(cover_image, ''), COALESCE(gallery_images, '[]'), category, client_name, project_url, COALESCE(metrics, '{}'), is_featured, is_visible, order_index, created_at, updated_at, COALESCE(industry, ''), COALESCE(year, 0), COALESCE(stack, ''), COALESCE(problem, ''), COALESCE(approach, '[]')`
	query := fmt.Sprintf(`SELECT %s FROM portfolios WHERE slug = $1`, fields)
	p := &models.Portfolio{}
	err := r.db.QueryRow(ctx, query, slug).Scan(&p.ID, &p.Title, &p.Slug, &p.Description, &p.CoverImage,
		&p.GalleryImages, &p.Category, &p.ClientName, &p.ProjectURL, &p.Metrics, &p.IsFeatured, &p.IsVisible,
		&p.OrderIndex, &p.CreatedAt, &p.UpdatedAt, &p.Industry, &p.Year, &p.Stack, &p.Problem, &p.Approach)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("portfolio not found")
		}
		return nil, err
	}
	p.Tagline = p.Description
	return p, nil
}

func (r *PortfolioRepository) GetByID(ctx context.Context, id uuid.UUID) (*models.Portfolio, error) {
	fields := `id, title, slug, COALESCE(description, ''), COALESCE(cover_image, ''), COALESCE(gallery_images, '[]'), category, client_name, project_url, COALESCE(metrics, '{}'), is_featured, is_visible, order_index, created_at, updated_at, COALESCE(industry, ''), COALESCE(year, 0), COALESCE(stack, ''), COALESCE(problem, ''), COALESCE(approach, '[]')`
	query := fmt.Sprintf(`SELECT %s FROM portfolios WHERE id = $1`, fields)
	p := &models.Portfolio{}
	err := r.db.QueryRow(ctx, query, id).Scan(&p.ID, &p.Title, &p.Slug, &p.Description, &p.CoverImage,
		&p.GalleryImages, &p.Category, &p.ClientName, &p.ProjectURL, &p.Metrics, &p.IsFeatured, &p.IsVisible,
		&p.OrderIndex, &p.CreatedAt, &p.UpdatedAt, &p.Industry, &p.Year, &p.Stack, &p.Problem, &p.Approach)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("portfolio not found")
		}
		return nil, err
	}
	p.Tagline = p.Description
	return p, nil
}

func (r *PortfolioRepository) GetFeatured(ctx context.Context) ([]models.Portfolio, error) {
	fields := `id, title, slug, COALESCE(description, ''), COALESCE(cover_image, ''), COALESCE(gallery_images, '[]'), category, client_name, project_url, COALESCE(metrics, '{}'), is_featured, is_visible, order_index, created_at, updated_at, COALESCE(industry, ''), COALESCE(year, 0), COALESCE(stack, ''), COALESCE(problem, ''), COALESCE(approach, '[]')`
	query := fmt.Sprintf(`SELECT %s FROM portfolios WHERE is_featured = true ORDER BY order_index`, fields)
	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var portfolios []models.Portfolio
	for rows.Next() {
		var p models.Portfolio
		if err := rows.Scan(&p.ID, &p.Title, &p.Slug, &p.Description, &p.CoverImage, &p.GalleryImages,
			&p.Category, &p.ClientName, &p.ProjectURL, &p.Metrics, &p.IsFeatured, &p.IsVisible,
			&p.OrderIndex, &p.CreatedAt, &p.UpdatedAt, &p.Industry, &p.Year, &p.Stack, &p.Problem, &p.Approach); err != nil {
			return nil, err
		}
		p.Tagline = p.Description
		portfolios = append(portfolios, p)
	}
	return portfolios, nil
}

func (r *PortfolioRepository) Create(ctx context.Context, p *models.Portfolio) error {
	p.ID = uuid.New()
	now := time.Now()
	p.CreatedAt = now
	p.UpdatedAt = now
	query := `INSERT INTO portfolios (id, title, slug, description, cover_image, gallery_images, category, client_name, project_url, metrics, is_featured, is_visible, order_index, created_at, updated_at, industry, year, stack, problem, approach)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19, $20)`
	_, err := r.db.Exec(ctx, query, p.ID, p.Title, p.Slug, p.Description, p.CoverImage, p.GalleryImages,
		p.Category, p.ClientName, p.ProjectURL, p.Metrics, p.IsFeatured, p.IsVisible, p.OrderIndex, p.CreatedAt, p.UpdatedAt, p.Industry, p.Year, p.Stack, p.Problem, p.Approach)
	return err
}

func (r *PortfolioRepository) Update(ctx context.Context, p *models.Portfolio) error {
	p.UpdatedAt = time.Now()
	query := `UPDATE portfolios SET title=$1, slug=$2, description=$3, cover_image=$4, gallery_images=$5, category=$6, client_name=$7, project_url=$8, metrics=$9, is_featured=$10, is_visible=$11, order_index=$12, updated_at=$13, industry=$14, year=$15, stack=$16, problem=$17, approach=$18
		WHERE id=$19`
	result, err := r.db.Exec(ctx, query, p.Title, p.Slug, p.Description, p.CoverImage, p.GalleryImages,
		p.Category, p.ClientName, p.ProjectURL, p.Metrics, p.IsFeatured, p.IsVisible, p.OrderIndex, p.UpdatedAt, p.Industry, p.Year, p.Stack, p.Problem, p.Approach, p.ID)
	if err != nil {
		return err
	}
	if result.RowsAffected() == 0 {
		return fmt.Errorf("portfolio not found")
	}
	return nil
}

func (r *PortfolioRepository) Delete(ctx context.Context, id uuid.UUID) error {
	result, err := r.db.Exec(ctx, `DELETE FROM portfolios WHERE id = $1`, id)
	if err != nil {
		return err
	}
	if result.RowsAffected() == 0 {
		return fmt.Errorf("portfolio not found")
	}
	return nil
}

func (r *PortfolioRepository) Reorder(ctx context.Context, ids []uuid.UUID) error {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)
	for i, id := range ids {
		if _, err := tx.Exec(ctx, `UPDATE portfolios SET order_index = $1 WHERE id = $2`, i, id); err != nil {
			return err
		}
	}
	return tx.Commit(ctx)
}
