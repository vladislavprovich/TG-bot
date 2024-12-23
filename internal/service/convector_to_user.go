package service

import (
	"github.com/vladislavprovich/TG-bot/internal/models"
	"github.com/vladislavprovich/TG-bot/internal/repository"
)

type converterToUser struct {
}

func newConverterToUser() *converterToUser {
	return &converterToUser{}
}

func (c *converterToUser) converterToNewUser(
	req models.CreateNewUserRequest,
	userID string,
) *repository.SaveUserRequest {
	return &repository.SaveUserRequest{
		TgID:     req.TgID,
		UserID:   userID,
		UserName: req.UserName,
	}
}
