package repository

import (
	"context"
	"database/sql"
	"github.com/vladislavprovich/TG-bot/internal/models"
)

type UrlRepo interface {
	SaveURL(ctx context.Context, url models.User) error
	GetList(ctx context.Context, userID string) (models.User, error)
	CreateUser(ctx context.Context, user models.User) (string, error)
}

type urlRepo struct {
	db *sql.DB
}

func NewRepo(db *sql.DB) UrlRepo {
	return &urlRepo{db: db}
}

func (r *urlRepo) SaveURL(ctx context.Context, url models.User) error {

}

func (r *urlRepo) GetList(ctx context.Context, username string) (models.User, error) {
	var user models.User
	query := "SELECT user_id, username FROM users WHERE username = $1"
	err := r.db.QueryRowContext(ctx, query, username).Scan(&user.ID, &user.Username)
	if err == sql.ErrNoRows {
		return user, nil
	}
	return user, err
}

func (r *urlRepo) CreateUser(ctx context.Context, user models.User) (string, error) {
	query := "INSERT INTO users (username) VALUES ($1) RETURNING user_id"
	err := r.db.QueryRowContext(ctx, query, user.Username).Scan(&user.ID)
	return user.ID, err
}
