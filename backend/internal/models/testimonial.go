package models

import (
	"time"

	"github.com/google/uuid"
)

type Testimonial struct {
	ID         uuid.UUID `json:"id"`
	Name       string    `json:"name"`
	Role       string    `json:"role"`
	Company    string    `json:"company"`
	Avatar     string    `json:"avatar"`
	Quote      string    `json:"quote"`
	Rating     int       `json:"rating"`
	VideoURL   string    `json:"video_url"`
	OrderIndex int       `json:"order_index"`
	IsVisible  bool      `json:"is_visible"`
	CreatedAt  time.Time `json:"created_at"`
}

func (Testimonial) TableName() string {
	return "testimonials"
}
