package test

import (
	"context"
	"database/sql"
	"errors"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	repository "github.com/vladislavprovich/TG-bot/internal/repository"
	mocks "github.com/vladislavprovich/TG-bot/mocks"
	"go.uber.org/mock/gomock"
	"testing"
)

// Example
func TestUrlRepository_SaveURL(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Create mock instances
	mockDB := mocks.NewMockDB(ctrl)
	logger := logrus.New()

	// Create repository with mocks
	repo := repository.NewBotRepository(mockDB, logger)

	ctx := context.Background()
	req := &repository.SaveUrlRequest{
		UserID: "user123",
		URL: &repository.URLCombined{
			OriginalURL: "http://example.com",
			ShortURL:    "http://short.url/abc",
		},
	}

	query := "INSERT INTO urls (user_id, original_url, short_url) VALUES (:user_id, :original_url, :short_url)"
	params := map[string]interface{}{
		"user_id":      req.UserID,
		"original_url": req.URL.OriginalURL,
		"short_url":    req.URL.ShortURL,
	}

	// Set expectations
	mockDB.
		EXPECT().
		NamedExecContext(ctx, query, params).
		Return(sql.Result(nil), nil)

	err := repo.SaveURL(ctx, req)

	assert.NoError(t, err)
}

func TestSaveURL(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockDB := mocks.NewMockExtContext(ctrl)

	mockDB.EXPECT().ExecContext(gomock.Any(), "INSERT INTO urls (user_id, original_url, short_url) VALUES (?, ?, ?)",
		"user123", "http://example.com", "http://short.url").Return(nil, nil)
	logger := logrus.New()
	repo := repository.NewBotRepository(mockDB, *logger)
	err := repo.SaveURL(context.TODO(), &repository.SaveUrlRequest{
		UserID: "user123",
		URL: &repository.URLCombined{
			OriginalURL: "http://example.com",
			ShortURL:    "http://short.url",
		},
	})
	require.NoError(t, err)
}

func TestGetListURL(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockDB := mocks.NewMockExtContext(ctrl)

	mockDB.EXPECT().QueryContext(gomock.Any(), "SELECT original_url, short_url FROM urls WHERE user_id = ?", "user123").
		Return([]*repository.URLCombined{
			{OriginalURL: "http://example.com", ShortURL: "http://short.url"},
		}, nil)

	logger := logrus.New()
	repo := repository.NewBotRepository(mockDB, *logger)

	_, err := repo.GetListURL(context.TODO(),
		&repository.GetListURLRequest{UserID: "user123"},
	)

	require.NoError(t, err)
}

func TestDeleteAllURL(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockDB := mocks.NewMockExtContext(ctrl)

	mockDB.EXPECT().ExecContext(gomock.Any(), "DELETE FROM urls WHERE user_id = ?", "user123").Return(nil, nil)

	logger := logrus.New()

	repo := repository.NewBotRepository(mockDB, *logger)

	err := repo.DeleteAllURL(context.TODO(), &repository.DeleteAllURLRequest{UserID: "user123"})

	require.NoError(t, err)
}

func TestDeleteURL(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockDB := mocks.NewMockExtContext(ctrl)

	mockDB.EXPECT().ExecContext(gomock.Any(), "DELETE FROM urls WHERE user_id = ? AND original_url = ?", "user123", "http://example.com").Return(nil, nil)

	logger := logrus.New()

	repo := repository.NewBotRepository(mockDB, *logger)

	err := repo.DeleteURL(context.TODO(), &repository.DeleteURLRequest{UserID: "user123", OriginalURL: "http://example.com"})

	require.NoError(t, err)
}

func TestDeleteAllUrl_Error(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockDB := mocks.NewMockExtContext(ctrl)

	mockDB.EXPECT().ExecContext(gomock.Any(), "DELETE FROM urls WHERE user_id = ?", "user123").Return(nil, errors.New("database error"))

	logger := logrus.New()

	repo := repository.NewBotRepository(mockDB, *logger)

	err := repo.DeleteAllURL(context.TODO(), &repository.DeleteAllURLRequest{UserID: "user123"})

	require.Error(t, err)
}

func TestDeleteURL_Error(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockDB := mocks.NewMockExtContext(ctrl)

	mockDB.EXPECT().ExecContext(gomock.Any(), "DELETE FROM urls WHERE user_id = ? AND original_url = ?", "user123", "http://example.com").Return(nil, errors.New("database error"))

	logger := logrus.New()

	repo := repository.NewBotRepository(mockDB, *logger)

	err := repo.DeleteURL(context.TODO(), &repository.DeleteURLRequest{UserID: "user123", OriginalURL: "http://example.com"})

	require.Error(t, err)
}
