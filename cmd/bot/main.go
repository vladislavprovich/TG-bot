package main

import (
	"context"
	"github.com/vladislavprovich/TG-bot/internal/handler"
	"github.com/vladislavprovich/TG-bot/internal/keyboard"
	"github.com/vladislavprovich/TG-bot/internal/repository"
	"github.com/vladislavprovich/TG-bot/internal/repository/postgres"
	"github.com/vladislavprovich/TG-bot/internal/service"
	"github.com/vladislavprovich/TG-bot/pkg/logger"
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
	logger := logger.NewLogger(cfg.Logger)
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
	//todo
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
	<-stop

	logger.Info("Shutting down gracefully...")
}
