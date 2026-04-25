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

type ComponentRepository struct {
	db *pgxpool.Pool
}

func NewComponentRepo(db *pgxpool.Pool) *ComponentRepository {
	return &ComponentRepository{db: db}
}

func (r *ComponentRepository) GetAll(ctx context.Context) ([]models.Component, error) {
	query := `SELECT id, page_id, component_type, is_visible, content, order_index, created_at, updated_at
		FROM components ORDER BY page_id, order_index`
	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var components []models.Component
	for rows.Next() {
		var c models.Component
		if err := rows.Scan(&c.ID, &c.PageID, &c.ComponentType, &c.IsVisible, &c.Content, &c.OrderIndex, &c.CreatedAt, &c.UpdatedAt); err != nil {
			return nil, err
		}
		components = append(components, c)
	}
	return components, nil
}

func (r *ComponentRepository) GetByPageID(ctx context.Context, pageID uuid.UUID) ([]models.Component, error) {
	query := `SELECT id, page_id, component_type, is_visible, content, order_index, created_at, updated_at
		FROM components WHERE page_id = $1 ORDER BY order_index`
	rows, err := r.db.Query(ctx, query, pageID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var components []models.Component
	for rows.Next() {
		var c models.Component
		if err := rows.Scan(&c.ID, &c.PageID, &c.ComponentType, &c.IsVisible, &c.Content, &c.OrderIndex, &c.CreatedAt, &c.UpdatedAt); err != nil {
			return nil, err
		}
		components = append(components, c)
	}
	return components, nil
}

func (r *ComponentRepository) GetByID(ctx context.Context, id uuid.UUID) (*models.Component, error) {
	query := `SELECT id, page_id, component_type, is_visible, content, order_index, created_at, updated_at
		FROM components WHERE id = $1`
	c := &models.Component{}
	err := r.db.QueryRow(ctx, query, id).Scan(&c.ID, &c.PageID, &c.ComponentType, &c.IsVisible, &c.Content, &c.OrderIndex, &c.CreatedAt, &c.UpdatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("component not found")
		}
		return nil, err
	}
	return c, nil
}

func (r *ComponentRepository) Create(ctx context.Context, c *models.Component) error {
	c.ID = uuid.New()
	now := time.Now()
	c.CreatedAt = now
	c.UpdatedAt = now
	query := `INSERT INTO components (id, page_id, component_type, is_visible, content, order_index, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`
	_, err := r.db.Exec(ctx, query, c.ID, c.PageID, c.ComponentType, c.IsVisible, c.Content, c.OrderIndex, c.CreatedAt, c.UpdatedAt)
	return err
}

func (r *ComponentRepository) Update(ctx context.Context, c *models.Component) error {
	c.UpdatedAt = time.Now()
	query := `UPDATE components SET page_id=$1, component_type=$2, is_visible=$3, content=$4, order_index=$5, updated_at=$6
		WHERE id=$7`
	result, err := r.db.Exec(ctx, query, c.PageID, c.ComponentType, c.IsVisible, c.Content, c.OrderIndex, c.UpdatedAt, c.ID)
	if err != nil {
		return err
	}
	if result.RowsAffected() == 0 {
		return fmt.Errorf("component not found")
	}
	return nil
}

func (r *ComponentRepository) Delete(ctx context.Context, id uuid.UUID) error {
	result, err := r.db.Exec(ctx, `DELETE FROM components WHERE id = $1`, id)
	if err != nil {
		return err
	}
	if result.RowsAffected() == 0 {
		return fmt.Errorf("component not found")
	}
	return nil
}

func (r *ComponentRepository) Reorder(ctx context.Context, pageID uuid.UUID, ids []uuid.UUID) error {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)
	for i, id := range ids {
		result, err := tx.Exec(ctx, `UPDATE components SET order_index = $1 WHERE id = $2 AND page_id = $3`, i, id, pageID)
		if err != nil {
			return err
		}
		if result.RowsAffected() == 0 {
			return fmt.Errorf("component %s not found on page %s", id, pageID)
		}
	}
	return tx.Commit(ctx)
}

func (r *ComponentRepository) ToggleVisibility(ctx context.Context, id uuid.UUID) error {
	query := `UPDATE components SET is_visible = NOT is_visible, updated_at = $2 WHERE id = $1`
	result, err := r.db.Exec(ctx, query, id, time.Now())
	if err != nil {
		return err
	}
	if result.RowsAffected() == 0 {
		return fmt.Errorf("component not found")
	}
	return nil
}
