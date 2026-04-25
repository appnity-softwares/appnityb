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

type AwardRepository struct {
	db *pgxpool.Pool
}

func NewAwardRepo(db *pgxpool.Pool) *AwardRepository {
	return &AwardRepository{db: db}
}

func (r *AwardRepository) GetAll(ctx context.Context, visibleOnly bool) ([]models.Award, error) {
	var query string
	if visibleOnly {
		query = `SELECT id, title, description, image, year, order_index, is_visible, created_at
			FROM awards WHERE is_visible = true ORDER BY order_index`
	} else {
		query = `SELECT id, title, description, image, year, order_index, is_visible, created_at
			FROM awards ORDER BY order_index`
	}
	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var awards []models.Award
	for rows.Next() {
		var a models.Award
		if err := rows.Scan(&a.ID, &a.Title, &a.Description, &a.Image, &a.Year, &a.OrderIndex, &a.IsVisible, &a.CreatedAt); err != nil {
			return nil, err
		}
		awards = append(awards, a)
	}
	return awards, nil
}

func (r *AwardRepository) GetByID(ctx context.Context, id uuid.UUID) (*models.Award, error) {
	query := `SELECT id, title, description, image, year, order_index, is_visible, created_at
		FROM awards WHERE id = $1`
	a := &models.Award{}
	err := r.db.QueryRow(ctx, query, id).Scan(&a.ID, &a.Title, &a.Description, &a.Image, &a.Year, &a.OrderIndex, &a.IsVisible, &a.CreatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("award not found")
		}
		return nil, err
	}
	return a, nil
}

func (r *AwardRepository) Create(ctx context.Context, a *models.Award) error {
	a.ID = uuid.New()
	a.CreatedAt = time.Now()
	query := `INSERT INTO awards (id, title, description, image, year, order_index, is_visible, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`
	_, err := r.db.Exec(ctx, query, a.ID, a.Title, a.Description, a.Image, a.Year, a.OrderIndex, a.IsVisible, a.CreatedAt)
	return err
}

func (r *AwardRepository) Update(ctx context.Context, a *models.Award) error {
	query := `UPDATE awards SET title=$1, description=$2, image=$3, year=$4, order_index=$5, is_visible=$6
		WHERE id=$7`
	result, err := r.db.Exec(ctx, query, a.Title, a.Description, a.Image, a.Year, a.OrderIndex, a.IsVisible, a.ID)
	if err != nil {
		return err
	}
	if result.RowsAffected() == 0 {
		return fmt.Errorf("award not found")
	}
	return nil
}

func (r *AwardRepository) Delete(ctx context.Context, id uuid.UUID) error {
	result, err := r.db.Exec(ctx, `DELETE FROM awards WHERE id = $1`, id)
	if err != nil {
		return err
	}
	if result.RowsAffected() == 0 {
		return fmt.Errorf("award not found")
	}
	return nil
}

func (r *AwardRepository) Reorder(ctx context.Context, ids []uuid.UUID) error {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)
	for i, id := range ids {
		if _, err := tx.Exec(ctx, `UPDATE awards SET order_index = $1 WHERE id = $2`, i, id); err != nil {
			return err
		}
	}
	return tx.Commit(ctx)
}
