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

type MediaRepository struct {
	db *pgxpool.Pool
}

func NewMediaRepo(db *pgxpool.Pool) *MediaRepository {
	return &MediaRepository{db: db}
}

func (r *MediaRepository) GetAll(ctx context.Context, folderFilter string, page, pageSize int) ([]models.Media, int64, error) {
	var countQuery string
	var dataQuery string
	var countArgs []interface{}
	var dataArgs []interface{}

	if folderFilter != "" {
		countQuery = `SELECT COUNT(*) FROM media WHERE folder = $1`
		dataQuery = `SELECT id, filename, original_name, mime_type, size, url, alt_text, folder, uploaded_by, created_at
			FROM media WHERE folder = $3 ORDER BY created_at DESC LIMIT $1 OFFSET $2`
		countArgs = append(countArgs, folderFilter)
		dataArgs = append(dataArgs, pageSize, (page-1)*pageSize, folderFilter)
	} else {
		countQuery = `SELECT COUNT(*) FROM media`
		dataQuery = `SELECT id, filename, original_name, mime_type, size, url, alt_text, folder, uploaded_by, created_at
			FROM media ORDER BY created_at DESC LIMIT $1 OFFSET $2`
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

	var media []models.Media
	for rows.Next() {
		var m models.Media
		if err := rows.Scan(&m.ID, &m.Filename, &m.OriginalName, &m.MimeType, &m.Size, &m.URL, &m.AltText, &m.Folder, &m.UploadedBy, &m.CreatedAt); err != nil {
			return nil, 0, err
		}
		media = append(media, m)
	}
	return media, total, nil
}

func (r *MediaRepository) GetByID(ctx context.Context, id uuid.UUID) (*models.Media, error) {
	query := `SELECT id, filename, original_name, mime_type, size, url, alt_text, folder, uploaded_by, created_at
		FROM media WHERE id = $1`
	m := &models.Media{}
	err := r.db.QueryRow(ctx, query, id).Scan(&m.ID, &m.Filename, &m.OriginalName, &m.MimeType, &m.Size, &m.URL, &m.AltText, &m.Folder, &m.UploadedBy, &m.CreatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("media not found")
		}
		return nil, err
	}
	return m, nil
}

func (r *MediaRepository) Create(ctx context.Context, m *models.Media) error {
	m.ID = uuid.New()
	m.CreatedAt = time.Now()
	query := `INSERT INTO media (id, filename, original_name, mime_type, size, url, alt_text, folder, uploaded_by, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)`
	_, err := r.db.Exec(ctx, query, m.ID, m.Filename, m.OriginalName, m.MimeType, m.Size, m.URL, m.AltText, m.Folder, m.UploadedBy, m.CreatedAt)
	return err
}

func (r *MediaRepository) Update(ctx context.Context, m *models.Media) error {
	query := `UPDATE media SET filename=$1, original_name=$2, mime_type=$3, size=$4, url=$5, alt_text=$6, folder=$7, uploaded_by=$8
		WHERE id=$9`
	result, err := r.db.Exec(ctx, query, m.Filename, m.OriginalName, m.MimeType, m.Size, m.URL, m.AltText, m.Folder, m.UploadedBy, m.ID)
	if err != nil {
		return err
	}
	if result.RowsAffected() == 0 {
		return fmt.Errorf("media not found")
	}
	return nil
}

func (r *MediaRepository) Delete(ctx context.Context, id uuid.UUID) error {
	result, err := r.db.Exec(ctx, `DELETE FROM media WHERE id = $1`, id)
	if err != nil {
		return err
	}
	if result.RowsAffected() == 0 {
		return fmt.Errorf("media not found")
	}
	return nil
}

func (r *MediaRepository) Count(ctx context.Context, folderFilter string) (int64, error) {
	var query string
	var args []interface{}
	if folderFilter != "" {
		query = `SELECT COUNT(*) FROM media WHERE folder = $1`
		args = append(args, folderFilter)
	} else {
		query = `SELECT COUNT(*) FROM media`
	}
	var count int64
	err := r.db.QueryRow(ctx, query, args...).Scan(&count)
	return count, err
}
