package service

import (
	"github.com/vladislavprovich/TG-bot/internal/models"
	"github.com/vladislavprovich/TG-bot/pkg/shortener"
)

type converterToGetStats struct {
}

func NewConverterToGetStats() *converterToGetStats {
	return &converterToGetStats{}
}

func (c *converterToGetStats) ConverterToNewStats(
	req models.GetUrlStatusRequest,
) *shortener.GetShortURLStatsRequest {
	return &shortener.GetShortURLStatsRequest{
		ShortURL: req.ShortUrl,
	}
}
