package models

import (
	"time"

	"github.com/google/uuid"
)

type Value struct {
	ID          uuid.UUID `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Icon        string    `json:"icon"`
	OrderIndex  int       `json:"order_index"`
	IsVisible   bool      `json:"is_visible"`
	CreatedAt   time.Time `json:"created_at"`
}

func (Value) TableName() string {
	return "values"
}
