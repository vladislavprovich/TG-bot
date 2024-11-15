package shortener

import "time"

type CreateShortURLRequest struct {
	// user send original url.
	URL string `json:"url"`
	// user send custom url (optional).
	CustomAlias *string `json:"custom_alias,omitempty"`
}

type CreateShortURLResponse struct {
	// returns short url.
	ShortURL string `json:"short_url"`
}

type GetShortURLStatsRequest struct {
	// user send short url.
	ShortURL string `json:"short_url"`
}

type GetShortURLStatsResponse struct {
	// the number of conversions for a short link.
	RedirectCount int `json:"redirect_count"`
	// time to created short url.
	CreatedAt time.Time `json:"created_at"`
}
