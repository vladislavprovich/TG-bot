package handler

import (
	"context"
	validation "github.com/go-ozzo/ozzo-validation/v4"
)

type BotConfig struct {
	BotToken string `envconfig:"BOT_TOKEN" `
	//TODO add params
}

func (b *BotConfig) ValidateWithContext(ctx context.Context) error {
	return validation.ValidateStructWithContext(ctx, &b,
		validation.Field(&b.BotToken, validation.Required),
	)
}

//func LoadBot() (*BotConfig, error) {
//	// load .env file
//	if err := godotenv.Load(); err != nil {
//		log.Println("No .env file found. Using system environment variables.")
//	}
//
//	botToken := os.Getenv("BOT_TOKEN")
//	if botToken == "" {
//		log.Fatal("BOT_TOKEN is not set in the environment variables")
//	}
//
//	return &BotConfig{
//		BotToken: botToken,
//	}, nil
//}
