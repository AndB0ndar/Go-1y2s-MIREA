package auth

import (
	"app/shared/httpx"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type Client struct {
	httpClient *httpx.Client
}

func NewClient(baseURL string, timeout time.Duration) *Client {
	return &Client{
		httpClient: httpx.NewClient(baseURL, timeout),
	}
}

type verifyResponse struct {
	Valid   bool   `json:"valid"`
	Subject string `json:"subject"`
	Error   string `json:"error"`
}

func (c *Client) VerifyToken(ctx context.Context, token, requestID string) error {
	headers := http.Header{}
	headers.Set("Authorization", "Bearer "+token)
	headers.Set("X-Request-ID", requestID)

	resp, err := c.httpClient.Get(ctx, "/auth/verify", headers)
	if err != nil {
		return fmt.Errorf("auth service call failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		var vResp verifyResponse
		if err := json.NewDecoder(resp.Body).Decode(&vResp); err == nil && vResp.Valid {
			return nil
		}
		return fmt.Errorf("auth returned invalid response")
	}
	if resp.StatusCode == http.StatusUnauthorized {
		return fmt.Errorf("unauthorized")
	}
	return fmt.Errorf("auth service error: status %d", resp.StatusCode)
}
