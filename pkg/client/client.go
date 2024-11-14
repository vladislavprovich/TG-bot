package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

type Client struct {
	OriginalURL string
	CustomURL   string
	httpClient  *http.Client
}

func NewURLShortenerClient(OriginalURL string, CustomURL string) *Client {
	return &Client{
		OriginalURL: OriginalURL,
		CustomURL:   CustomURL,
		httpClient: &http.Client{
			Timeout: time.Second * 10,
		},
	}
}

func (c *Client) ShortenerPOST(originalURL string, custom string) (string, error) {

	request := Request{
		URL:         originalURL,
		CustomAlias: &custom,
	}

	jsonReq, err := json.Marshal(request)
	if err != nil {
		return "", fmt.Errorf("error marshalling request: %v", err)
	}

	baseURL := "/shorten"

	req, err := http.NewRequest("POST", baseURL, bytes.NewBuffer(jsonReq))
	if err != nil {
		return "", fmt.Errorf("Failed to create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("http request failed: %v", err)
	}

	defer func() {
		if err = resp.Body.Close(); err != nil {
			panic(err)
		}
	}()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("bad status code: %d", resp.StatusCode)
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("error reading body: %v", err)
	}

	var response Response
	if err = json.Unmarshal(body, &response); err != nil {
		return "", fmt.Errorf("bad response body: %s", string(body))
	}

	if custom != "" {
		response.ShortURL = custom
	}

	return response.ShortURL, nil
}
