package repository

import "time"

type DeleteURLRequest struct {
	TgID        int64  `json:"tg_id"`
	UserID      string `json:"user_id"`
	OriginalURL string `json:"original_url"`
	ShortURL    string `json:"short_url"`
}

type DeleteAllURLRequest struct {
	TgID   int64  `json:"tg_id"`
	UserID string `json:"user_id"`
}

type GetListURLRequest struct {
	TgID   int64  `json:"tg_id"`
	UserID string `json:"user_id"`
}

type SaveUrlRequest struct {
	UserID string       `json:"user_id"`
	URL    *URLCombined `json:"url_combined"`
}

type URLCombined struct {
	OriginalURL string `json:"original_url"`
	ShortURL    string `json:"short_url"`
}

type GetUrlStatsRequest struct {
	TgID     int64  `json:"tg_id"`
	UserID   string `json:"user_id"`
	ShortURL string `json:"short_url"`
}

type GetUrlStatsResponse struct {
	ShortURL      string    `json:"short_url"`
	RedirectCount int       `json:"redirect_count"`
	CreatedAt     time.Time `json:"created_at"`
}
