package service

import (
	"context"
	"errors"
	"github.com/sirupsen/logrus"
	"github.com/vladislavprovich/TG-bot/internal/repository"
)

type UrlService interface {
	CreateShortUrl(ctx context.Context, req CreateShortUrlRequest) (CreateShortUrlResponse, error)
	GetListUrl(ctx context.Context, ID GetListRequest) ([]*GetListResponse, error)
	DeleteShortUrl(ctx context.Context, url DeleteShortUrl) error
	DeleteAllUrl(ctx context.Context, ID DeleteAllUrl) error
}

type urlService struct {
	repo   repository.URLRepository
	logger *logrus.Logger
}

func NewUrlService(repo repository.URLRepository, logger *logrus.Logger) UrlService {
	return &urlService{
		repo:   repo,
		logger: logger,
	}
}

func (s *urlService) CreateShortUrl(ctx context.Context, req CreateShortUrlRequest) (CreateShortUrlResponse, error) {
	if req.OriginalUrl == "" {
		return CreateShortUrlResponse{}, errors.New("origin url is empty")
	}

	existingUrls, err := s.repo.GetListURL(ctx, &repository.GetListURLRequest{UserID: req.UserID})
	if err != nil {
		s.logger.Errorf("failed to check existing URLs: %v", err)
		return CreateShortUrlResponse{}, err
	}

	for _, url := range existingUrls {
		if url.OriginalURL == req.OriginalUrl {
			return CreateShortUrlResponse{ShortUrl: url.ShortURL}, nil
		}
	}


	saveReq := &repository.SaveUrlRequest{
		UserID: req.UserID,
		URL: &repository.URLCombined{
			OriginalURL: req.OriginalUrl,
			ShortURL:    ??????,
		},
	}
	if err = s.repo.SaveURL(ctx, saveReq); err != nil {
		s.logger.Errorf("failed to save URL: %v", err)
		return CreateShortUrlResponse{}, err
	}

	return CreateShortUrlResponse{ShortUrl: ??????}, nil
}

func (s *urlService) GetListUrl(ctx context.Context, ID GetListRequest) ([]*GetListResponse, error) {
	urls, err := s.repo.GetListURL(ctx, &repository.GetListURLRequest{UserID: ID.ID})
	if err != nil {
		s.logger.Errorf("failed to get list of URLs: %v", err)
		return nil, err
	}

	var response []*GetListResponse
	for _, url := range urls {
		response = append(response, &GetListResponse{
			OriginalUrl: url.OriginalURL,
			ShortUrl:    url.ShortURL,
		})
	}
	return response, nil
}

func (s *urlService) DeleteShortUrl(ctx context.Context, url DeleteShortUrl) error {
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

func (s *urlService) DeleteAllUrl(ctx context.Context, ID DeleteAllUrl) error {
	err := s.repo.DeleteAllURL(ctx, &repository.DeleteAllURLRequest{
		UserID: ID.UserID,
	})
	if err != nil {
		s.logger.Errorf("failed to delete all URLs: %v", err)
	}
	return nil
}
