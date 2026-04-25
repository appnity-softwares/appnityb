package models

import (
	"time"

	"github.com/google/uuid"
)

type Media struct {
	ID           uuid.UUID `json:"id"`
	Filename     string    `json:"filename"`
	OriginalName string    `json:"original_name"`
	MimeType     string    `json:"mime_type"`
	Size         int64     `json:"size"`
	URL          string    `json:"url"`
	AltText      string    `json:"alt_text"`
	Folder       string    `json:"folder"`
	UploadedBy   uuid.UUID `json:"uploaded_by"`
	CreatedAt    time.Time `json:"created_at"`
}

func (Media) TableName() string {
	return "media"
}
