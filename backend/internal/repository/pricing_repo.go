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

type PricingRepository struct {
	db *pgxpool.Pool
}

func NewPricingRepo(db *pgxpool.Pool) *PricingRepository {
	return &PricingRepository{db: db}
}

func (r *PricingRepository) GetAll(ctx context.Context, visibleOnly bool) ([]models.Pricing, error) {
	var query string
	if visibleOnly {
		query = `SELECT id, plan_name, COALESCE(price_inr, 0), COALESCE(price_usd, 0), COALESCE(tagline, ''), COALESCE(description, ''), COALESCE(features, '[]'), COALESCE(is_highlighted, false), COALESCE(order_index, 0), COALESCE(is_visible, true), COALESCE(icon, 'star'), COALESCE(gradient, 'from-white to-gray-50'), COALESCE(cta_button, 'Get Started'), created_at, updated_at
			FROM pricing WHERE is_visible = true ORDER BY order_index`
	} else {
		query = `SELECT id, plan_name, COALESCE(price_inr, 0), COALESCE(price_usd, 0), COALESCE(tagline, ''), COALESCE(description, ''), COALESCE(features, '[]'), COALESCE(is_highlighted, false), COALESCE(order_index, 0), COALESCE(is_visible, true), COALESCE(icon, 'star'), COALESCE(gradient, 'from-white to-gray-50'), COALESCE(cta_button, 'Get Started'), created_at, updated_at
			FROM pricing ORDER BY order_index`
	}
	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var pricings []models.Pricing
	for rows.Next() {
		var p models.Pricing
		if err := rows.Scan(&p.ID, &p.PlanName, &p.PriceINR, &p.PriceUSD, &p.Tagline, &p.Description, &p.Features, &p.IsHighlighted, &p.OrderIndex, &p.IsVisible, &p.Icon, &p.Gradient, &p.CTAButton, &p.CreatedAt, &p.UpdatedAt); err != nil {
			return nil, err
		}
		pricings = append(pricings, p)
	}
	return pricings, nil
}

func (r *PricingRepository) GetByID(ctx context.Context, id uuid.UUID) (*models.Pricing, error) {
	query := `SELECT id, plan_name, COALESCE(price_inr, 0), COALESCE(price_usd, 0), COALESCE(tagline, ''), COALESCE(description, ''), COALESCE(features, '[]'), COALESCE(is_highlighted, false), COALESCE(order_index, 0), COALESCE(is_visible, true), COALESCE(icon, 'star'), COALESCE(gradient, 'from-white to-gray-50'), COALESCE(cta_button, 'Get Started'), created_at, updated_at
		FROM pricing WHERE id = $1`
	p := &models.Pricing{}
	err := r.db.QueryRow(ctx, query, id).Scan(&p.ID, &p.PlanName, &p.PriceINR, &p.PriceUSD, &p.Tagline, &p.Description, &p.Features, &p.IsHighlighted, &p.OrderIndex, &p.IsVisible, &p.Icon, &p.Gradient, &p.CTAButton, &p.CreatedAt, &p.UpdatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("pricing not found")
		}
		return nil, err
	}
	return p, nil
}

func (r *PricingRepository) Create(ctx context.Context, p *models.Pricing) error {
	p.ID = uuid.New()
	now := time.Now()
	p.CreatedAt = now
	p.UpdatedAt = now
	if p.Icon == "" {
		p.Icon = "star"
	}
	if p.Gradient == "" {
		p.Gradient = "from-white to-gray-50"
	}
	if p.CTAButton == "" {
		p.CTAButton = "Get Started"
	}
	query := `INSERT INTO pricing (id, plan_name, price_inr, price_usd, tagline, description, features, is_highlighted, order_index, is_visible, icon, gradient, cta_button, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15)`
	_, err := r.db.Exec(ctx, query, p.ID, p.PlanName, p.PriceINR, p.PriceUSD, p.Tagline, p.Description, p.Features, p.IsHighlighted, p.OrderIndex, p.IsVisible, p.Icon, p.Gradient, p.CTAButton, p.CreatedAt, p.UpdatedAt)
	return err
}

func (r *PricingRepository) Update(ctx context.Context, p *models.Pricing) error {
	p.UpdatedAt = time.Now()
	if p.Icon == "" {
		p.Icon = "star"
	}
	if p.Gradient == "" {
		p.Gradient = "from-white to-gray-50"
	}
	if p.CTAButton == "" {
		p.CTAButton = "Get Started"
	}
	query := `UPDATE pricing SET plan_name=$1, price_inr=$2, price_usd=$3, tagline=$4, description=$5, features=$6, is_highlighted=$7, order_index=$8, is_visible=$9, icon=$10, gradient=$11, cta_button=$12, updated_at=$13
		WHERE id=$14`
	result, err := r.db.Exec(ctx, query, p.PlanName, p.PriceINR, p.PriceUSD, p.Tagline, p.Description, p.Features, p.IsHighlighted, p.OrderIndex, p.IsVisible, p.Icon, p.Gradient, p.CTAButton, p.UpdatedAt, p.ID)
	if err != nil {
		return err
	}
	if result.RowsAffected() == 0 {
		return fmt.Errorf("pricing not found")
	}
	return nil
}

func (r *PricingRepository) Delete(ctx context.Context, id uuid.UUID) error {
	result, err := r.db.Exec(ctx, `DELETE FROM pricing WHERE id = $1`, id)
	if err != nil {
		return err
	}
	if result.RowsAffected() == 0 {
		return fmt.Errorf("pricing not found")
	}
	return nil
}

func (r *PricingRepository) Reorder(ctx context.Context, ids []uuid.UUID) error {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)
	for i, id := range ids {
		if _, err := tx.Exec(ctx, `UPDATE pricing SET order_index = $1 WHERE id = $2`, i, id); err != nil {
			return err
		}
	}
	return tx.Commit(ctx)
}
