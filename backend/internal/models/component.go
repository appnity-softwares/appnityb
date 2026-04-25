package models

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

type Component struct {
	ID            uuid.UUID       `json:"id"`
	PageID        uuid.UUID       `json:"page_id"`
	ComponentType string          `json:"component_type"`
	IsVisible     bool            `json:"is_visible"`
	Content       json.RawMessage `json:"content"`
	OrderIndex    int             `json:"order_index"`
	CreatedAt     time.Time       `json:"created_at"`
	UpdatedAt     time.Time       `json:"updated_at"`
}

func (Component) TableName() string {
	return "components"
}
