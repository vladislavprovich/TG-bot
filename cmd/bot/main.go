package main

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/vladislavprovich/TG-bot/cmd/config"
	"github.com/vladislavprovich/TG-bot/internal/handler"
)

func main() {
	cfg := config.LoadBotConfig()

	bot, err := tgbotapi.NewBotAPI(cfg.BotToken)
	if err != nil {
		panic(err)
	}

	bot.Debug = true

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message != nil && update.Message.Text == "/start" {
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Вітаю! Оберіть опцію:")
			msg.ReplyMarkup = handler.MainMenu()
			bot.Send(msg)
		}

		if update.CallbackQuery != nil {
			var responseText string
			var replyMarkup tgbotapi.InlineKeyboardMarkup

			switch update.CallbackQuery.Data {
			case "create_short_url":
				responseText = "Виберіть опцію"
				replyMarkup = handler.CreateURL()
			case "list_short_urls":
				responseText = "Ось список всіх ваших скорочених URL:"
				replyMarkup = handler.BackMenu()
			case "show_url_stats":
				responseText = "URL статистика:"
				replyMarkup = handler.BackMenu()
			case "settings":
				responseText = "Налаштування: Ви можете очистити історію URL."
				replyMarkup = handler.ClearAndBack()
			case "back_to_main":
				responseText = "Вітаю! Оберіть опцію:"
				replyMarkup = handler.MainMenu()
			case "clear_history":
				responseText = "Історію видалено.\nВітаю! Оберіть опцію:"
				replyMarkup = handler.MainMenu()
			case "rand_url":
				responseText = "Твоя коротка URL:"
				replyMarkup = handler.MainMenu()
			case "cust_url":
				responseText = "Твоя custom URL:"
				replyMarkup = handler.MainMenu()
			default:
				responseText = "Невідома команда."
				replyMarkup = handler.MainMenu()
			}

			// Оновлюємо повідомлення
			editMsg := tgbotapi.NewEditMessageTextAndMarkup(
				update.CallbackQuery.Message.Chat.ID,
				update.CallbackQuery.Message.MessageID,
				responseText,
				replyMarkup,
			)
			bot.Send(editMsg)

			callback := tgbotapi.CallbackConfig{
				CallbackQueryID: update.CallbackQuery.ID,
			}
			bot.Request(callback)
		}
	}
}
