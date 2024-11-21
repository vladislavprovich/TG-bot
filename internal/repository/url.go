package repository

type DeleteURLRequest struct {
	UserID      string `json:"user_id"`
	OriginalURL string `json:"original_url"`
}

type DeleteAllURLRequest struct {
	UserID string `json:"user_id"`
}

type GetListURLRequest struct {
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
