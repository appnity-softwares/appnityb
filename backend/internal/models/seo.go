package models

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

type SEO struct {
	ID                uuid.UUID       `json:"id"`
	SiteTitle         string          `json:"site_title"`
	SiteDescription   string          `json:"site_description"`
	Keywords          []string        `json:"keywords"`
	OgImage           string          `json:"og_image"`
	TwitterHandle     string          `json:"twitter_handle"`
	CanonicalURL      string          `json:"canonical_url"`
	RobotsTxt         string          `json:"robots_txt"`
	GoogleAnalyticsID string          `json:"google_analytics_id"`
	GoogleTagManager  string          `json:"google_tag_manager"`
	JSONLDSchema      json.RawMessage `json:"json_ld_schema"`
	Favicon           string          `json:"favicon"`
	CreatedAt         time.Time       `json:"created_at"`
	UpdatedAt         time.Time       `json:"updated_at"`
}

func (SEO) TableName() string {
	return "seo_settings"
}
