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

type UserRepository struct {
	db *pgxpool.Pool
}

func NewUserRepo(db *pgxpool.Pool) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) Create(ctx context.Context, u *models.User) error {
	u.ID = uuid.New()
	now := time.Now()
	u.CreatedAt = now
	u.UpdatedAt = now
	query := `INSERT INTO users (id, email, password_hash, full_name, role, is_active, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`
	_, err := r.db.Exec(ctx, query, u.ID, u.Email, u.PasswordHash, u.FullName, u.Role, u.IsActive, u.CreatedAt, u.UpdatedAt)
	return err
}

func (r *UserRepository) GetByEmail(ctx context.Context, email string) (*models.User, error) {
	query := `SELECT id, email, password_hash, 
		COALESCE(full_name, ''), 
		COALESCE(role, 'admin'), 
		COALESCE(is_active, true), 
		COALESCE(created_at, NOW()), 
		COALESCE(updated_at, NOW()), 
		last_login
		FROM users WHERE email = $1`
	u := &models.User{}
	err := r.db.QueryRow(ctx, query, email).Scan(
		&u.ID, &u.Email, &u.PasswordHash, &u.FullName, &u.Role, &u.IsActive,
		&u.CreatedAt, &u.UpdatedAt, &u.LastLogin,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("user not found")
		}
		return nil, fmt.Errorf("error scanning user: %w", err)
	}
	return u, nil
}

func (r *UserRepository) GetByID(ctx context.Context, id uuid.UUID) (*models.User, error) {
	query := `SELECT id, email, password_hash, 
		COALESCE(full_name, ''), 
		COALESCE(role, 'admin'), 
		COALESCE(is_active, true), 
		COALESCE(created_at, NOW()), 
		COALESCE(updated_at, NOW()), 
		last_login
		FROM users WHERE id = $1`
	u := &models.User{}
	err := r.db.QueryRow(ctx, query, id).Scan(
		&u.ID, &u.Email, &u.PasswordHash, &u.FullName, &u.Role, &u.IsActive,
		&u.CreatedAt, &u.UpdatedAt, &u.LastLogin,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("user not found")
		}
		return nil, err
	}
	return u, nil
}

func (r *UserRepository) Update(ctx context.Context, u *models.User) error {
	u.UpdatedAt = time.Now()
	query := `UPDATE users SET email=$1, full_name=$2, role=$3, is_active=$4, updated_at=$5
		WHERE id=$6`
	result, err := r.db.Exec(ctx, query, u.Email, u.FullName, u.Role, u.IsActive, u.UpdatedAt, u.ID)
	if err != nil {
		return err
	}
	if result.RowsAffected() == 0 {
		return fmt.Errorf("user not found")
	}
	return nil
}

func (r *UserRepository) UpdateLastLogin(ctx context.Context, id uuid.UUID) error {
	now := time.Now()
	query := `UPDATE users SET last_login=$1 WHERE id=$2`
	_, err := r.db.Exec(ctx, query, now, id)
	return err
}

func (r *UserRepository) UpdatePassword(ctx context.Context, id uuid.UUID, passwordHash string) error {
	now := time.Now()
	query := `UPDATE users SET password_hash=$1, updated_at=$2 WHERE id=$3`
	result, err := r.db.Exec(ctx, query, passwordHash, now, id)
	if err != nil {
		return err
	}
	if result.RowsAffected() == 0 {
		return fmt.Errorf("user not found")
	}
	return nil
}

func (r *UserRepository) Delete(ctx context.Context, id uuid.UUID) error {
	result, err := r.db.Exec(ctx, `DELETE FROM users WHERE id=$1`, id)
	if err != nil {
		return err
	}
	if result.RowsAffected() == 0 {
		return fmt.Errorf("user not found")
	}
	return nil
}

func (r *UserRepository) GetAll(ctx context.Context) ([]models.User, error) {
	query := `SELECT id, email, password_hash, 
		COALESCE(full_name, ''), 
		COALESCE(role, 'admin'), 
		COALESCE(is_active, true), 
		COALESCE(created_at, NOW()), 
		COALESCE(updated_at, NOW()), 
		last_login
		FROM users ORDER BY created_at DESC`
	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []models.User
	for rows.Next() {
		u := models.User{}
		err := rows.Scan(
			&u.ID, &u.Email, &u.PasswordHash, &u.FullName, &u.Role, &u.IsActive,
			&u.CreatedAt, &u.UpdatedAt, &u.LastLogin,
		)
		if err != nil {
			return nil, err
		}
		users = append(users, u)
	}
	return users, nil
}
