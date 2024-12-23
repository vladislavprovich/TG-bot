package service

import (
	"github.com/vladislavprovich/TG-bot/internal/models"
	"github.com/vladislavprovich/TG-bot/pkg/shortener"
)

type converterToGetStats struct {
}

func newConverterToGetStats() *converterToGetStats {
	return &converterToGetStats{}
}

func (c *converterToGetStats) ConverterToNewStats(
	req models.GetURLStatusRequest,
) *shortener.GetShortURLStatsRequest {
	return &shortener.GetShortURLStatsRequest{
		ShortURL: req.ShortURL,
	}
}
