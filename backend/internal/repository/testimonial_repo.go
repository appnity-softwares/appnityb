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

type TestimonialRepository struct {
	db *pgxpool.Pool
}

func NewTestimonialRepo(db *pgxpool.Pool) *TestimonialRepository {
	return &TestimonialRepository{db: db}
}

func (r *TestimonialRepository) GetAll(ctx context.Context, visibleOnly bool) ([]models.Testimonial, error) {
	var query string
	if visibleOnly {
		query = `SELECT id, name, COALESCE(role, ''), COALESCE(company, ''), COALESCE(avatar, ''), COALESCE(quote, ''), COALESCE(rating, 5), COALESCE(video_url, ''), COALESCE(order_index, 0), COALESCE(is_visible, true), created_at
			FROM testimonials WHERE is_visible = true ORDER BY order_index`
	} else {
		query = `SELECT id, name, COALESCE(role, ''), COALESCE(company, ''), COALESCE(avatar, ''), COALESCE(quote, ''), COALESCE(rating, 5), COALESCE(video_url, ''), COALESCE(order_index, 0), COALESCE(is_visible, true), created_at
			FROM testimonials ORDER BY order_index`
	}
	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var testimonials []models.Testimonial
	for rows.Next() {
		var t models.Testimonial
		if err := rows.Scan(&t.ID, &t.Name, &t.Role, &t.Company, &t.Avatar, &t.Quote, &t.Rating, &t.VideoURL, &t.OrderIndex, &t.IsVisible, &t.CreatedAt); err != nil {
			return nil, err
		}
		testimonials = append(testimonials, t)
	}
	return testimonials, nil
}

func (r *TestimonialRepository) GetByID(ctx context.Context, id uuid.UUID) (*models.Testimonial, error) {
	query := `SELECT id, name, COALESCE(role, ''), COALESCE(company, ''), COALESCE(avatar, ''), COALESCE(quote, ''), COALESCE(rating, 5), COALESCE(video_url, ''), COALESCE(order_index, 0), COALESCE(is_visible, true), created_at
		FROM testimonials WHERE id = $1`
	t := &models.Testimonial{}
	err := r.db.QueryRow(ctx, query, id).Scan(&t.ID, &t.Name, &t.Role, &t.Company, &t.Avatar, &t.Quote, &t.Rating, &t.VideoURL, &t.OrderIndex, &t.IsVisible, &t.CreatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("testimonial not found")
		}
		return nil, err
	}
	return t, nil
}

func (r *TestimonialRepository) Create(ctx context.Context, t *models.Testimonial) error {
	t.ID = uuid.New()
	t.CreatedAt = time.Now()
	query := `INSERT INTO testimonials (id, name, role, company, avatar, quote, rating, video_url, order_index, is_visible, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)`
	_, err := r.db.Exec(ctx, query, t.ID, t.Name, t.Role, t.Company, t.Avatar, t.Quote, t.Rating, t.VideoURL, t.OrderIndex, t.IsVisible, t.CreatedAt)
	return err
}

func (r *TestimonialRepository) Update(ctx context.Context, t *models.Testimonial) error {
	query := `UPDATE testimonials SET name=$1, role=$2, company=$3, avatar=$4, quote=$5, rating=$6, video_url=$7, order_index=$8, is_visible=$9
		WHERE id=$10`
	result, err := r.db.Exec(ctx, query, t.Name, t.Role, t.Company, t.Avatar, t.Quote, t.Rating, t.VideoURL, t.OrderIndex, t.IsVisible, t.ID)
	if err != nil {
		return err
	}
	if result.RowsAffected() == 0 {
		return fmt.Errorf("testimonial not found")
	}
	return nil
}

func (r *TestimonialRepository) Delete(ctx context.Context, id uuid.UUID) error {
	result, err := r.db.Exec(ctx, `DELETE FROM testimonials WHERE id = $1`, id)
	if err != nil {
		return err
	}
	if result.RowsAffected() == 0 {
		return fmt.Errorf("testimonial not found")
	}
	return nil
}

func (r *TestimonialRepository) Reorder(ctx context.Context, ids []uuid.UUID) error {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)
	for i, id := range ids {
		if _, err := tx.Exec(ctx, `UPDATE testimonials SET order_index = $1 WHERE id = $2`, i, id); err != nil {
			return err
		}
	}
	return tx.Commit(ctx)
}
