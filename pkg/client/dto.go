package client

type Request struct {
	URL         string  `json:"url"`
	CustomAlias *string `json:"custom_alias,omitempty"`
}

type Response struct {
	ShortURL string `json:"short_url"`
}
