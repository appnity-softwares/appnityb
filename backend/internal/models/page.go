package models

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

type Page struct {
	ID              uuid.UUID       `json:"id"`
	Slug            string          `json:"slug"`
	Title           string          `json:"title"`
	MetaDescription string          `json:"meta_description"`
	MetaKeywords    string          `json:"meta_keywords"`
	OgImage         string          `json:"og_image"`
	CanonicalURL    string          `json:"canonical_url"`
	BackgroundColor string          `json:"background_color"`
	CustomCSS       string          `json:"custom_css"`
	IsPublished     bool            `json:"is_published"`
	Sections        json.RawMessage `json:"sections"`
	CreatedAt       time.Time       `json:"created_at"`
	UpdatedAt       time.Time       `json:"updated_at"`
}

func (Page) TableName() string {
	return "pages"
}
