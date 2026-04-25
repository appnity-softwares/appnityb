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

type ValueRepository struct {
	db *pgxpool.Pool
}

func NewValueRepo(db *pgxpool.Pool) *ValueRepository {
	return &ValueRepository{db: db}
}

func (r *ValueRepository) GetAll(ctx context.Context, visibleOnly bool) ([]models.Value, error) {
	var query string
	if visibleOnly {
		query = `SELECT id, title, description, icon, order_index, is_visible, created_at
			FROM "values" WHERE is_visible = true ORDER BY order_index`
	} else {
		query = `SELECT id, title, description, icon, order_index, is_visible, created_at
			FROM "values" ORDER BY order_index`
	}
	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var values []models.Value
	for rows.Next() {
		var v models.Value
		if err := rows.Scan(&v.ID, &v.Title, &v.Description, &v.Icon, &v.OrderIndex, &v.IsVisible, &v.CreatedAt); err != nil {
			return nil, err
		}
		values = append(values, v)
	}
	return values, nil
}

func (r *ValueRepository) GetByID(ctx context.Context, id uuid.UUID) (*models.Value, error) {
	query := `SELECT id, title, description, icon, order_index, is_visible, created_at
		FROM "values" WHERE id = $1`
	v := &models.Value{}
	err := r.db.QueryRow(ctx, query, id).Scan(&v.ID, &v.Title, &v.Description, &v.Icon, &v.OrderIndex, &v.IsVisible, &v.CreatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("value not found")
		}
		return nil, err
	}
	return v, nil
}

func (r *ValueRepository) Create(ctx context.Context, v *models.Value) error {
	v.ID = uuid.New()
	v.CreatedAt = time.Now()
	query := `INSERT INTO "values" (id, title, description, icon, order_index, is_visible, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)`
	_, err := r.db.Exec(ctx, query, v.ID, v.Title, v.Description, v.Icon, v.OrderIndex, v.IsVisible, v.CreatedAt)
	return err
}

func (r *ValueRepository) Update(ctx context.Context, v *models.Value) error {
	query := `UPDATE "values" SET title=$1, description=$2, icon=$3, order_index=$4, is_visible=$5
		WHERE id=$6`
	result, err := r.db.Exec(ctx, query, v.Title, v.Description, v.Icon, v.OrderIndex, v.IsVisible, v.ID)
	if err != nil {
		return err
	}
	if result.RowsAffected() == 0 {
		return fmt.Errorf("value not found")
	}
	return nil
}

func (r *ValueRepository) Delete(ctx context.Context, id uuid.UUID) error {
	result, err := r.db.Exec(ctx, `DELETE FROM "values" WHERE id = $1`, id)
	if err != nil {
		return err
	}
	if result.RowsAffected() == 0 {
		return fmt.Errorf("value not found")
	}
	return nil
}

func (r *ValueRepository) Reorder(ctx context.Context, ids []uuid.UUID) error {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)
	for i, id := range ids {
		if _, err := tx.Exec(ctx, `UPDATE "values" SET order_index = $1 WHERE id = $2`, i, id); err != nil {
			return err
		}
	}
	return tx.Commit(ctx)
}
