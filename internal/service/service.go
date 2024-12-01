package service

import (
	"context"
	"errors"
	"github.com/sirupsen/logrus"
	"github.com/vladislavprovich/TG-bot/internal/models"
	"github.com/vladislavprovich/TG-bot/internal/repository"
	"github.com/vladislavprovich/TG-bot/pkg/shortener"
)

type UrlService interface {
	CreateShortUrl(ctx context.Context, req models.CreateShortUrlRequest) (models.CreateShortUrlResponse, error)
	GetListUrl(ctx context.Context, ID models.GetListRequest) ([]*models.GetListResponse, error)
	DeleteShortUrl(ctx context.Context, url models.DeleteShortUrl) error
	DeleteAllUrl(ctx context.Context, ID models.DeleteAllUrl) error
}

type (
	Service struct {
		client             shortener.Client
		repo               repository.URLRepository
		logger             *logrus.Logger
		convertToShortener *converterToShortener
		convertToStorage   *converterToStorage
	}

	Params struct {
		repo   repository.URLRepository
		logger *logrus.Logger
		client shortener.Client
	}
)

func NewService(params Params) UrlService {
	return &Service{
		client:             params.client,
		repo:               params.repo,
		logger:             params.logger,
		convertToShortener: NewConverterToShortener(),
		convertToStorage:   NewConverterToStorage(),
	}
}

func (s *Service) CreateShortUrl(ctx context.Context, req models.CreateShortUrlRequest) (models.CreateShortUrlResponse, error) {
	if req.OriginalUrl == "" {
		return models.CreateShortUrlResponse{}, errors.New("origin url is empty")
	}

	existingUrls, err := s.repo.GetListURL(ctx, &repository.GetListURLRequest{UserID: req.UserID})
	if err != nil {
		s.logger.Errorf("failed to check existing URLs: %v", err)
		return models.CreateShortUrlResponse{}, err
	}

	for _, url := range existingUrls {
		if url.OriginalURL == req.OriginalUrl {
			return models.CreateShortUrlResponse{ShortUrl: url.ShortURL}, nil
		}
	}

	convertedShortUrlReq := s.convertToShortener.ConvertToCreateShortURLRequest(req)

	shortUrlResp, err := s.client.CreateShortUrl(ctx, convertedShortUrlReq)

	if err != nil {
		s.logger.Errorf("failed to create short url: %v", err)
		return models.CreateShortUrlResponse{}, err
	}

	saveReq := s.convertToStorage.ConvertToSaveUrlReq(req, shortUrlResp.ShortURL)

	if err = s.repo.SaveURL(ctx, saveReq); err != nil {
		s.logger.Errorf("failed to save URL: %v", err)
		return models.CreateShortUrlResponse{}, err
	}

	return models.CreateShortUrlResponse{ShortUrl: shortUrlResp.ShortURL}, nil
}

func (s *Service) GetListUrl(ctx context.Context, ID models.GetListRequest) ([]*models.GetListResponse, error) {
	urls, err := s.repo.GetListURL(ctx, &repository.GetListURLRequest{UserID: ID.TgID})
	if err != nil {
		s.logger.Errorf("failed to get list of URLs: %v", err)
		return nil, err
	}

	var response []*models.GetListResponse
	for _, url := range urls {
		response = append(response, &models.GetListResponse{
			OriginalUrl: url.OriginalURL,
			ShortUrl:    url.ShortURL,
		})
	}
	return response, nil
}

func (s *Service) DeleteShortUrl(ctx context.Context, url models.DeleteShortUrl) error {
	err := s.repo.DeleteURL(ctx, &repository.DeleteURLRequest{
		UserID:      url.UserID,
		OriginalURL: url.OriginalUrl,
	})
	if err != nil {
		s.logger.Errorf("failed to delete short URL: %v", err)
		return err
	}
	return nil
}

func (s *Service) DeleteAllUrl(ctx context.Context, ID models.DeleteAllUrl) error {
	err := s.repo.DeleteAllURL(ctx, &repository.DeleteAllURLRequest{
		UserID: ID.TgID,
	})
	if err != nil {
		s.logger.Errorf("failed to delete all URLs: %v", err)
	}
	return nil
}
