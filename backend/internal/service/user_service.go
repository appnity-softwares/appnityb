package service

import (
	"context"
	"errors"
	"fmt"

	"github.com/appnity/backend/internal/models"
	"github.com/appnity/backend/internal/repository"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	userRepo *repository.UserRepository
}

func NewUserService(userRepo *repository.UserRepository) *UserService {
	return &UserService{
		userRepo: userRepo,
	}
}

func (s *UserService) GetAllUsers(ctx context.Context) ([]models.User, error) {
	return s.userRepo.GetAll(ctx)
}

func (s *UserService) GetUserByID(ctx context.Context, id uuid.UUID) (*models.User, error) {
	return s.userRepo.GetByID(ctx, id)
}

func (s *UserService) CreateUser(ctx context.Context, email, password, fullName, role string) (*models.User, error) {
	existing, _ := s.userRepo.GetByEmail(ctx, email)
	if existing != nil {
		return nil, errors.New("user with this email already exists")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	user := &models.User{
		Email:        email,
		PasswordHash: string(hashedPassword),
		FullName:     fullName,
		Role:         role,
		IsActive:     true,
	}

	if err := s.userRepo.Create(ctx, user); err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	return user, nil
}

func (s *UserService) UpdateUser(ctx context.Context, id uuid.UUID, updates map[string]interface{}) (*models.User, error) {
	user, err := s.userRepo.GetByID(ctx, id)
	if err != nil {
		return nil, errors.New("user not found")
	}

	if email, ok := updates["email"].(string); ok {
		user.Email = email
	}
	if fullName, ok := updates["full_name"].(string); ok {
		user.FullName = fullName
	}
	if role, ok := updates["role"].(string); ok {
		user.Role = role
	}
	if isActive, ok := updates["is_active"].(bool); ok {
		user.IsActive = isActive
	}
	if password, ok := updates["password"].(string); ok && password != "" {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
		if err != nil {
			return nil, fmt.Errorf("failed to hash password: %w", err)
		}
		user.PasswordHash = string(hashedPassword)
	}

	if err := s.userRepo.Update(ctx, user); err != nil {
		return nil, fmt.Errorf("failed to update user: %w", err)
	}

	return user, nil
}

func (s *UserService) UpdateUserRole(ctx context.Context, id uuid.UUID, role string) (*models.User, error) {
	user, err := s.userRepo.GetByID(ctx, id)
	if err != nil {
		return nil, errors.New("user not found")
	}

	user.Role = role

	if err := s.userRepo.Update(ctx, user); err != nil {
		return nil, fmt.Errorf("failed to update user role: %w", err)
	}

	return user, nil
}

func (s *UserService) DeleteUser(ctx context.Context, id uuid.UUID) error {
	user, err := s.userRepo.GetByID(ctx, id)
	if err != nil {
		return errors.New("user not found")
	}

	if user.Role == "super_admin" {
		return errors.New("cannot delete super admin user")
	}

	return s.userRepo.Delete(ctx, id)
}
