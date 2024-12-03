package service

import (
	"github.com/vladislavprovich/TG-bot/internal/models"
	"github.com/vladislavprovich/TG-bot/internal/repository"
)

type converterToUser struct {
}

func NewConverterToUser() *converterToUser {
	return &converterToUser{}
}

func (c *converterToUser) converterToNewUser(
	req models.CreateNewUserRequest,
) *repository.SaveUserRequest {
	return &repository.SaveUserRequest{
		TgID:     req.TgID,
		UserID:   req.UserID,
		UserName: req.UserName,
	}
}
