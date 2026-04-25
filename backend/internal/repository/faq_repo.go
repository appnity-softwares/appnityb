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

type FAQRepository struct {
	db *pgxpool.Pool
}

func NewFAQRepo(db *pgxpool.Pool) *FAQRepository {
	return &FAQRepository{db: db}
}

func (r *FAQRepository) GetAll(ctx context.Context, visibleOnly bool) ([]models.FAQ, error) {
	var query string
	if visibleOnly {
		query = `SELECT id, question, answer, order_index, is_visible, created_at
			FROM faqs WHERE is_visible = true ORDER BY order_index`
	} else {
		query = `SELECT id, question, answer, order_index, is_visible, created_at
			FROM faqs ORDER BY order_index`
	}
	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var faqs []models.FAQ
	for rows.Next() {
		var f models.FAQ
		if err := rows.Scan(&f.ID, &f.Question, &f.Answer, &f.OrderIndex, &f.IsVisible, &f.CreatedAt); err != nil {
			return nil, err
		}
		faqs = append(faqs, f)
	}
	return faqs, nil
}

func (r *FAQRepository) GetByID(ctx context.Context, id uuid.UUID) (*models.FAQ, error) {
	query := `SELECT id, question, answer, order_index, is_visible, created_at
		FROM faqs WHERE id = $1`
	f := &models.FAQ{}
	err := r.db.QueryRow(ctx, query, id).Scan(&f.ID, &f.Question, &f.Answer, &f.OrderIndex, &f.IsVisible, &f.CreatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("faq not found")
		}
		return nil, err
	}
	return f, nil
}

func (r *FAQRepository) Create(ctx context.Context, f *models.FAQ) error {
	f.ID = uuid.New()
	f.CreatedAt = time.Now()
	query := `INSERT INTO faqs (id, question, answer, order_index, is_visible, created_at)
		VALUES ($1, $2, $3, $4, $5, $6)`
	_, err := r.db.Exec(ctx, query, f.ID, f.Question, f.Answer, f.OrderIndex, f.IsVisible, f.CreatedAt)
	return err
}

func (r *FAQRepository) Update(ctx context.Context, f *models.FAQ) error {
	query := `UPDATE faqs SET question=$1, answer=$2, order_index=$3, is_visible=$4
		WHERE id=$5`
	result, err := r.db.Exec(ctx, query, f.Question, f.Answer, f.OrderIndex, f.IsVisible, f.ID)
	if err != nil {
		return err
	}
	if result.RowsAffected() == 0 {
		return fmt.Errorf("faq not found")
	}
	return nil
}

func (r *FAQRepository) Delete(ctx context.Context, id uuid.UUID) error {
	result, err := r.db.Exec(ctx, `DELETE FROM faqs WHERE id = $1`, id)
	if err != nil {
		return err
	}
	if result.RowsAffected() == 0 {
		return fmt.Errorf("faq not found")
	}
	return nil
}

func (r *FAQRepository) Reorder(ctx context.Context, ids []uuid.UUID) error {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)
	for i, id := range ids {
		if _, err := tx.Exec(ctx, `UPDATE faqs SET order_index = $1 WHERE id = $2`, i, id); err != nil {
			return err
		}
	}
	return tx.Commit(ctx)
}
