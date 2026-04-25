package models

import (
	"time"

	"github.com/google/uuid"
)

type Navigation struct {
	ID         uuid.UUID  `json:"id"`
	Label      string     `json:"label"`
	URL        string     `json:"url"`
	ParentID   *uuid.UUID `json:"parent_id,omitempty"`
	OrderIndex int        `json:"order_index"`
	IsVisible  bool       `json:"is_visible"`
	IsExternal bool       `json:"is_external"`
	Icon       string     `json:"icon"`
	CreatedAt  time.Time  `json:"created_at"`
}

func (Navigation) TableName() string {
	return "navigations"
}
