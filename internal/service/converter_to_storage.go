package service

import (
	"github.com/vladislavprovich/TG-bot/internal/models"
	"github.com/vladislavprovich/TG-bot/internal/repository"
)

type converterToStorage struct {
}

func newConverterToStorage() *converterToStorage {
	return &converterToStorage{}
}

func (c *converterToStorage) ConvertToSaveURLReq(
	req models.CreateShortURLRequest,
	shortURL string,
	userID string,
) *repository.SaveURLRequest {
	return &repository.SaveURLRequest{
		UserID: userID,
		URL: &repository.URLCombined{
			OriginalURL: req.OriginalURL,
			ShortURL:    shortURL,
		},
	}
}
