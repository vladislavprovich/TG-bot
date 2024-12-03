package main

import (
	"context"
	"github.com/sirupsen/logrus"
	"github.com/vladislavprovich/TG-bot/internal/handler"
	"github.com/vladislavprovich/TG-bot/internal/keyboard"
	"github.com/vladislavprovich/TG-bot/internal/repository"
	"github.com/vladislavprovich/TG-bot/internal/repository/postgres"
	"github.com/vladislavprovich/TG-bot/internal/service"
	"github.com/vladislavprovich/TG-bot/pkg/shortener"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	ctx := context.Background()
	cfg, err := LoadConfig(ctx)
	if err != nil {
		log.Fatalf("Error loading config: %v", err)
	}

	bot := handler.BotInit(cfg.Bot)

	logger := logrus.New()

	db, err := postgres.PrepareConnection(ctx, cfg.Database, logger)
	if err != nil {
		logger.Fatal(err)
	}

	httpClient := &http.Client{
		Timeout: 10 * time.Second,
	}

	client := shortener.NewBasicClient(cfg.Client, httpClient, logger)

	UrlRepo := repository.NewBotRepository(db, logger)
	UserRepo := repository.NewUserRepository(db, logger)

	params := service.Params{
		UrlRepo,
		UserRepo,
		logger,
		&client,
	}

	urlService := service.NewService(params)

	buttonHandler := keyboard.NewHandleButtons(urlService, logger)
	messageHandler := keyboard.NewMessageHandler(urlService, logger)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go handler.ProcessUpdates(ctx, bot, buttonHandler, messageHandler, logger)

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
	<-stop

	logger.Info("Shutting down gracefully...")

	//for update := range updates {
	//	if update.CallbackQuery != nil {
	//		keyboard.HandleCallbackQuery(bot, update, urlService)
	//	} else if update.Message != nil {
	//		keyboard.HandleMessage(bot, update, urlService)
	//	}
	//}

	//for update := range updates {
	//	if update.Message != nil && update.Message.Text == "/start" { //start in const
	//		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Вітаю! Оберіть опцію:")
	//		msg.ReplyMarkup = keyboard.MainMenu()
	//		bot.Send(msg)
	//	}
	//
	//	if update.CallbackQuery != nil {
	//		var responseText string
	//		var replyMarkup tgbotapi.InlineKeyboardMarkup
	//
	//		switch update.CallbackQuery.Data {
	//		case "create_short_url":
	//			responseText = "Виберіть опцію"
	//			replyMarkup = keyboard.CreateURL()
	//		case "list_short_urls":
	//			responseText = "Ось список всіх ваших скорочених URL:"
	//			replyMarkup = keyboard.BackMenu()
	//		case "show_url_stats":
	//			responseText = "URL статистика:"
	//			replyMarkup = keyboard.BackMenu()
	//		case "settings":
	//			responseText = "Налаштування: Ви можете очистити історію URL."
	//			replyMarkup = keyboard.ClearAndBack()
	//		case "back_to_main":
	//			responseText = "Вітаю! Оберіть опцію:"
	//			replyMarkup = keyboard.MainMenu()
	//		case "clear_history":
	//			responseText = "Історію видалено.\nВітаю! Оберіть опцію:"
	//			replyMarkup = keyboard.MainMenu()
	//		case "rand_url":
	//			responseText = "Твоя коротка URL:"
	//			replyMarkup = keyboard.MainMenu()
	//		case "cust_url":
	//			responseText = "Напиши свою custom URL:"
	//			replyMarkup = keyboard.MainMenu()
	//		default:
	//			responseText = "Невідома команда."
	//			replyMarkup = keyboard.MainMenu()
	//		}
	//
	//		// Update msg
	//		editMsg := tgbotapi.NewEditMessageTextAndMarkup(
	//			update.CallbackQuery.Message.Chat.ID,
	//			update.CallbackQuery.Message.MessageID,
	//			responseText,
	//			replyMarkup,
	//		)
	//		_, err := bot.Send(editMsg)
	//		if err != nil {
	//			panic(err)
	//		}
	//
	//		callback := tgbotapi.CallbackConfig{
	//			CallbackQueryID: update.CallbackQuery.ID,
	//		}
	//		_, err = bot.Request(callback)
	//		if err != nil {
	//			panic(err)
	//		}
	//
	//	}
	//}

}
