package shortener

import "time"

type CreateShortUrlRequest struct {
	// user send original usr.
	URL string `json:"url"`
	// user send custom url (option).
	CustomAlias *string `json:"custom_alias,omitempty"`
}

type CreateShortUrlResponse struct {
	// returns short url.
	ShortURL string `json:"short_url"`
}

type GetShortUrlRequest struct {
	// user send short url.
	ShortURL string `json:"short_url"`
}

type GetShortUrlResponse struct {
	// the number of conversions for a short link.
	RedirectCount int `json:"redirect_count"`
	// time to created short url.
	CreatedAt time.Time `json:"created_at"`
}
