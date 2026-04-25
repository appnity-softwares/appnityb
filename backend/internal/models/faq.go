package models

import (
	"time"

	"github.com/google/uuid"
)

type FAQ struct {
	ID         uuid.UUID `json:"id"`
	Question   string    `json:"question"`
	Answer     string    `json:"answer"`
	OrderIndex int       `json:"order_index"`
	IsVisible  bool      `json:"is_visible"`
	CreatedAt  time.Time `json:"created_at"`
}

func (FAQ) TableName() string {
	return "faqs"
}
