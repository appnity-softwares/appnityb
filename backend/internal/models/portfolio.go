package models

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

type Portfolio struct {
	ID            uuid.UUID       `json:"id"`
	Title         string          `json:"title"`
	Slug          string          `json:"slug"`
	Description   string          `json:"description"`
	Tagline       string          `json:"tagline"` // Map description to tagline if needed
	CoverImage    string          `json:"cover_image"`
	GalleryImages json.RawMessage `json:"gallery_images"`
	Category      string          `json:"category"`
	Industry      string          `json:"industry"`
	Year          int             `json:"year"`
	Stack         string          `json:"stack"`
	Problem       string          `json:"problem"`
	Approach      json.RawMessage `json:"approach"`
	ClientName    string          `json:"client_name"`
	ProjectURL    string          `json:"project_url"`
	Metrics       json.RawMessage `json:"metrics"`
	IsFeatured    bool            `json:"is_featured"`
	IsVisible     bool            `json:"is_visible"`
	OrderIndex    int             `json:"order_index"`
	CreatedAt     time.Time       `json:"created_at"`
	UpdatedAt     time.Time       `json:"updated_at"`
}

func (Portfolio) TableName() string {
	return "portfolios"
}
