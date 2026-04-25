package service

import (
	"context"
	"encoding/json"

	"github.com/appnity/backend/internal/models"
	"github.com/appnity/backend/internal/repository"
	"github.com/google/uuid"
)

type ThemeService struct {
	themeRepo *repository.ThemeRepository
}

func NewThemeService(themeRepo *repository.ThemeRepository) *ThemeService {
	return &ThemeService{themeRepo: themeRepo}
}

func (s *ThemeService) GetActiveTheme(ctx context.Context) (*models.Theme, error) {
	return s.themeRepo.GetActive(ctx)
}

func (s *ThemeService) GetAllThemes(ctx context.Context) ([]models.Theme, error) {
	return s.themeRepo.GetAll(ctx)
}

func (s *ThemeService) CreateTheme(ctx context.Context, theme *models.Theme) error {
	return s.themeRepo.Create(ctx, theme)
}

func (s *ThemeService) UpdateTheme(ctx context.Context, theme *models.Theme) error {
	return s.themeRepo.Update(ctx, theme)
}

func (s *ThemeService) SetActiveTheme(ctx context.Context, id uuid.UUID) error {
	return s.themeRepo.SetActive(ctx, id)
}

func (s *ThemeService) ResetToDefaults(ctx context.Context) (*models.Theme, error) {
	// Try to get existing active theme to update it
	existing, err := s.themeRepo.GetActive(ctx)
	
	defaultTheme := &models.Theme{
		Name:         "Default Sunset",
		IsActive:     true,
		Colors:       s.defaultColors(),
		Fonts:        s.defaultFonts(),
		Typography:   s.defaultTypography(),
		BorderRadius: s.defaultBorderRadius(),
		Spacing:      s.defaultSpacing(),
		Shadows:      s.defaultShadows(),
		Gradients:    s.defaultGradients(),
		Animations:   s.defaultAnimations(),
	}

	if err == nil && existing != nil {
		defaultTheme.ID = existing.ID
		if err := s.themeRepo.Update(ctx, defaultTheme); err != nil {
			return nil, err
		}
	} else {
		if err := s.themeRepo.Create(ctx, defaultTheme); err != nil {
			return nil, err
		}
		if err := s.themeRepo.SetActive(ctx, defaultTheme.ID); err != nil {
			return nil, err
		}
	}

	return defaultTheme, nil
}

func (s *ThemeService) defaultColors() json.RawMessage {
	colors := map[string]string{
		"primary":        "#f97316",
		"secondary":      "#111111",
		"accent":         "#f97316",
		"background":     "#e6e6e6",
		"text":           "#111111",
		"surface":        "#ffffff",
		"border":         "#dbdbdb",
		"gradient_start": "#f97316",
		"gradient_end":   "#ff4d00",
		"error":          "#ef4444",
		"success":        "#22c55e",
		"warning":        "#f59e0b",
		"muted":          "#6b7280",
	}
	data, _ := json.Marshal(colors)
	return data
}

func (s *ThemeService) defaultFonts() json.RawMessage {
	fonts := map[string]string{
		"heading": "Plus Jakarta Sans",
		"body":    "Plus Jakarta Sans",
		"cursive": "Pacifico",
	}
	data, _ := json.Marshal(fonts)
	return data
}

func (s *ThemeService) defaultTypography() json.RawMessage {
	typography := map[string]interface{}{
		"h1":    map[string]interface{}{"fontSize": 64, "fontWeight": 800, "lineHeight": 1.1, "letterSpacing": -1.5, "textTransform": "none"},
		"h2":    map[string]interface{}{"fontSize": 48, "fontWeight": 700, "lineHeight": 1.2, "letterSpacing": -0.5, "textTransform": "none"},
		"h3":    map[string]interface{}{"fontSize": 38, "fontWeight": 700, "lineHeight": 1.3, "letterSpacing": 0, "textTransform": "none"},
		"h4":    map[string]interface{}{"fontSize": 30, "fontWeight": 600, "lineHeight": 1.3, "letterSpacing": 0, "textTransform": "none"},
		"h5":    map[string]interface{}{"fontSize": 24, "fontWeight": 600, "lineHeight": 1.4, "letterSpacing": 0, "textTransform": "none"},
		"h6":    map[string]interface{}{"fontSize": 20, "fontWeight": 600, "lineHeight": 1.4, "letterSpacing": 0.5, "textTransform": "uppercase"},
		"p":     map[string]interface{}{"fontSize": 16, "fontWeight": 400, "lineHeight": 1.6, "letterSpacing": 0, "textTransform": "none"},
		"small": map[string]interface{}{"fontSize": 14, "fontWeight": 400, "lineHeight": 1.4, "letterSpacing": 0, "textTransform": "none"},
	}
	data, _ := json.Marshal(typography)
	return data
}

func (s *ThemeService) defaultBorderRadius() json.RawMessage {
	radius := map[string]int{
		"sm":   4,
		"md":   8,
		"lg":   12,
		"xl":   16,
		"full": 9999,
	}
	data, _ := json.Marshal(radius)
	return data
}

func (s *ThemeService) defaultSpacing() json.RawMessage {
	spacing := map[string]int{
		"xs":  4,
		"sm":  8,
		"md":  16,
		"lg":  24,
		"xl":  32,
		"2xl": 48,
	}
	data, _ := json.Marshal(spacing)
	return data
}

func (s *ThemeService) defaultShadows() json.RawMessage {
	shadows := map[string]interface{}{
		"small":  map[string]interface{}{"blur": 4, "spread": 0, "opacity": 0.05},
		"medium": map[string]interface{}{"blur": 12, "spread": 0, "opacity": 0.1},
		"large":  map[string]interface{}{"blur": 24, "spread": 0, "opacity": 0.15},
	}
	data, _ := json.Marshal(shadows)
	return data
}

func (s *ThemeService) defaultGradients() json.RawMessage {
	gradients := map[string]interface{}{
		"primary":   map[string]interface{}{"enabled": true, "start": "#f97316", "end": "#ff4d00", "angle": 135},
		"secondary": map[string]interface{}{"enabled": false, "start": "#6366f1", "end": "#8b5cf6", "angle": 135},
	}
	data, _ := json.Marshal(gradients)
	return data
}

func (s *ThemeService) defaultAnimations() json.RawMessage {
	animations := map[string]interface{}{
		"duration": 300,
		"easing":   "ease-out",
	}
	data, _ := json.Marshal(animations)
	return data
}
