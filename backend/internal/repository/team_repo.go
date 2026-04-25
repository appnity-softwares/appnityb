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

type TeamRepository struct {
	db *pgxpool.Pool
}

func NewTeamRepo(db *pgxpool.Pool) *TeamRepository {
	return &TeamRepository{db: db}
}

func (r *TeamRepository) GetAll(ctx context.Context, visibleOnly bool) ([]models.Team, error) {
	var query string
	if visibleOnly {
		query = `SELECT id, COALESCE(full_name, ''), COALESCE(role, ''), COALESCE(bio, ''), COALESCE(photo, ''), COALESCE(social_links, '[]'), order_index, is_visible, created_at, updated_at
			FROM team WHERE is_visible = true ORDER BY order_index`
	} else {
		query = `SELECT id, COALESCE(full_name, ''), COALESCE(role, ''), COALESCE(bio, ''), COALESCE(photo, ''), COALESCE(social_links, '[]'), order_index, is_visible, created_at, updated_at
			FROM team ORDER BY order_index`
	}
	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var teams []models.Team
	for rows.Next() {
		var t models.Team
		if err := rows.Scan(&t.ID, &t.FullName, &t.Role, &t.Bio, &t.Photo, &t.SocialLinks, &t.OrderIndex, &t.IsVisible, &t.CreatedAt, &t.UpdatedAt); err != nil {
			return nil, err
		}
		teams = append(teams, t)
	}
	return teams, nil
}

func (r *TeamRepository) GetByID(ctx context.Context, id uuid.UUID) (*models.Team, error) {
	query := `SELECT id, COALESCE(full_name, ''), COALESCE(role, ''), COALESCE(bio, ''), COALESCE(photo, ''), COALESCE(social_links, '[]'), order_index, is_visible, created_at, updated_at
		FROM team WHERE id = $1`
	t := &models.Team{}
	err := r.db.QueryRow(ctx, query, id).Scan(&t.ID, &t.FullName, &t.Role, &t.Bio, &t.Photo, &t.SocialLinks, &t.OrderIndex, &t.IsVisible, &t.CreatedAt, &t.UpdatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("team member not found")
		}
		return nil, err
	}
	return t, nil
}

func (r *TeamRepository) Create(ctx context.Context, t *models.Team) error {
	t.ID = uuid.New()
	now := time.Now()
	t.CreatedAt = now
	t.UpdatedAt = now
	query := `INSERT INTO team (id, full_name, role, bio, photo, social_links, order_index, is_visible, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)`
	_, err := r.db.Exec(ctx, query, t.ID, t.FullName, t.Role, t.Bio, t.Photo, t.SocialLinks, t.OrderIndex, t.IsVisible, t.CreatedAt, t.UpdatedAt)
	return err
}

func (r *TeamRepository) Update(ctx context.Context, t *models.Team) error {
	t.UpdatedAt = time.Now()
	query := `UPDATE team SET full_name=$1, role=$2, bio=$3, photo=$4, social_links=$5, order_index=$6, is_visible=$7, updated_at=$8
		WHERE id=$9`
	result, err := r.db.Exec(ctx, query, t.FullName, t.Role, t.Bio, t.Photo, t.SocialLinks, t.OrderIndex, t.IsVisible, t.UpdatedAt, t.ID)
	if err != nil {
		return err
	}
	if result.RowsAffected() == 0 {
		return fmt.Errorf("team member not found")
	}
	return nil
}

func (r *TeamRepository) Delete(ctx context.Context, id uuid.UUID) error {
	result, err := r.db.Exec(ctx, `DELETE FROM team WHERE id = $1`, id)
	if err != nil {
		return err
	}
	if result.RowsAffected() == 0 {
		return fmt.Errorf("team member not found")
	}
	return nil
}

func (r *TeamRepository) Reorder(ctx context.Context, ids []uuid.UUID) error {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)
	for i, id := range ids {
		if _, err := tx.Exec(ctx, `UPDATE team SET order_index = $1 WHERE id = $2`, i, id); err != nil {
			return err
		}
	}
	return tx.Commit(ctx)
}
