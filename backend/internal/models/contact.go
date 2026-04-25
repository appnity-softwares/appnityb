package models

import (
	"time"

	"github.com/google/uuid"
)

type Contact struct {
	ID         uuid.UUID `json:"id"`
	Name       string    `json:"name"`
	Email      string    `json:"email"`
	Phone      string    `json:"phone"`
	Service    string    `json:"service"`
	Subject    string    `json:"subject"`
	Message    string    `json:"message"`
	Status     string    `json:"status"`
	AdminNotes *string   `json:"admin_notes,omitempty"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

func (Contact) TableName() string {
	return "contacts"
}
