package service

import (
	"github.com/vladislavprovich/TG-bot/internal/repository"
)

type converterToTgID struct {
}

func NewConverterToTgID() *converterToTgID {
	return &converterToTgID{}
}

func (c *converterToTgID) converterToTgID(
	TgID int64,
) *repository.GetUserByTgIDRequest {
	return &repository.GetUserByTgIDRequest{
		TgID: TgID,
	}
}
