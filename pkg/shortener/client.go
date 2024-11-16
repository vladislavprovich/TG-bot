package shortener

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/sirupsen/logrus"
	"io"
	"net/http"
)

type BasicClient struct {
	config     *Config
	httpClient *http.Client
	logger     *logrus.Logger
}

type Client interface {
	CreateShortUrl(ctx context.Context, req *CreateShortURLRequest) (*CreateShortURLResponse, error)
	GetStatsUrl(ctx context.Context, req *GetShortURLStatsRequest) (*GetShortURLStatsResponse, error)
}

func NewBasicClient(config *Config, httpClient *http.Client, logger *logrus.Logger) BasicClient {
	return BasicClient{
		config:     config,
		httpClient: httpClient,
		logger:     logger,
	}
}

func (c *BasicClient) CreateShortUrl(ctx context.Context, req *CreateShortURLRequest) (*CreateShortURLResponse, error) {
	jsonReq, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("error marshalling request: %w", err)
	}

	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, c.config.BaseURL, bytes.NewBuffer(jsonReq))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	httpReq.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("http request failed: %w", err)
	}
	defer func() {
		if err = resp.Body.Close(); err != nil {
			panic(err)
		}
	}()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("bad status code: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading body: %v", err)
	}

	var res CreateShortURLResponse
	if err = json.Unmarshal(body, &res); err != nil {
		return nil, fmt.Errorf("bad response body: %s", string(body))
	}

	return &res, nil
}

func (c *BasicClient) GetStatsUrl(ctx context.Context, req *GetShortURLStatsRequest) (*GetShortURLStatsResponse, error) {
	jsonReq, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("error marshalling request: %w", err)
	}

	httpReq, err := http.NewRequestWithContext(ctx, http.MethodGet, c.config.BaseGetURL, bytes.NewBuffer(jsonReq))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	httpReq.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("http request failed: %w", err)
	}
	defer func() {
		if err = resp.Body.Close(); err != nil {
			panic(err)
		}
	}()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("bad status code: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading body: %v", err)
	}

	var res GetShortURLStatsResponse
	if err = json.Unmarshal(body, &res); err != nil {
		return nil, fmt.Errorf("bad response body: %s", string(body))
	}

	return &res, nil
}
