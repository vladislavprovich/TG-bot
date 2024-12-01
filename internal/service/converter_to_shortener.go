package service

import (
	"github.com/vladislavprovich/TG-bot/internal/models"
	"github.com/vladislavprovich/TG-bot/pkg/shortener"
)

type converterToShortener struct {
}

func NewConverterToShortener() *converterToShortener {
	return &converterToShortener{}
}

func (c *converterToShortener) ConvertToCreateShortURLRequest(
	req models.CreateShortUrlRequest,
) *shortener.CreateShortURLRequest {
	return &shortener.CreateShortURLRequest{
		URL:         req.OriginalUrl,
		CustomAlias: req.CustomAlias,
	}
}
