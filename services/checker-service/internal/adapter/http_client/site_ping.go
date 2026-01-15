package http_client

import (
	"context"
	"net/http"
	"time"
)

type HTTPChecker struct {
	client *http.Client
}

func NewHTTPChecker() *HTTPChecker {
	return &HTTPChecker{
		client: &http.Client{
			Timeout: 15 * time.Second,
		},
	}
}

func (h *HTTPChecker) Check(ctx context.Context, url string) (bool, int, time.Duration, error) {
	start := time.Now()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return false, 0, 0, err
	}

	resp, err := h.client.Do(req)
	latency := time.Since(start)

	if err != nil {
		return false, 0, latency, err
	}
	defer resp.Body.Close()

	isUp := resp.StatusCode >= 200 && resp.StatusCode < 300
	return isUp, resp.StatusCode, latency, nil
}
