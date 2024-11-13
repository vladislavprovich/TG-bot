package config

import (
	"github.com/joho/godotenv"
	"log"
	"os"
)

type BotConfig struct {
	BotToken string
	//TODO add params
}

func LoadBotConfig() *BotConfig {
	// load .env file
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found. Using system environment variables.")
	}

	botToken := os.Getenv("BOT_TOKEN")
	if botToken == "" {
		log.Fatal("BOT_TOKEN is not set in the environment variables")
	}

	return &BotConfig{
		BotToken: botToken,
	}
}
