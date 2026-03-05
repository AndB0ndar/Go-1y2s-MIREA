package httpx

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"time"
)

type Client struct {
	BaseURL    string
	HTTPClient *http.Client
}

func NewClient(baseURL string, timeout time.Duration) *Client {
	return &Client{
		BaseURL: baseURL,
		HTTPClient: &http.Client{
			Timeout: timeout,
		},
	}
}

func (c *Client) Get(ctx context.Context, path string, headers http.Header) (*http.Response, error) {
	return c.do(ctx, http.MethodGet, path, nil, headers)
}

func (c *Client) Post(ctx context.Context, path string, body interface{}, headers http.Header) (*http.Response, error) {
	return c.do(ctx, http.MethodPost, path, body, headers)
}

func (c *Client) do(ctx context.Context, method, path string, body interface{}, headers http.Header) (*http.Response, error) {
	url := c.BaseURL + path
	var reqBody io.Reader
	if body != nil {
		b, err := json.Marshal(body)
		if err != nil {
			return nil, err
		}
		reqBody = bytes.NewReader(b)
	}
	req, err := http.NewRequestWithContext(ctx, method, url, reqBody)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	for k, v := range headers {
		req.Header[k] = v
	}
	return c.HTTPClient.Do(req)
}
