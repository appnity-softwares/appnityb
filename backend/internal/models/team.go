package models

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

type Team struct {
	ID          uuid.UUID       `json:"id"`
	FullName    string          `json:"full_name"`
	Role        string          `json:"role"`
	Bio         string          `json:"bio"`
	Photo       string          `json:"photo"`
	SocialLinks json.RawMessage `json:"social_links"`
	OrderIndex  int             `json:"order_index"`
	IsVisible   bool            `json:"is_visible"`
	CreatedAt   time.Time       `json:"created_at"`
	UpdatedAt   time.Time       `json:"updated_at"`
}

func (Team) TableName() string {
	return "teams"
}
