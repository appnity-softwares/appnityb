package models

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

type Pricing struct {
	ID            uuid.UUID       `json:"id"`
	PlanName      string          `json:"plan_name"`
	PriceINR      float64         `json:"price_inr"`
	PriceUSD      float64         `json:"price_usd"`
	Tagline       string          `json:"tagline"`
	Description   string          `json:"description"`
	Features      json.RawMessage `json:"features"`
	IsHighlighted bool            `json:"is_highlighted"`
	OrderIndex    int             `json:"order_index"`
	IsVisible     bool            `json:"is_visible"`
	Icon          string          `json:"icon"`
	Gradient      string          `json:"gradient"`
	CTAButton     string          `json:"cta_button"`
	CreatedAt     time.Time       `json:"created_at"`
	UpdatedAt     time.Time       `json:"updated_at"`
}

func (Pricing) TableName() string {
	return "pricing"
}
