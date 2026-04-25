package models

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

type Service struct {
	ID          uuid.UUID       `json:"id"`
	Title       string          `json:"title"`
	Description string          `json:"description"`
	Icon        string          `json:"icon"`
	Image       string          `json:"image"`
	Features    json.RawMessage `json:"features"`
	OrderIndex  int             `json:"order_index"`
	IsVisible   bool            `json:"is_visible"`
	CreatedAt   time.Time       `json:"created_at"`
}

func (Service) TableName() string {
	return "services"
}
