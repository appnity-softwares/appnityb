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

type SocialRepository struct {
	db *pgxpool.Pool
}

func NewSocialRepo(db *pgxpool.Pool) *SocialRepository {
	return &SocialRepository{db: db}
}

func (r *SocialRepository) GetAll(ctx context.Context, visibleOnly bool) ([]models.Social, error) {
	var query string
	if visibleOnly {
		query = `SELECT id, platform, url, is_visible, COALESCE(icon, ''), order_index, created_at
			FROM social_links WHERE is_visible = true ORDER BY order_index`
	} else {
		query = `SELECT id, platform, url, is_visible, COALESCE(icon, ''), order_index, created_at
			FROM social_links ORDER BY order_index`
	}
	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var socials []models.Social
	for rows.Next() {
		var s models.Social
		if err := rows.Scan(&s.ID, &s.Platform, &s.URL, &s.IsVisible, &s.Icon, &s.OrderIndex, &s.CreatedAt); err != nil {
			return nil, err
		}
		socials = append(socials, s)
	}
	return socials, nil
}

func (r *SocialRepository) GetByID(ctx context.Context, id uuid.UUID) (*models.Social, error) {
	query := `SELECT id, platform, url, is_visible, COALESCE(icon, ''), order_index, created_at
		FROM social_links WHERE id = $1`
	s := &models.Social{}
	err := r.db.QueryRow(ctx, query, id).Scan(&s.ID, &s.Platform, &s.URL, &s.IsVisible, &s.Icon, &s.OrderIndex, &s.CreatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("social not found")
		}
		return nil, err
	}
	return s, nil
}

func (r *SocialRepository) Create(ctx context.Context, s *models.Social) error {
	s.ID = uuid.New()
	s.CreatedAt = time.Now()
	query := `INSERT INTO social_links (id, platform, url, is_visible, icon, order_index, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)`
	_, err := r.db.Exec(ctx, query, s.ID, s.Platform, s.URL, s.IsVisible, s.Icon, s.OrderIndex, s.CreatedAt)
	return err
}

func (r *SocialRepository) Update(ctx context.Context, s *models.Social) error {
	query := `UPDATE social_links SET platform=$1, url=$2, is_visible=$3, icon=$4, order_index=$5
		WHERE id=$6`
	result, err := r.db.Exec(ctx, query, s.Platform, s.URL, s.IsVisible, s.Icon, s.OrderIndex, s.ID)
	if err != nil {
		return err
	}
	if result.RowsAffected() == 0 {
		return fmt.Errorf("social not found")
	}
	return nil
}

func (r *SocialRepository) Delete(ctx context.Context, id uuid.UUID) error {
	result, err := r.db.Exec(ctx, `DELETE FROM social_links WHERE id = $1`, id)
	if err != nil {
		return err
	}
	if result.RowsAffected() == 0 {
		return fmt.Errorf("social not found")
	}
	return nil
}

func (r *SocialRepository) Reorder(ctx context.Context, ids []uuid.UUID) error {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)
	for i, id := range ids {
		if _, err := tx.Exec(ctx, `UPDATE social_links SET order_index = $1 WHERE id = $2`, i, id); err != nil {
			return err
		}
	}
	return tx.Commit(ctx)
}
