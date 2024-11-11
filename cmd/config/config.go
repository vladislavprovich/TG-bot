package config

type BotConfig struct {
	BotToken string
	//TODO add params
}

func LoadBotConfig() *BotConfig {
	return &BotConfig{
		BotToken: "7645002559:AAFHUKs3uI4rZ0zF3pf70fs9KHVPEYY3lyk",
	}
}

//func Config(key string) string {
//	// load .env file
//	err := godotenv.Load(".env")
//	if err != nil {
//		fmt.Print("Error loading .env file")
//	}
//	return os.Getenv(key)
//}
