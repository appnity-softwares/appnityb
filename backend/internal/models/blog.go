package models

import (
	"time"

	"github.com/google/uuid"
)

type Blog struct {
	ID          uuid.UUID  `json:"id"`
	Title       string     `json:"title"`
	Slug        string     `json:"slug"`
	Excerpt     string     `json:"excerpt"`
	Content     string     `json:"content"`
	CoverImage  string     `json:"cover_image"`
	Author      string     `json:"author"`
	Category    string     `json:"category"`
	Tags        []string   `json:"tags"`
	ReadTime    int        `json:"read_time"`
	IsPublished bool       `json:"is_published"`
	PublishedAt *time.Time `json:"published_at,omitempty"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
}

func (Blog) TableName() string {
	return "blogs"
}
