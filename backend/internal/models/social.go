package models

import (
	"time"

	"github.com/google/uuid"
)

type Social struct {
	ID         uuid.UUID `json:"id"`
	Platform   string    `json:"platform"`
	URL        string    `json:"url"`
	IsVisible  bool      `json:"is_visible"`
	Icon       string    `json:"icon"`
	OrderIndex int       `json:"order_index"`
	CreatedAt  time.Time `json:"created_at"`
}

func (Social) TableName() string {
	return "social_links"
}
