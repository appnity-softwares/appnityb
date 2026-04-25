package models

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

type JobOpening struct {
	ID           uuid.UUID       `json:"id"`
	Title        string          `json:"title"`
	Department   string          `json:"department"`
	Location     string          `json:"location"`
	Type         string          `json:"type"`
	Description  string          `json:"description"`
	Requirements json.RawMessage `json:"requirements"`
	IsActive     bool            `json:"is_active"`
	ApplyURL     string          `json:"apply_url"`
	CreatedAt    time.Time       `json:"created_at"`
	UpdatedAt    time.Time       `json:"updated_at"`
}

func (JobOpening) TableName() string {
	return "job_openings"
}
