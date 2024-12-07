package service

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/google/uuid"
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
	CreateUserByTgID(ctx context.Context, req models.CreateNewUserRequest) error
}

type (
	Service struct {
		client             shortener.Client
		repo               repository.URLRepository
		repoUser           repository.UserRepository
		logger             *logrus.Logger
		convertToShortener *converterToShortener
		convertToStorage   *converterToStorage
		convertToUser      *converterToUser
	}

	Params struct {
		Repo     repository.URLRepository
		RepoUser repository.UserRepository
		Logger   *logrus.Logger
		Client   shortener.Client
	}
)

func NewService(params Params) UrlService {
	return &Service{
		client:             params.Client,
		repo:               params.Repo,
		repoUser:           params.RepoUser,
		logger:             params.Logger,
		convertToShortener: NewConverterToShortener(),
		convertToStorage:   NewConverterToStorage(),
		convertToUser:      NewConverterToUser(),
	}
}

func (s *Service) CreateShortUrl(ctx context.Context, req models.CreateShortUrlRequest) (models.CreateShortUrlResponse, error) {
	if req.OriginalUrl == "" {
		return models.CreateShortUrlResponse{}, errors.New("origin url is empty")
	}

	existingUrls, err := s.repo.GetListURL(ctx, &repository.GetListURLRequest{TgID: req.TgID})
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

	saveReq := s.convertToStorage.ConvertToSaveUrlReq(req, shortUrlResp.ShortURL, req.UserID)

	if err = s.repo.SaveURL(ctx, saveReq); err != nil {
		s.logger.Errorf("failed to save URL: %v", err)
		return models.CreateShortUrlResponse{}, err
	}

	return models.CreateShortUrlResponse{ShortUrl: shortUrlResp.ShortURL}, nil
}

func (s *Service) GetListUrl(ctx context.Context, ID models.GetListRequest) ([]*models.GetListResponse, error) {
	urls, err := s.repo.GetListURL(ctx, &repository.GetListURLRequest{TgID: ID.TgID})
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
		TgID:        url.TgID,
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
		TgID: ID.TgID,
	})
	if err != nil {
		s.logger.Errorf("failed to delete all URLs: %v", err)
	}
	return nil
}

func (s *Service) CreateUserByTgID(ctx context.Context, req models.CreateNewUserRequest) error {

	User, err := s.repoUser.GetUserByTgID(ctx, &repository.GetUserByTgIDRequest{TgID: req.TgID})
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return fmt.Errorf("error checking user existence: %w", err)
	}

	if User != nil {
		return errors.New("user already exists")
	}

	userID := uuid.New().String()

	saveUserReq := s.convertToUser.converterToNewUser(req, userID)

	err = s.repoUser.SaveUser(ctx, saveUserReq)
	if err != nil {
		s.logger.Errorf("failed to save user: %v", err)
	}

	return nil
}
