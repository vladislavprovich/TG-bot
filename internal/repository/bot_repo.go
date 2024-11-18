package repository

import (
	"context"
	"database/sql"
	"github.com/sirupsen/logrus"
)

type BOTRepository interface {
	SaveURL(ctx context.Context, req SaveUrlRequest) error
	GetListURL(ctx context.Context, userID string) ([]*URLCombined, error)
	DeleteAllURL(ctx context.Context, userID string) error
	DeleteURL(ctx context.Context, originalURL string, userID string) error
}

type SaveUrlRequest struct {
	UserID string       `json:"user_id"`
	URL    *URLCombined `json:"url_combinate"`
}

type URLCombined struct {
	OriginalURL string `json:"original_url"`
	ShortURL    string `json:"short_url"`
}

type botRepository struct {
	db     DBExecutor
	logger logrus.Logger
}

func NewBotRepository(db DBExecutor, logger logrus.Logger) BOTRepository {
	return &botRepository{
		db:     db,
		logger: logger,
	}
}

// moke interface
type DBExecutor interface {
	ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error)
	QueryContext(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error)
}

func (r *botRepository) SaveURL(ctx context.Context, req SaveUrlRequest) error {
	return nil
}

func (r *botRepository) GetListURL(ctx context.Context, userID string) ([]*URLCombined, error) {
	return nil, nil
}

func (r *botRepository) DeleteAllURL(ctx context.Context, userID string) error {
	return nil
}

func (r *botRepository) DeleteURL(ctx context.Context, originalURL string, userID string) error {
	return nil
}
