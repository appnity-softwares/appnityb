package repository

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/appnity/backend/internal/models"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type ContactRepository struct {
	db *pgxpool.Pool
}

func NewContactRepo(db *pgxpool.Pool) *ContactRepository {
	return &ContactRepository{db: db}
}

func (r *ContactRepository) GetAll(ctx context.Context, statusFilter string, page, pageSize int) ([]models.Contact, int64, error) {
	var countQuery string
	var dataQuery string
	var countArgs []interface{}
	var dataArgs []interface{}

	if statusFilter != "" {
		countQuery = `SELECT COUNT(*) FROM contacts WHERE status = $1`
		dataQuery = `SELECT id, name, email, service, message, status, admin_notes, created_at, updated_at
			FROM contacts WHERE status = $3 ORDER BY created_at DESC LIMIT $1 OFFSET $2`
		countArgs = append(countArgs, statusFilter)
		dataArgs = append(dataArgs, pageSize, (page-1)*pageSize, statusFilter)
	} else {
		countQuery = `SELECT COUNT(*) FROM contacts`
		dataQuery = `SELECT id, name, email, service, message, status, admin_notes, created_at, updated_at
			FROM contacts ORDER BY created_at DESC LIMIT $1 OFFSET $2`
		dataArgs = append(dataArgs, pageSize, (page-1)*pageSize)
	}

	var total int64
	if err := r.db.QueryRow(ctx, countQuery, countArgs...).Scan(&total); err != nil {
		return nil, 0, err
	}

	rows, err := r.db.Query(ctx, dataQuery, dataArgs...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var contacts []models.Contact
	for rows.Next() {
		var c models.Contact
		if err := rows.Scan(&c.ID, &c.Name, &c.Email, &c.Service, &c.Message, &c.Status, &c.AdminNotes, &c.CreatedAt, &c.UpdatedAt); err != nil {
			return nil, 0, err
		}
		contacts = append(contacts, c)
	}
	return contacts, total, nil
}

func (r *ContactRepository) GetByID(ctx context.Context, id uuid.UUID) (*models.Contact, error) {
	query := `SELECT id, name, email, service, message, status, admin_notes, created_at, updated_at
		FROM contacts WHERE id = $1`
	c := &models.Contact{}
	err := r.db.QueryRow(ctx, query, id).Scan(&c.ID, &c.Name, &c.Email, &c.Service, &c.Message, &c.Status, &c.AdminNotes, &c.CreatedAt, &c.UpdatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("contact not found")
		}
		return nil, err
	}
	return c, nil
}

func (r *ContactRepository) Create(ctx context.Context, c *models.Contact) error {
	c.ID = uuid.New()
	now := time.Now()
	c.CreatedAt = now
	c.UpdatedAt = now
	if c.Status == "" {
		c.Status = "new"
	}
	query := `INSERT INTO contacts (id, name, email, service, message, status, admin_notes, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)`
	_, err := r.db.Exec(ctx, query, c.ID, c.Name, c.Email, c.Service, c.Message, c.Status, c.AdminNotes, c.CreatedAt, c.UpdatedAt)
	return err
}

func (r *ContactRepository) Update(ctx context.Context, c *models.Contact) error {
	c.UpdatedAt = time.Now()
	query := `UPDATE contacts SET name=$1, email=$2, service=$3, message=$4, status=$5, admin_notes=$6, updated_at=$7
		WHERE id=$8`
	result, err := r.db.Exec(ctx, query, c.Name, c.Email, c.Service, c.Message, c.Status, c.AdminNotes, c.UpdatedAt, c.ID)
	if err != nil {
		return err
	}
	if result.RowsAffected() == 0 {
		return fmt.Errorf("contact not found")
	}
	return nil
}

func (r *ContactRepository) Delete(ctx context.Context, id uuid.UUID) error {
	result, err := r.db.Exec(ctx, `DELETE FROM contacts WHERE id = $1`, id)
	if err != nil {
		return err
	}
	if result.RowsAffected() == 0 {
		return fmt.Errorf("contact not found")
	}
	return nil
}

func (r *ContactRepository) CountByStatus(ctx context.Context) (map[string]int64, error) {
	query := `SELECT status, COUNT(*) FROM contacts GROUP BY status`
	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	counts := make(map[string]int64)
	for rows.Next() {
		var status string
		var count int64
		if err := rows.Scan(&status, &count); err != nil {
			return nil, err
		}
		counts[status] = count
	}
	return counts, nil
}
