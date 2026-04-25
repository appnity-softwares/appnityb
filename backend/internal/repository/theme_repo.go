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

type ThemeRepository struct {
	db *pgxpool.Pool
}

func NewThemeRepo(db *pgxpool.Pool) *ThemeRepository {
	return &ThemeRepository{db: db}
}

func (r *ThemeRepository) GetActive(ctx context.Context) (*models.Theme, error) {
	query := `SELECT id, name, is_active, colors, fonts, typography, border_radius, spacing, shadows, breakpoints, background_images, gradients, animations, created_at, updated_at
		FROM themes WHERE is_active = true LIMIT 1`
	t := &models.Theme{}
	err := r.db.QueryRow(ctx, query).Scan(
		&t.ID, &t.Name, &t.IsActive, &t.Colors, &t.Fonts, &t.Typography,
		&t.BorderRadius, &t.Spacing, &t.Shadows, &t.Breakpoints, &t.BackgroundImages, &t.Gradients, &t.Animations,
		&t.CreatedAt, &t.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("active theme not found")
		}
		return nil, err
	}
	return t, nil
}

func (r *ThemeRepository) GetAll(ctx context.Context) ([]models.Theme, error) {
	query := `SELECT id, name, is_active, colors, fonts, typography, border_radius, spacing, shadows, breakpoints, background_images, gradients, animations, created_at, updated_at
		FROM themes ORDER BY created_at DESC`
	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var themes []models.Theme
	for rows.Next() {
		var t models.Theme
		if err := rows.Scan(
			&t.ID, &t.Name, &t.IsActive, &t.Colors, &t.Fonts, &t.Typography,
			&t.BorderRadius, &t.Spacing, &t.Shadows, &t.Breakpoints, &t.BackgroundImages, &t.Gradients, &t.Animations,
			&t.CreatedAt, &t.UpdatedAt,
		); err != nil {
			return nil, err
		}
		themes = append(themes, t)
	}
	return themes, nil
}

func (r *ThemeRepository) Create(ctx context.Context, t *models.Theme) error {
	t.ID = uuid.New()
	now := time.Now()
	t.CreatedAt = now
	t.UpdatedAt = now
	query := `INSERT INTO themes (id, name, is_active, colors, fonts, typography, border_radius, spacing, shadows, breakpoints, background_images, gradients, animations, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15)`
	_, err := r.db.Exec(ctx, query, t.ID, t.Name, t.IsActive, t.Colors, t.Fonts, t.Typography,
		t.BorderRadius, t.Spacing, t.Shadows, t.Breakpoints, t.BackgroundImages, t.Gradients, t.Animations, t.CreatedAt, t.UpdatedAt)
	return err
}

func (r *ThemeRepository) Update(ctx context.Context, t *models.Theme) error {
	t.UpdatedAt = time.Now()
	query := `UPDATE themes SET name=$1, is_active=$2, colors=$3, fonts=$4, typography=$5, border_radius=$6, spacing=$7, shadows=$8, breakpoints=$9, background_images=$10, gradients=$11, animations=$12, updated_at=$13
		WHERE id=$14`
	result, err := r.db.Exec(ctx, query, t.Name, t.IsActive, t.Colors, t.Fonts, t.Typography,
		t.BorderRadius, t.Spacing, t.Shadows, t.Breakpoints, t.BackgroundImages, t.Gradients, t.Animations, t.UpdatedAt, t.ID)
	if err != nil {
		return err
	}
	if result.RowsAffected() == 0 {
		return fmt.Errorf("theme not found")
	}
	return nil
}

func (r *ThemeRepository) SetActive(ctx context.Context, id uuid.UUID) error {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)
	if _, err := tx.Exec(ctx, `UPDATE themes SET is_active = false`); err != nil {
		return err
	}
	result, err := tx.Exec(ctx, `UPDATE themes SET is_active = true WHERE id = $1`, id)
	if err != nil {
		return err
	}
	if result.RowsAffected() == 0 {
		return fmt.Errorf("theme not found")
	}
	return tx.Commit(ctx)
}

func (r *ThemeRepository) Delete(ctx context.Context, id uuid.UUID) error {
	result, err := r.db.Exec(ctx, `DELETE FROM themes WHERE id = $1`, id)
	if err != nil {
		return err
	}
	if result.RowsAffected() == 0 {
		return fmt.Errorf("theme not found")
	}
	return nil
}
