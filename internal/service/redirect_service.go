package service

import (
	"url-shortener/internal/logger"
	"url-shortener/internal/repository"
)

type RedirectService struct {
	repo   *repository.RedisRepository
	logger *logger.Logger
}

func NewRedirectService(redisRepo *repository.RedisRepository, logger *logger.Logger) *RedirectService {
	return &RedirectService{
		repo:   redisRepo,
		logger: logger,
	}
}

func (svc *RedirectService) Redirect(shortURL string) string {
	return svc.repo.Retrieve(shortURL)
}
