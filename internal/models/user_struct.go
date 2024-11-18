package models

type User struct {
	UserID string `json:"id"`
	//Username    string `json:"username"`
	OriginalURL string `json:"original_url"`
	ShortURL    string `json:"short_url"`
}
