package service

import (
	"context"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/vladislavprovich/TG-bot/internal/models"
	"github.com/vladislavprovich/TG-bot/internal/repository"
)

type UrlService interface {
	SaveURL(ctx context.Context, username string, originalURL string) error
	GetList(ctx context.Context, userID string) (models.User, error)
}

type urlService struct {
	bot  *tgbotapi.BotAPI
	repo repository.UrlRepo
}

func NewURlService(bot *tgbotapi.BotAPI, repo repository.UrlRepo) UrlService {
	return &urlService{
		bot:  bot,
		repo: repo,
	}
}

func (u *urlService) SaveURL(ctx context.Context, username, originalURL string) error {
	user, err := u.repo.GetList(ctx, username)
	if err != nil {
		return err
	}

	if user.ID == "" {
		user = models.User{
			Username: username,
		}
		userID, err = u.repo.CreateUser(ctx, user)
		if err != nil {
			return err
		}
		user.ID = userID //error
	}

	url := models.User{
		ID:          user.ID,
		OriginalURL: originalURL,
		ShortURL:    shorturl, //error shorturl
	}
	return u.repo.SaveURL(ctx, url)
}

func (u *urlService) GetList(ctx context.Context, userID string) (models.User, error) {
	return u.repo.GetList(ctx, userID)
}
