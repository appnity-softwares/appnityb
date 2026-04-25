package models

import (
	"time"

	"github.com/google/uuid"
)

type Award struct {
	ID          uuid.UUID `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Image       string    `json:"image"`
	Year        int       `json:"year"`
	OrderIndex  int       `json:"order_index"`
	IsVisible   bool      `json:"is_visible"`
	CreatedAt   time.Time `json:"created_at"`
}

func (Award) TableName() string {
	return "awards"
}
