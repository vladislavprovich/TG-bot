package repository

import (
	"context"
	"database/sql"
	"github.com/sirupsen/logrus"
)

type URLRepository interface {
	SaveURL(ctx context.Context, req *SaveUrlRequest) error
	GetListURL(ctx context.Context, req *GetListURLRequest) ([]*URLCombined, error)
	DeleteAllURL(ctx context.Context, req *DeleteAllURLRequest) error
	DeleteURL(ctx context.Context, req *DeleteURLRequest) error
}

type urlRepository struct {
	db     *sql.DB
	logger *logrus.Logger
}

func NewBotRepository(db *sql.DB, logger *logrus.Logger) URLRepository {
	return &urlRepository{
		db:     db,
		logger: logger,
	}
}

func (r *urlRepository) SaveURL(ctx context.Context, req *SaveUrlRequest) error {
	_, err := r.db.ExecContext(ctx, `INSERT INTO urls (user_id, original_url, short_url) VALUES ($1, $2, $3)`,
		req.UserID, req.URL.OriginalURL, req.URL.ShortURL)
	if err != nil {
		return err
	}
	return nil
}

func (r *urlRepository) GetListURL(ctx context.Context, req *GetListURLRequest) ([]*URLCombined, error) {
	query := `SELECT original_url, short_url FROM urls WHERE user_id = $1`
	rows, err := r.db.QueryContext(ctx, query, req.UserID)
	if err != nil {
		return nil, err
	}
	defer func() {
		if err = rows.Close(); err != nil {
			r.logger.Errorf("rows.Close(): %v", err)
		}
	}()

	var urls []*URLCombined
	for rows.Next() {
		var url URLCombined
		if err = rows.Scan(&url.OriginalURL, &url.ShortURL); err != nil {
			return nil, err
		}
		urls = append(urls, &url)
	}

	if err = rows.Err(); err != nil {
		r.logger.Errorf("rows.Err(): %v", err)
		return nil, err
	}

	return urls, nil
}

func (r *urlRepository) DeleteAllURL(ctx context.Context, req *DeleteAllURLRequest) error {
	query := `DELETE FROM urls WHERE user_id = $1`
	_, err := r.db.ExecContext(ctx, query, req.UserID)
	if err != nil {
		r.logger.Errorf("Failed to delete URLs for user %s: %v", req.UserID, err)
		return err
	}
	return nil
}

func (r *urlRepository) DeleteURL(ctx context.Context, req *DeleteURLRequest) error {
	query := `DELETE FROM urls WHERE (user_id = $1 AND short_url = $2 AND original_url = $3)`
	_, err := r.db.ExecContext(ctx, query, req.UserID, req.ShortURL, req.OriginalURL)
	r.logger.Errorf("DELETE WORKKKK", req.UserID, req.ShortURL, req.OriginalURL)
	if err != nil {
		r.logger.Errorf("Failed to delete URLs for user %s: %v", req.UserID, err)
		return err
	}
	return nil
}
