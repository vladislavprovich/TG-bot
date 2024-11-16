package main

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/vladislavprovich/TG-bot/internal/handler"
	"github.com/vladislavprovich/TG-bot/internal/keyboard"
)

func main() {
	cfg := handler.LoadBotConfig()
	bot := handler.BotInit(*cfg)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message != nil && update.Message.Text == "/start" {
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Вітаю! Оберіть опцію:")
			msg.ReplyMarkup = keyboard.MainMenu()
			bot.Send(msg)
		}

		if update.CallbackQuery != nil {
			var responseText string
			var replyMarkup tgbotapi.InlineKeyboardMarkup

			switch update.CallbackQuery.Data {
			case "create_short_url":
				responseText = "Виберіть опцію"
				replyMarkup = keyboard.CreateURL()
			case "list_short_urls":
				responseText = "Ось список всіх ваших скорочених URL:"
				replyMarkup = keyboard.BackMenu()
			case "show_url_stats":
				responseText = "URL статистика:"
				replyMarkup = keyboard.BackMenu()
			case "settings":
				responseText = "Налаштування: Ви можете очистити історію URL."
				replyMarkup = keyboard.ClearAndBack()
			case "back_to_main":
				responseText = "Вітаю! Оберіть опцію:"
				replyMarkup = keyboard.MainMenu()
			case "clear_history":
				responseText = "Історію видалено.\nВітаю! Оберіть опцію:"
				replyMarkup = keyboard.MainMenu()
			case "rand_url":
				responseText = "Твоя коротка URL:"
				replyMarkup = keyboard.MainMenu()
			case "cust_url":
				responseText = "Напиши свою custom URL:"
				replyMarkup = keyboard.MainMenu()
			default:
				responseText = "Невідома команда."
				replyMarkup = keyboard.MainMenu()
			}

			// Update msg
			editMsg := tgbotapi.NewEditMessageTextAndMarkup(
				update.CallbackQuery.Message.Chat.ID,
				update.CallbackQuery.Message.MessageID,
				responseText,
				replyMarkup,
			)
			_, err := bot.Send(editMsg)
			if err != nil {
				panic(err)
			}

			callback := tgbotapi.CallbackConfig{
				CallbackQueryID: update.CallbackQuery.ID,
			}
			_, err = bot.Request(callback)
			if err != nil {
				panic(err)
			}

		}
	}
}
