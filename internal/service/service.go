package service

import (
	"context"

	"url-shortener/internal/repository"

	"go.uber.org/zap"
)

type Shortener interface {
	CreateShortURL(ctx context.Context, originalURL string) (string, error)
	GetOriginalURL(ctx context.Context, shortCode string) (string, error)
}

type shortenerService struct {
	repo   repository.URLRepository
	cache  repository.Cache
	logger *zap.Logger
}

func NewShortenerService(repo repository.URLRepository, cache repository.Cache, l *zap.Logger) Shortener {
	return &shortenerService{
		repo:   repo,
		cache:  cache,
		logger: l,
	}
}

func (s *shortenerService) CreateShortURL(ctx context.Context, originalURL string) (string, error) {

	if cached, err := s.cache.Get(ctx, originalURL); err == nil {
		return cached, nil
	}

	// Создание в базе
	id, err := s.repo.SaveURL(ctx, originalURL)
	if err != nil {
		return "", err
	}

	// Генерация short code
	shortCode := utils.GenerateShortCode(id)

	// Сохранение в кэш
	if err := s.cache.Set(ctx, shortCode, originalURL); err != nil {
		s.logger.Warn("Cache set failed", zap.Error(err))
	}

	return shortCode, nil
}

func (s *shortenerService) GetOriginalURL(ctx context.Context, shortCode string) (string, error) {

	if cached, err := s.cache.Get(ctx, shortCode); err == nil {
		return cached, nil
	}


	originalURL, err := s.repo.FindByShortCode(ctx, shortCode)
	if err != nil {
		return "", err
	}


	go func() {
		if err := s.cache.Set(context.Background(), shortCode, originalURL); err != nil {
			s.logger.Warn("Cache update failed", zap.Error(err))
		}
	}()

	return originalURL, nil
}
