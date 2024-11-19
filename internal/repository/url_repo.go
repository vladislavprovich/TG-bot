package repository

import (
	"context"
	"github.com/sirupsen/logrus"
	"github.com/vladislavprovich/TG-bot/internal/repository/test"
)

type URLRepository interface {
	SaveURL(ctx context.Context, req *SaveUrlRequest) error
	GetListURL(ctx context.Context, req *GetListURLRequest) ([]*URLCombined, error)
	DeleteAllURL(ctx context.Context, req *DeleteAllURLRequest) error
	DeleteURL(ctx context.Context, req *DeleteURLRequest) error
}

type urlRepository struct {
	db     test.DBExecutor
	logger logrus.Logger
}

func NewBotRepository(db test.DBExecutor, logger logrus.Logger) URLRepository {
	return &urlRepository{
		db:     db,
		logger: logger,
	}
}

func (r *urlRepository) SaveURL(ctx context.Context, req *SaveUrlRequest) error {
	return nil
}

func (r *urlRepository) GetListURL(ctx context.Context, req *GetListURLRequest) ([]*URLCombined, error) {
	return nil, nil
}

func (r *urlRepository) DeleteAllURL(ctx context.Context, req *DeleteAllURLRequest) error {
	return nil
}

func (r *urlRepository) DeleteURL(ctx context.Context, req *DeleteURLRequest) error {
	return nil
}
