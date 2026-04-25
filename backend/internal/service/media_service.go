package service

import (
	"context"
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"

	"github.com/appnity/backend/internal/models"
	"github.com/appnity/backend/internal/repository"
	"github.com/google/uuid"
)

type MediaService struct {
	mediaRepo     *repository.MediaRepository
	uploadDir     string
	maxUploadSize int64
	baseURL       string
}

func NewMediaService(mediaRepo *repository.MediaRepository, uploadDir string, maxUploadSize int64, baseURL string) *MediaService {
	return &MediaService{
		mediaRepo:     mediaRepo,
		uploadDir:     uploadDir,
		maxUploadSize: maxUploadSize,
		baseURL:       baseURL,
	}
}

func (s *MediaService) UploadFile(ctx context.Context, file multipart.File, header *multipart.FileHeader) (*models.Media, error) {
	if header.Size > s.maxUploadSize {
		return nil, fmt.Errorf("file size exceeds maximum allowed size of %d bytes", s.maxUploadSize)
	}

	ext := filepath.Ext(header.Filename)
	uniqueName := fmt.Sprintf("%s%s", uuid.New().String(), ext)

	folderPath := filepath.Join(s.uploadDir)
	if err := os.MkdirAll(folderPath, 0755); err != nil {
		return nil, fmt.Errorf("failed to create upload directory: %w", err)
	}

	filePath := filepath.Join(folderPath, uniqueName)
	dst, err := os.Create(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to create file: %w", err)
	}
	defer dst.Close()

	if _, err := io.Copy(dst, file); err != nil {
		os.Remove(filePath)
		return nil, fmt.Errorf("failed to save file: %w", err)
	}

	media := &models.Media{
		Filename:     uniqueName,
		OriginalName: header.Filename,
		MimeType:     header.Header.Get("Content-Type"),
		Size:         header.Size,
		URL:          fmt.Sprintf("%s/uploads/%s", s.baseURL, uniqueName),
		Folder:       "",
		AltText:      "",
	}

	if err := s.mediaRepo.Create(ctx, media); err != nil {
		os.Remove(filePath)
		return nil, fmt.Errorf("failed to save media record: %w", err)
	}

	return media, nil
}

func (s *MediaService) GetMedia(ctx context.Context, folder string, page, pageSize int) ([]models.Media, int64, error) {
	return s.mediaRepo.GetAll(ctx, folder, page, pageSize)
}

func (s *MediaService) DeleteMedia(ctx context.Context, id uuid.UUID) error {
	media, err := s.mediaRepo.GetByID(ctx, id)
	if err != nil {
		return err
	}

	filePath := filepath.Join(s.uploadDir, media.Filename)
	os.Remove(filePath)

	return s.mediaRepo.Delete(ctx, id)
}

func (s *MediaService) UpdateMedia(ctx context.Context, id uuid.UUID, updates map[string]interface{}) (*models.Media, error) {
	media, err := s.mediaRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if altText, ok := updates["alt_text"].(string); ok {
		media.AltText = altText
	}
	if folder, ok := updates["folder"].(string); ok {
		media.Folder = folder
	}
	if originalName, ok := updates["original_name"].(string); ok {
		media.OriginalName = originalName
	}

	if err := s.mediaRepo.Update(ctx, media); err != nil {
		return nil, err
	}

	return media, nil
}

func (s *MediaService) GetMediaByID(ctx context.Context, id uuid.UUID) (*models.Media, error) {
	return s.mediaRepo.GetByID(ctx, id)
}

func (s *MediaService) Count(ctx context.Context, folder string) (int64, error) {
	return s.mediaRepo.Count(ctx, folder)
}
