package models

type User struct {
	ID          int    `json:"id"`
	OriginalURL string `json:"original_url"`
	ShortURL    string `json:"short_url"`
}

type Request struct {
	OriginalURL string `json:"original_url"`
	CustomURL   string `json:"custom_url,omitempty"`
}

type Response struct {
	ShortURL string `json:"short_url"`
}
