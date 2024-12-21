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
	"github.com/vladislavprovich/TG-bot/pkg"
	"github.com/vladislavprovich/TG-bot/pkg/shortener"
)

type UrlService interface {
	CreateShortUrl(ctx context.Context, req models.CreateShortUrlRequest) (models.CreateShortUrlResponse, error)
	GetListUrl(ctx context.Context, ID models.GetListRequest) ([]*models.GetListResponse, error)
	DeleteShortUrl(ctx context.Context, url models.DeleteShortUrl) error
	DeleteAllUrl(ctx context.Context, ID models.DeleteAllUrl) error
	CreateUserByTgID(ctx context.Context, req models.CreateNewUserRequest) (string, error)
	GetUrlStatus(ctx context.Context, req models.GetUrlStatusRequest) ([]*models.GetUrlStatusResponse, error)
}

type (
	Service struct {
		client              shortener.Client
		repo                repository.URLRepository
		repoUser            repository.UserRepository
		logger              *logrus.Logger
		convertToShortener  *converterToShortener
		convertToStorage    *converterToStorage
		convertToUser       *converterToUser
		convertToTgID       *converterToTgID
		converterToGetStats *converterToGetStats
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
		client:              params.Client,
		repo:                params.Repo,
		repoUser:            params.RepoUser,
		logger:              params.Logger,
		convertToShortener:  NewConverterToShortener(),
		convertToStorage:    NewConverterToStorage(),
		convertToUser:       NewConverterToUser(),
		convertToTgID:       NewConverterToTgID(),
		converterToGetStats: NewConverterToGetStats(),
	}
}

func (s *Service) CreateShortUrl(ctx context.Context, req models.CreateShortUrlRequest) (models.CreateShortUrlResponse, error) {
	if req.OriginalUrl == "" {
		return models.CreateShortUrlResponse{}, errors.New("origin url is empty")
	}
	if req.UserID == "" {
		newReq := s.convertToTgID.converterToTgID(req.TgID)
		userIDresp, err := s.repoUser.GetUserByTgID(ctx, newReq)
		if err != nil {
			return models.CreateShortUrlResponse{}, err
		}
		s.logger.Infof("user id check %s", userIDresp.User.UserID)
		req.UserID = userIDresp.User.UserID
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

	saveReq := s.convertToStorage.ConvertToSaveUrlReq(req, shortUrlResp.ShortURL, req.UserID)

	if err = s.repo.SaveURL(ctx, saveReq); err != nil {
		s.logger.Errorf("failed to save URL: %v", err)
		return models.CreateShortUrlResponse{}, err
	}

	return models.CreateShortUrlResponse{ShortUrl: shortUrlResp.ShortURL}, nil
}

func (s *Service) GetListUrl(ctx context.Context, ID models.GetListRequest) ([]*models.GetListResponse, error) {
	userID, err := s.repoUser.GetUserByTgID(ctx, &repository.GetUserByTgIDRequest{TgID: ID.TgID})
	userID.User.UserID = ID.UserID

	urlsRes, err := s.repo.GetListURL(ctx, &repository.GetListURLRequest{UserID: ID.UserID})
	if err != nil {
		s.logger.Errorf("failed to get list of URLs: %v", err)
		return nil, err
	}

	var response []*models.GetListResponse
	for _, url := range urlsRes {
		response = append(response, &models.GetListResponse{
			OriginalUrl: url.OriginalURL,
			ShortUrl:    url.ShortURL,
		})
	}
	return response, nil
}

func (s *Service) DeleteShortUrl(ctx context.Context, url models.DeleteShortUrl) error {
	userID, err := s.repoUser.GetUserByTgID(ctx, &repository.GetUserByTgIDRequest{TgID: url.TgID})
	userID.User.UserID = url.UserID
	err = s.repo.DeleteURL(ctx, &repository.DeleteURLRequest{
		UserID:      url.UserID,
		OriginalURL: url.OriginalUrl,
		ShortURL:    url.ShortUrl,
	})
	if err != nil {
		s.logger.Errorf("failed to delete short URL: %v", err)
		return err
	}
	return nil
}

func (s *Service) DeleteAllUrl(ctx context.Context, ID models.DeleteAllUrl) error {
	userID, err := s.repoUser.GetUserByTgID(ctx, &repository.GetUserByTgIDRequest{TgID: ID.TgID})
	userID.User.UserID = ID.UserID

	err = s.repo.DeleteAllURL(ctx, &repository.DeleteAllURLRequest{
		TgID:   ID.TgID,
		UserID: ID.UserID,
	})
	if err != nil {
		s.logger.Errorf("failed to delete all URLs: %v", err)
	}
	return nil
}

func (s *Service) CreateUserByTgID(ctx context.Context, req models.CreateNewUserRequest) (string, error) {
	reqNew := s.convertToTgID.converterToTgID(req.TgID)

	userResp, err := s.repoUser.GetUserByTgID(ctx, reqNew)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return "", fmt.Errorf("DataBase Error")
		s.logger.Errorf("error checking user existence: %w", err)
	}

	if userResp != nil {
		s.logger.Errorf("INFO ERROR", userResp.User.UserID)
		return userResp.User.UserID, nil
	}
	userID := uuid.New().String()

	saveUserReq := s.convertToUser.converterToNewUser(req, userID)

	err = s.repoUser.SaveUser(ctx, saveUserReq)
	if err != nil {
		s.logger.Errorf("failed to save user: %v", err)
		return "", err
	}

	return saveUserReq.UserID, nil
}

func (s *Service) GetUrlStatus(ctx context.Context, req models.GetUrlStatusRequest) ([]*models.GetUrlStatusResponse, error) {
	userID, err := s.repoUser.GetUserByTgID(ctx, &repository.GetUserByTgIDRequest{TgID: req.TgID})
	if err != nil {
		s.logger.Errorf("failed to get user ID for TG ID %d: %v", req.TgID, err)
		return nil, fmt.Errorf("user ID check error: %w", err)
	}
	req.UserID = userID.User.UserID

	urls, err := s.repo.GetListURL(ctx, &repository.GetListURLRequest{TgID: req.TgID, UserID: req.UserID})
	if err != nil {
		s.logger.Errorf("failed to get list of URLs for user ID %s: %v", req.UserID, err)
		return nil, err
	}

	var responses []*models.GetUrlStatusResponse

	for _, url := range urls {
		shortInfoUrls := pkg.ShortInfo(url.ShortURL)
		statsReq := &shortener.GetShortURLStatsRequest{
			ShortURL: shortInfoUrls,
		}

		stats, err := s.client.GetStatsUrl(ctx, statsReq)
		if err != nil {
			s.logger.Errorf("failed to get stats for short URL %s: %v", shortInfoUrls, err)
			return nil, err
		}

		responses = append(responses, &models.GetUrlStatusResponse{
			ShortUrl:      shortInfoUrls,
			RedirectCount: stats.RedirectCount,
			CreatedAt:     stats.CreatedAt,
		})
	}

	if len(responses) == 0 {
		s.logger.Warnf("no valid stats found for user ID %s", req.UserID)
		return nil, fmt.Errorf("no valid stats found for user URLs")
	}

	return responses, nil
}
