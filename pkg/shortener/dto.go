package shortener

import "time"

type CreateShortURLRequest struct {
	// user send original url.
	URL string `json:"url"`
	// user send custom url (optional).
	CustomAlias *string `json:"custom_alias,omitempty"`
}

type CreateShortURLResponse struct {
	// return created short url.
	ShortURL string `json:"short_url"`
}

type GetShortURLStatsRequest struct {
	// user send short url.
	ShortURL string `json:"short_url"`
}

type GetShortURLStatsResponse struct {
	// the number of redirect for a short link
	RedirectCount int `json:"redirect_count"`
	// short url creation date.
	CreatedAt time.Time `json:"created_at"`
}
