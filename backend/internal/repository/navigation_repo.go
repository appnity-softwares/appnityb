package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/appnity/backend/internal/models"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type NavigationRepository struct {
	db *pgxpool.Pool
}

func NewNavigationRepo(db *pgxpool.Pool) *NavigationRepository {
	return &NavigationRepository{db: db}
}

func (r *NavigationRepository) GetAll(ctx context.Context, visibleOnly bool) ([]models.Navigation, error) {
	var query string
	if visibleOnly {
		query = `SELECT id, label, url, parent_id, COALESCE(order_index, 0), is_visible, is_external, COALESCE(icon, ''), created_at
			FROM navigation WHERE is_visible = true ORDER BY order_index`
	} else {
		query = `SELECT id, label, url, parent_id, COALESCE(order_index, 0), is_visible, is_external, COALESCE(icon, ''), created_at
			FROM navigation ORDER BY order_index`
	}
	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var navigations []models.Navigation
	for rows.Next() {
		var n models.Navigation
		if err := rows.Scan(&n.ID, &n.Label, &n.URL, &n.ParentID, &n.OrderIndex, &n.IsVisible, &n.IsExternal, &n.Icon, &n.CreatedAt); err != nil {
			return nil, err
		}
		navigations = append(navigations, n)
	}
	return navigations, nil
}

func (r *NavigationRepository) Create(ctx context.Context, n *models.Navigation) error {
	n.ID = uuid.New()
	n.CreatedAt = time.Now()
	query := `INSERT INTO navigation (id, label, url, parent_id, order_index, is_visible, is_external, icon, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)`
	_, err := r.db.Exec(ctx, query, n.ID, n.Label, n.URL, n.ParentID, n.OrderIndex, n.IsVisible, n.IsExternal, n.Icon, n.CreatedAt)
	return err
}

func (r *NavigationRepository) Update(ctx context.Context, n *models.Navigation) error {
	query := `UPDATE navigation SET label=$1, url=$2, parent_id=$3, order_index=$4, is_visible=$5, is_external=$6, icon=$7
		WHERE id=$8`
	result, err := r.db.Exec(ctx, query, n.Label, n.URL, n.ParentID, n.OrderIndex, n.IsVisible, n.IsExternal, n.Icon, n.ID)
	if err != nil {
		return err
	}
	if result.RowsAffected() == 0 {
		return fmt.Errorf("navigation not found")
	}
	return nil
}

func (r *NavigationRepository) Delete(ctx context.Context, id uuid.UUID) error {
	result, err := r.db.Exec(ctx, `DELETE FROM navigation WHERE id = $1`, id)
	if err != nil {
		return err
	}
	if result.RowsAffected() == 0 {
		return fmt.Errorf("navigation not found")
	}
	return nil
}

func (r *NavigationRepository) Reorder(ctx context.Context, ids []uuid.UUID) error {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)
	for i, id := range ids {
		if _, err := tx.Exec(ctx, `UPDATE navigation SET order_index = $1 WHERE id = $2`, i, id); err != nil {
			return err
		}
	}
	return tx.Commit(ctx)
}
