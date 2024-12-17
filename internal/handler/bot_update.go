package handler

import (
	"context"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/sirupsen/logrus"
	"github.com/vladislavprovich/TG-bot/internal/keyboard"
)

const TimeUpdate = 60

func ProcessUpdates(ctx context.Context, bot *tgbotapi.BotAPI, buttonHandler *keyboard.HandleButtons, messageHandler *keyboard.HandleButtons, logger *logrus.Logger) {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = TimeUpdate

	updates := bot.GetUpdatesChan(u)

	for {
		select {
		case update := <-updates:
			if update.CallbackQuery != nil {
				go buttonHandler.HandleCallbackQuery(ctx, bot, update)
			} else if update.Message != nil {
				go messageHandler.HandleMessage(ctx, bot, update)
			}
		case <-ctx.Done():
			return
		}
	}
}
