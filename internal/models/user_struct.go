package models

import "time"

type User struct {
	Action      string     `json:"action"`
	ID          string     `json:"id"`
	TelegramID  string     `json:"telegram_id"`
	OriginalURL string     `json:"original_url"`
	CustomUrl   string     `json:"custom_url"`
	ShortURL    string     `json:"short_url"`
	userAction  UserAction `json:"user_action"`
}

type UserAction struct {
	Action      string `json:"action"` // "rand_url" or "cust_url"
	OriginalURL string `json:"original_url"`
	CustomUrl   string `json:"custom_url"`
	ShortURL    string `json:"short_url"`
}

type CreateShortUrlRequest struct {
	TgID        int64   `json:"tg_id"`
	UserID      string  `json:"user_id"`
	OriginalUrl string  `json:"original_url"`
	CustomAlias *string `json:"custom_alias,omitempty"`
}

type CreateShortUrlResponse struct {
	ShortUrl string `json:"short_url"`
}

type DeleteShortUrl struct {
	TgID        int64  `json:"tg_id"`
	UserID      string `json:"user_id"`
	OriginalUrl string `json:"original_url"`
	ShortUrl    string `json:"short_url"`
}

type GetListRequest struct {
	TgID   int64  `json:"tg_id"`
	UserID string `json:"user_id"`
}

type GetListResponse struct {
	OriginalUrl string `json:"original_url"`
	ShortUrl    string `json:"short_url"`
}

type DeleteAllUrl struct {
	TgID   int64  `json:"tg_id"`
	UserID string `json:"user_id"`
}

type CreateNewUserRequest struct {
	TgID     int64  `json:"tg_id"`
	UserName string `json:"username"`
}

type GetUrlStatusRequest struct {
	UserID   string `json:"user_id"`
	TgID     int64  `json:"tg_id"`
	ShortUrl string `json:"short_url"`
}

type GetUrlStatusResponse struct {
	ShortUrl      string    `json:"short_url"`
	RedirectCount int       `json:"redirect_count"`
	CreatedAt     time.Time `json:"created_at"`
}
