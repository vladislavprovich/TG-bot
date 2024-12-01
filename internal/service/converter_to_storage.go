package service

import (
	"github.com/vladislavprovich/TG-bot/internal/models"
	"github.com/vladislavprovich/TG-bot/internal/repository"
)

type converterToStorage struct {
}

func NewConverterToStorage() *converterToStorage {
	return &converterToStorage{}
}

func (c *converterToStorage) ConvertToSaveUrlReq(
	req models.CreateShortUrlRequest,
	shortURL string,
	userID string,
) *repository.SaveUrlRequest {
	return &repository.SaveUrlRequest{
		UserID: userID,
		URL: &repository.URLCombined{
			OriginalURL: req.OriginalUrl,
			ShortURL:    shortURL,
		},
	}
}
