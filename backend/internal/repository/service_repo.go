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

type ServiceRepository struct {
	db *pgxpool.Pool
}

func NewServiceRepo(db *pgxpool.Pool) *ServiceRepository {
	return &ServiceRepository{db: db}
}

func (r *ServiceRepository) GetAll(ctx context.Context, visibleOnly bool) ([]models.Service, error) {
	var query string
	if visibleOnly {
		query = `SELECT id, title, COALESCE(description, ''), COALESCE(icon, ''), COALESCE(image, ''), COALESCE(features, '[]'), COALESCE(order_index, 0), COALESCE(is_visible, true), created_at
			FROM services WHERE is_visible = true ORDER BY order_index`
	} else {
		query = `SELECT id, title, COALESCE(description, ''), COALESCE(icon, ''), COALESCE(image, ''), COALESCE(features, '[]'), COALESCE(order_index, 0), COALESCE(is_visible, true), created_at
			FROM services ORDER BY order_index`
	}
	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var services []models.Service
	for rows.Next() {
		var s models.Service
		if err := rows.Scan(&s.ID, &s.Title, &s.Description, &s.Icon, &s.Image, &s.Features, &s.OrderIndex, &s.IsVisible, &s.CreatedAt); err != nil {
			return nil, err
		}
		services = append(services, s)
	}
	return services, nil
}

func (r *ServiceRepository) GetByID(ctx context.Context, id uuid.UUID) (*models.Service, error) {
	query := `SELECT id, title, COALESCE(description, ''), COALESCE(icon, ''), COALESCE(image, ''), COALESCE(features, '[]'), COALESCE(order_index, 0), COALESCE(is_visible, true), created_at
		FROM services WHERE id = $1`
	s := &models.Service{}
	err := r.db.QueryRow(ctx, query, id).Scan(&s.ID, &s.Title, &s.Description, &s.Icon, &s.Image, &s.Features, &s.OrderIndex, &s.IsVisible, &s.CreatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("service not found")
		}
		return nil, err
	}
	return s, nil
}

func (r *ServiceRepository) Create(ctx context.Context, s *models.Service) error {
	s.ID = uuid.New()
	s.CreatedAt = time.Now()
	query := `INSERT INTO services (id, title, description, icon, image, features, order_index, is_visible, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)`
	_, err := r.db.Exec(ctx, query, s.ID, s.Title, s.Description, s.Icon, s.Image, s.Features, s.OrderIndex, s.IsVisible, s.CreatedAt)
	return err
}

func (r *ServiceRepository) Update(ctx context.Context, s *models.Service) error {
	query := `UPDATE services SET title=$1, description=$2, icon=$3, image=$4, features=$5, order_index=$6, is_visible=$7
		WHERE id=$8`
	result, err := r.db.Exec(ctx, query, s.Title, s.Description, s.Icon, s.Image, s.Features, s.OrderIndex, s.IsVisible, s.ID)
	if err != nil {
		return err
	}
	if result.RowsAffected() == 0 {
		return fmt.Errorf("service not found")
	}
	return nil
}

func (r *ServiceRepository) Delete(ctx context.Context, id uuid.UUID) error {
	result, err := r.db.Exec(ctx, `DELETE FROM services WHERE id = $1`, id)
	if err != nil {
		return err
	}
	if result.RowsAffected() == 0 {
		return fmt.Errorf("service not found")
	}
	return nil
}

func (r *ServiceRepository) Reorder(ctx context.Context, ids []uuid.UUID) error {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)
	for i, id := range ids {
		if _, err := tx.Exec(ctx, `UPDATE services SET order_index = $1 WHERE id = $2`, i, id); err != nil {
			return err
		}
	}
	return tx.Commit(ctx)
}
