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

type BlogRepository struct {
	db *pgxpool.Pool
}

func NewBlogRepo(db *pgxpool.Pool) *BlogRepository {
	return &BlogRepository{db: db}
}

func (r *BlogRepository) GetAll(ctx context.Context, publishedOnly bool, page, pageSize int) ([]models.Blog, int64, error) {
	var countQuery string
	var dataQuery string

	if publishedOnly {
		countQuery = `SELECT COUNT(*) FROM blogs WHERE is_published = true`
		dataQuery = `SELECT id, title, slug, excerpt, content, cover_image, author, category, tags, read_time, is_published, published_at, created_at, updated_at
			FROM blogs WHERE is_published = true ORDER BY COALESCE(published_at, created_at) DESC LIMIT $1 OFFSET $2`
	} else {
		countQuery = `SELECT COUNT(*) FROM blogs`
		dataQuery = `SELECT id, title, slug, excerpt, content, cover_image, author, category, tags, read_time, is_published, published_at, created_at, updated_at
			FROM blogs ORDER BY COALESCE(published_at, created_at) DESC LIMIT $1 OFFSET $2`
	}

	var total int64
	if err := r.db.QueryRow(ctx, countQuery).Scan(&total); err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	rows, err := r.db.Query(ctx, dataQuery, pageSize, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var blogs []models.Blog
	for rows.Next() {
		var b models.Blog
		if err := rows.Scan(&b.ID, &b.Title, &b.Slug, &b.Excerpt, &b.Content, &b.CoverImage,
			&b.Author, &b.Category, &b.Tags, &b.ReadTime, &b.IsPublished, &b.PublishedAt,
			&b.CreatedAt, &b.UpdatedAt); err != nil {
			return nil, 0, err
		}
		blogs = append(blogs, b)
	}
	return blogs, total, nil
}

func (r *BlogRepository) GetBySlug(ctx context.Context, slug string) (*models.Blog, error) {
	query := `SELECT id, title, slug, excerpt, content, cover_image, author, category, tags, read_time, is_published, published_at, created_at, updated_at
		FROM blogs WHERE slug = $1`
	b := &models.Blog{}
	err := r.db.QueryRow(ctx, query, slug).Scan(&b.ID, &b.Title, &b.Slug, &b.Excerpt, &b.Content,
		&b.CoverImage, &b.Author, &b.Category, &b.Tags, &b.ReadTime, &b.IsPublished, &b.PublishedAt,
		&b.CreatedAt, &b.UpdatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("blog not found")
		}
		return nil, err
	}
	return b, nil
}

func (r *BlogRepository) GetByID(ctx context.Context, id uuid.UUID) (*models.Blog, error) {
	query := `SELECT id, title, slug, excerpt, content, cover_image, author, category, tags, read_time, is_published, published_at, created_at, updated_at
		FROM blogs WHERE id = $1`
	b := &models.Blog{}
	err := r.db.QueryRow(ctx, query, id).Scan(&b.ID, &b.Title, &b.Slug, &b.Excerpt, &b.Content,
		&b.CoverImage, &b.Author, &b.Category, &b.Tags, &b.ReadTime, &b.IsPublished, &b.PublishedAt,
		&b.CreatedAt, &b.UpdatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("blog not found")
		}
		return nil, err
	}
	return b, nil
}

func (r *BlogRepository) Create(ctx context.Context, b *models.Blog) error {
	b.ID = uuid.New()
	now := time.Now()
	b.CreatedAt = now
	b.UpdatedAt = now
	query := `INSERT INTO blogs (id, title, slug, excerpt, content, cover_image, author, category, tags, read_time, is_published, published_at, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14)`
	_, err := r.db.Exec(ctx, query, b.ID, b.Title, b.Slug, b.Excerpt, b.Content, b.CoverImage,
		b.Author, b.Category, b.Tags, b.ReadTime, b.IsPublished, b.PublishedAt, b.CreatedAt, b.UpdatedAt)
	return err
}

func (r *BlogRepository) Update(ctx context.Context, b *models.Blog) error {
	b.UpdatedAt = time.Now()
	query := `UPDATE blogs SET title=$1, slug=$2, excerpt=$3, content=$4, cover_image=$5, author=$6, category=$7, tags=$8, read_time=$9, is_published=$10, published_at=$11, updated_at=$12
		WHERE id=$13`
	result, err := r.db.Exec(ctx, query, b.Title, b.Slug, b.Excerpt, b.Content, b.CoverImage,
		b.Author, b.Category, b.Tags, b.ReadTime, b.IsPublished, b.PublishedAt, b.UpdatedAt, b.ID)
	if err != nil {
		return err
	}
	if result.RowsAffected() == 0 {
		return fmt.Errorf("blog not found")
	}
	return nil
}

func (r *BlogRepository) Delete(ctx context.Context, id uuid.UUID) error {
	result, err := r.db.Exec(ctx, `DELETE FROM blogs WHERE id = $1`, id)
	if err != nil {
		return err
	}
	if result.RowsAffected() == 0 {
		return fmt.Errorf("blog not found")
	}
	return nil
}

func (r *BlogRepository) GetRelated(ctx context.Context, id uuid.UUID, category string, limit int) ([]models.Blog, error) {
	query := `SELECT id, title, slug, excerpt, content, cover_image, author, category, tags, read_time, is_published, published_at, created_at, updated_at
		FROM blogs WHERE category = $1 AND id != $2 AND is_published = true ORDER BY COALESCE(published_at, created_at) DESC LIMIT $3`
	rows, err := r.db.Query(ctx, query, category, id, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var blogs []models.Blog
	for rows.Next() {
		var b models.Blog
		if err := rows.Scan(&b.ID, &b.Title, &b.Slug, &b.Excerpt, &b.Content, &b.CoverImage,
			&b.Author, &b.Category, &b.Tags, &b.ReadTime, &b.IsPublished, &b.PublishedAt,
			&b.CreatedAt, &b.UpdatedAt); err != nil {
			return nil, err
		}
		blogs = append(blogs, b)
	}
	return blogs, nil
}

func (r *BlogRepository) GetRelatedBySlug(ctx context.Context, slug string, limit int) ([]models.Blog, error) {
	// First get the blog to find its category and ID
	query := `SELECT id, category FROM blogs WHERE slug = $1`
	var id uuid.UUID
	var category string
	err := r.db.QueryRow(ctx, query, slug).Scan(&id, &category)
	if err != nil {
		return nil, err
	}

	return r.GetRelated(ctx, id, category, limit)
}
