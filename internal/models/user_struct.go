package models

type User struct {
	ID         string `json:"id"`
	TelegramID string `json:"telegram_id"`
	//Username    string `json:"username"`
	OriginalURL string `json:"original_url"`
	ShortURL    string `json:"short_url"`
}
