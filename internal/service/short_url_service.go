package service

import (
	"url-shortener/internal/logger"
	"url-shortener/internal/repository"
)

type URLShortener struct {
	hashService *HashService
	repo        *repository.RedisRepository
	logger      *logger.Logger
	baseUrl     string
}

func NewURLShortenerService(hashService *HashService, redisRepo *repository.RedisRepository, logger *logger.Logger, baseUrl string) *URLShortener {
	return &URLShortener{
		hashService: hashService,
		repo:        redisRepo,
		logger:      logger,
		baseUrl:     baseUrl,
	}
}

func (svc *URLShortener) Create(req *Request) (Response, error) {
	hash := svc.createHash()

	err := svc.store(hash, req.URL)
	if err != nil {
		return Response{}, err
	}

	return Response{ShortURL: svc.baseUrl + "/" + hash}, nil
}

func (svc *URLShortener) createHash() string {
	return svc.hashService.getHash()
}

func (svc *URLShortener) store(hash, originUrl string) error {
	err := svc.repo.Store(hash, originUrl)
	if err != nil {
		return err
	}

	return nil
}

type Request struct {
	URL string `json:"url"`
}

type Response struct {
	ShortURL string `json:"shortURL"`
}
