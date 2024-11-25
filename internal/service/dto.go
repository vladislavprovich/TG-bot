package service

type CreateShortUrlRequest struct {
	UserID      string `json:"user_id"`
	OriginalUrl string `json:"original_url"`
}

type CreateShortUrlResponse struct {
	ShortUrl string `json:"short_url"`
}

type DeleteShortUrl struct {
	UserID      string `json:"user_id"`
	OriginalUrl string `json:"original_url"`
}

type GetListRequest struct {
	ID string `json:"id"`
}

type GetListResponse struct {
	OriginalUrl string `json:"original_url"`
	ShortUrl    string `json:"short_url"`
}

type DeleteAllUrl struct {
	UserID string `json:"user_id"`
}
