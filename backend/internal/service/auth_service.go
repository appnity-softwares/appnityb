package service

import (
	"context"
	"errors"
	"fmt"

	"github.com/appnity/backend/internal/config"
	"github.com/appnity/backend/internal/models"
	"github.com/appnity/backend/internal/repository"
	"github.com/appnity/backend/pkg/jwt"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	userRepo   *repository.UserRepository
	jwtService *jwt.JWTService
	config     config.Config
}

func NewAuthService(userRepo *repository.UserRepository, jwtService *jwt.JWTService, cfg config.Config) *AuthService {
	return &AuthService{
		userRepo:   userRepo,
		jwtService: jwtService,
		config:     cfg,
	}
}

func (s *AuthService) Login(ctx context.Context, email, password string) (string, string, error) {
	user, err := s.userRepo.GetByEmail(ctx, email)
	if err != nil {
		return "", "", errors.New("invalid credentials")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		return "", "", errors.New("invalid credentials")
	}

	if !user.IsActive {
		return "", "", errors.New("account is disabled")
	}

	if err := s.userRepo.UpdateLastLogin(ctx, user.ID); err != nil {
		return "", "", fmt.Errorf("failed to update last login: %w", err)
	}

	accessToken, refreshToken, err := s.jwtService.GenerateToken(user.ID, user.Email, user.Role)
	if err != nil {
		return "", "", fmt.Errorf("failed to generate tokens: %w", err)
	}

	return accessToken, refreshToken, nil
}

func (s *AuthService) RefreshToken(ctx context.Context, refreshToken string) (string, string, error) {
	claims, err := s.jwtService.ValidateToken(refreshToken)
	if err != nil {
		return "", "", fmt.Errorf("invalid refresh token: %w", err)
	}

	user, err := s.userRepo.GetByID(ctx, claims.UserID)
	if err != nil {
		return "", "", fmt.Errorf("user not found: %w", err)
	}

	if !user.IsActive {
		return "", "", errors.New("account is disabled")
	}

	newAccess, newRefresh, err := s.jwtService.GenerateToken(user.ID, user.Email, user.Role)
	if err != nil {
		return "", "", fmt.Errorf("failed to generate tokens: %w", err)
	}

	return newAccess, newRefresh, nil
}

func (s *AuthService) GetProfile(ctx context.Context, userID uuid.UUID) (*models.User, error) {
	return s.userRepo.GetByID(ctx, userID)
}

func (s *AuthService) UpdateProfile(ctx context.Context, userID uuid.UUID, updates map[string]interface{}) (*models.User, error) {
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("user not found: %w", err)
	}

	if fullName, ok := updates["full_name"].(string); ok {
		user.FullName = fullName
	}
	if email, ok := updates["email"].(string); ok {
		user.Email = email
	}

	if err := s.userRepo.Update(ctx, user); err != nil {
		return nil, fmt.Errorf("failed to update user: %w", err)
	}

	return user, nil
}

func (s *AuthService) ChangePassword(ctx context.Context, userID uuid.UUID, oldPass, newPass string) error {
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return fmt.Errorf("user not found: %w", err)
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(oldPass)); err != nil {
		return errors.New("incorrect current password")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPass), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("failed to hash password: %w", err)
	}

	user.PasswordHash = string(hashedPassword)
	return s.userRepo.Update(ctx, user)
}
