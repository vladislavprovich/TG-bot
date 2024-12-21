package shortener

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/sirupsen/logrus"
	"io"
	"net/http"
	"net/url"
)

const (
	shorten = "shorten"
	stats   = "stats"
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
		return nil, fmt.Errorf("Service error")
		c.logger.Errorf("error marshalling request: %w", err)
	}

	urlPost := &url.URL{
		Scheme: "http",
		Host:   fmt.Sprintf("%s:%s", c.config.BaseURL, c.config.Port),
		Path:   shorten,
	}

	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, urlPost.String(), bytes.NewBuffer(jsonReq))
	if err != nil {
		return nil, fmt.Errorf("Service error")
		c.logger.Errorf("failed to create request: %w", err)
	}

	httpReq.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("Service error")
		c.logger.Errorf("http request failed: %w", err)

	}
	defer func() {
		if err = resp.Body.Close(); err != nil {
			panic(err)
		}
	}()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Service error")
		c.logger.Errorf("bad status code: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("Service error")
		c.logger.Errorf("error reading body: %v", err)

	}

	var res CreateShortURLResponse
	if err = json.Unmarshal(body, &res); err != nil {
		return nil, fmt.Errorf("Service error")
		c.logger.Errorf("bad response body:", string(body))
	}

	return &res, nil
}

func (c *BasicClient) GetStatsUrl(ctx context.Context, req *GetShortURLStatsRequest) (*GetShortURLStatsResponse, error) {
	urlGet := &url.URL{
		Scheme: "http",
		Host:   fmt.Sprintf("%s:%s", c.config.BaseURL, c.config.Port),
		Path:   fmt.Sprintf("%s/%s", req.ShortURL, stats),
	}

	httpReq, err := http.NewRequestWithContext(ctx, http.MethodGet, urlGet.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("Service error")
		c.logger.Errorf("failed to create request: %w", err)
	}

	httpReq.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("Service error")
		c.logger.Errorf("http request failed: %v", err)
	}
	defer func() {
		if err = resp.Body.Close(); err != nil {
			panic(err)
		}
	}()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Service error")
		c.logger.Errorf("bad status code: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("Service error")
		c.logger.Errorf("error reading body: %v", err)
	}

	var res GetShortURLStatsResponse
	if err = json.Unmarshal(body, &res); err != nil {
		return nil, fmt.Errorf("Service error")
		c.logger.Errorf("bad response body: %s", string(body))
	}

	return &res, nil
}
