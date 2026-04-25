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

type JobOpeningRepository struct {
	db *pgxpool.Pool
}

func NewJobOpeningRepo(db *pgxpool.Pool) *JobOpeningRepository {
	return &JobOpeningRepository{db: db}
}

func (r *JobOpeningRepository) GetAll(ctx context.Context, activeOnly bool) ([]models.JobOpening, error) {
	var query string
	if activeOnly {
		query = `SELECT id, title, department, location, type, description, requirements, is_active, apply_url, created_at, updated_at
			FROM job_openings WHERE is_active = true ORDER BY created_at DESC`
	} else {
		query = `SELECT id, title, department, location, type, description, requirements, is_active, apply_url, created_at, updated_at
			FROM job_openings ORDER BY created_at DESC`
	}
	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var jobs []models.JobOpening
	for rows.Next() {
		var j models.JobOpening
		if err := rows.Scan(&j.ID, &j.Title, &j.Department, &j.Location, &j.Type, &j.Description, &j.Requirements, &j.IsActive, &j.ApplyURL, &j.CreatedAt, &j.UpdatedAt); err != nil {
			return nil, err
		}
		jobs = append(jobs, j)
	}
	return jobs, nil
}

func (r *JobOpeningRepository) GetByID(ctx context.Context, id uuid.UUID) (*models.JobOpening, error) {
	query := `SELECT id, title, department, location, type, description, requirements, is_active, apply_url, created_at, updated_at
		FROM job_openings WHERE id = $1`
	j := &models.JobOpening{}
	err := r.db.QueryRow(ctx, query, id).Scan(&j.ID, &j.Title, &j.Department, &j.Location, &j.Type, &j.Description, &j.Requirements, &j.IsActive, &j.ApplyURL, &j.CreatedAt, &j.UpdatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("job opening not found")
		}
		return nil, err
	}
	return j, nil
}

func (r *JobOpeningRepository) Create(ctx context.Context, j *models.JobOpening) error {
	j.ID = uuid.New()
	now := time.Now()
	j.CreatedAt = now
	j.UpdatedAt = now
	query := `INSERT INTO job_openings (id, title, department, location, type, description, requirements, is_active, apply_url, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)`
	_, err := r.db.Exec(ctx, query, j.ID, j.Title, j.Department, j.Location, j.Type, j.Description, j.Requirements, j.IsActive, j.ApplyURL, j.CreatedAt, j.UpdatedAt)
	return err
}

func (r *JobOpeningRepository) Update(ctx context.Context, j *models.JobOpening) error {
	j.UpdatedAt = time.Now()
	query := `UPDATE job_openings SET title=$1, department=$2, location=$3, type=$4, description=$5, requirements=$6, is_active=$7, apply_url=$8, updated_at=$9
		WHERE id=$10`
	result, err := r.db.Exec(ctx, query, j.Title, j.Department, j.Location, j.Type, j.Description, j.Requirements, j.IsActive, j.ApplyURL, j.UpdatedAt, j.ID)
	if err != nil {
		return err
	}
	if result.RowsAffected() == 0 {
		return fmt.Errorf("job opening not found")
	}
	return nil
}

func (r *JobOpeningRepository) Delete(ctx context.Context, id uuid.UUID) error {
	result, err := r.db.Exec(ctx, `DELETE FROM job_openings WHERE id = $1`, id)
	if err != nil {
		return err
	}
	if result.RowsAffected() == 0 {
		return fmt.Errorf("job opening not found")
	}
	return nil
}
