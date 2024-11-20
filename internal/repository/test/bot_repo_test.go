package test

import (
	"context"
	"database/sql"
	"errors"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	repository "github.com/vladislavprovich/TG-bot/internal/repository"
	mocks "github.com/vladislavprovich/TG-bot/mocks"
	"go.uber.org/mock/gomock"
	"testing"
)

func TestSaveURLTable(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockDB := mocks.NewMockDB(ctrl)
	logger := logrus.New()
	repo := repository.NewBotRepository(mockDB, logger)

	tests := []struct {
		name       string
		input      *repository.SaveUrlRequest
		mockExpect func()
		wantErr    bool
	}{
		{
			name: "Save URL Success",
			input: &repository.SaveUrlRequest{
				UserID: "user123",
				URL: &repository.URLCombined{
					OriginalURL: "http://example.com",
					ShortURL:    "http://short.url",
				},
			},
			mockExpect: func() {
				mockDB.EXPECT().
					ExecContext(gomock.Any(), "INSERT INTO urls (user_id, original_url, short_url) VALUES (?, ?, ?)",
						"user123", "http://example.com", "http://short.url").
					Return(nil, nil)
			},
			wantErr: false,
		},
		{
			name: "Save URL Failure",
			input: &repository.SaveUrlRequest{
				UserID: "user123",
				URL: &repository.URLCombined{
					OriginalURL: "http://example.com",
					ShortURL:    "http://short.url",
				},
			},
			mockExpect: func() {
				mockDB.EXPECT().
					ExecContext(gomock.Any(), "INSERT INTO urls (user_id, original_url, short_url) VALUES (?, ?, ?)",
						"user123", "http://example.com", "http://short.url").
					Return(nil, errors.New("db error"))
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockExpect()

			err := repo.SaveURL(context.Background(), tt.input)
			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}

}

func TestGetListURLTable(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockDB := mocks.NewMockDB(ctrl)
	logger := logrus.New()
	repo := repository.NewBotRepository(mockDB, logger)

	tests := []struct {
		name       string
		input      *repository.GetListURLRequest
		mockExpect func()
		expected   []*repository.URLCombined
		wantErr    bool
	}{
		{
			name:  "Get List Success",
			input: &repository.GetListURLRequest{UserID: "user123"},
			mockExpect: func() {
				rows := sqlmock.NewRows([]string{"original_url", "short_url"}).
					AddRow("http://example.com", "http://short.url").
					AddRow("http://example2.com", "http://short2.url")

				mockDB.EXPECT().
					QueryContext(gomock.Any(), "SELECT original_url, short_url FROM urls WHERE user_id = ?", "user123").
					Return(rows, nil)
			},
			expected: []*repository.URLCombined{
				{OriginalURL: "http://example.com", ShortURL: "http://short.url"},
				{OriginalURL: "http://example2.com", ShortURL: "http://short2.url"},
			},
			wantErr: false,
		},
		{
			name:  "Get List Failure",
			input: &repository.GetListURLRequest{UserID: "user123"},
			mockExpect: func() {
				mockDB.EXPECT().
					QueryContext(gomock.Any(), "SELECT original_url, short_url FROM urls WHERE user_id = ?", "user123").
					Return(nil, errors.New("db error"))
			},
			expected: nil,
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockExpect()

			result, err := repo.GetListURL(context.Background(), tt.input)
			if tt.wantErr {
				require.Error(t, err)
				require.Nil(t, result)
			} else {
				require.NoError(t, err)
				require.Equal(t, tt.expected, result)
			}
		})
	}
}

func TestDeleteAllURLTable(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockDB := mocks.NewMockDB(ctrl)
	logger := logrus.New()
	repo := repository.NewBotRepository(mockDB, logger)

	tests := []struct {
		name       string
		input      *repository.DeleteAllURLRequest
		mockExpect func()
		wantErr    bool
	}{
		{
			name:  "Delete All Success",
			input: &repository.DeleteAllURLRequest{UserID: "user123"},
			mockExpect: func() {
				mockDB.EXPECT().
					ExecContext(gomock.Any(), "DELETE FROM urls WHERE user_id = ?", "user123").
					Return(nil, nil)
			},
			wantErr: false,
		},
		{
			name:  "Delete All Failure",
			input: &repository.DeleteAllURLRequest{UserID: "user123"},
			mockExpect: func() {
				mockDB.EXPECT().
					ExecContext(gomock.Any(), "DELETE FROM urls WHERE user_id = ?", "user123").
					Return(nil, errors.New("db error"))
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockExpect()

			err := repo.DeleteAllURL(context.Background(), tt.input)
			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestDeleteURLTable(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockDB := mocks.NewMockDB(ctrl)
	logger := logrus.New()
	repo := repository.NewBotRepository(mockDB, logger)

	tests := []struct {
		name       string
		input      *repository.DeleteURLRequest
		mockExpect func()
		wantErr    bool
	}{
		{
			name: "Delete URL Success",
			input: &repository.DeleteURLRequest{
				UserID:      "user123",
				OriginalURL: "http://example.com",
			},
			mockExpect: func() {
				mockDB.EXPECT().
					ExecContext(gomock.Any(), "DELETE FROM urls WHERE user_id = ? AND original_url = ?", "user123", "http://example.com").
					Return(nil, nil)
			},
			wantErr: false,
		},
		{
			name: "Delete URL Failure",
			input: &repository.DeleteURLRequest{
				UserID:      "user123",
				OriginalURL: "http://example.com",
			},
			mockExpect: func() {
				mockDB.EXPECT().
					ExecContext(gomock.Any(), "DELETE FROM urls WHERE user_id = ? AND original_url = ?", "user123", "http://example.com").
					Return(nil, errors.New("db error"))
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockExpect()

			err := repo.DeleteURL(context.Background(), tt.input)
			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

// Example test
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
		EXPECT().ExecContext(ctx, query, params).
		Return(sql.Result(nil), nil)

	err := repo.SaveURL(ctx, req)

	assert.NoError(t, err)
}
