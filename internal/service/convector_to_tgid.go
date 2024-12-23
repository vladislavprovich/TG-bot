package service

import (
	"github.com/vladislavprovich/TG-bot/internal/repository"
)

type converterToTgID struct {
}

func newConverterToTgID() *converterToTgID {
	return &converterToTgID{}
}

func (c *converterToTgID) converterToTgID(
	tgID int64,
) *repository.GetUserByTgIDRequest {
	return &repository.GetUserByTgIDRequest{
		TgID: tgID,
	}
}
