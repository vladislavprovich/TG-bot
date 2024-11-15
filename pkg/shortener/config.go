package shortener

import (
	"context"
	validation "github.com/go-ozzo/ozzo-validation/v4"
)

type Config struct {
	Port       string `envconfig:"SERVER_PORT" default:"8080"`
	BaseURL    string `envconfig:"BASE_URL" default:"http://localhost:8080/shorten"`
	BaseGetURL string `envconfig:"BASE_GET_URL" default:"http://localhost:8080/:url/shorten"`
	// todo add if needed
}

func (c Config) ValidateWithContext(ctx context.Context) error {
	return validation.ValidateStructWithContext(ctx, &c,
		validation.Field(&c.Port, validation.Required),
		validation.Field(&c.BaseURL, validation.Required),
		validation.Field(&c.BaseGetURL, validation.Required),
	)
}
