package models

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

type Theme struct {
	ID               uuid.UUID       `json:"id"`
	Name             string          `json:"name"`
	IsActive         bool            `json:"is_active"`
	Colors           json.RawMessage `json:"colors"`
	Fonts            json.RawMessage `json:"fonts"`
	Typography       json.RawMessage `json:"typography"`
	BorderRadius     json.RawMessage `json:"borderRadius"`
	Spacing          json.RawMessage `json:"spacing"`
	Shadows          json.RawMessage `json:"shadows"`
	Breakpoints      json.RawMessage `json:"breakpoints"`
	BackgroundImages json.RawMessage `json:"backgroundImages"`
	Gradients        json.RawMessage `json:"gradients"`
	Animations       json.RawMessage `json:"animations"`
	CreatedAt        time.Time       `json:"created_at"`
	UpdatedAt        time.Time       `json:"updated_at"`
}

func (Theme) TableName() string {
	return "themes"
}
