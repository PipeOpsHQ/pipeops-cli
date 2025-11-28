package client

import (
	"context"
	"fmt"
	"math"
	"net/http"
	"time"
)

// HTTPClient wraps http.Client with additional features
type HTTPClient struct {
	client     *http.Client
	maxRetries int
	retryDelay time.Duration
	timeout    time.Duration
}

// NewHTTPClient creates a new HTTP client with sensible defaults
func NewHTTPClient() *HTTPClient {
	return &HTTPClient{
		client: &http.Client{
			Timeout: 30 * time.Second,
			Transport: &http.Transport{
				MaxIdleConns:        10,
				MaxIdleConnsPerHost: 10,
				IdleConnTimeout:     90 * time.Second,
			},
		},
		maxRetries: 3,
		retryDelay: 1 * time.Second,
		timeout:    30 * time.Second,
	}
}

// WithTimeout sets a custom timeout for the client
func (c *HTTPClient) WithTimeout(timeout time.Duration) *HTTPClient {
	c.timeout = timeout
	c.client.Timeout = timeout
	return c
}

// WithRetries sets custom retry parameters
func (c *HTTPClient) WithRetries(maxRetries int, retryDelay time.Duration) *HTTPClient {
	c.maxRetries = maxRetries
	c.retryDelay = retryDelay
	return c
}

// Do executes an HTTP request with retry logic
func (c *HTTPClient) Do(req *http.Request) (*http.Response, error) {
	var resp *http.Response
	var err error

	for attempt := 0; attempt <= c.maxRetries; attempt++ {
		// Clone request for retry
		reqClone := req.Clone(req.Context())

		resp, err = c.client.Do(reqClone)

		// Success or non-retryable error
		if err == nil && !shouldRetry(resp.StatusCode) {
			return resp, nil
		}

		// Last attempt, return error
		if attempt == c.maxRetries {
			if err != nil {
				return nil, fmt.Errorf("request failed after %d attempts: %w", attempt+1, err)
			}
			return resp, nil
		}

		// Calculate backoff delay (exponential backoff)
		backoff := time.Duration(math.Pow(2, float64(attempt))) * c.retryDelay
		if backoff > 30*time.Second {
			backoff = 30 * time.Second
		}

		// Close response body before retry
		if resp != nil && resp.Body != nil {
			resp.Body.Close()
		}

		// Wait before retry
		select {
		case <-time.After(backoff):
			// Continue to next attempt
		case <-req.Context().Done():
			return nil, req.Context().Err()
		}
	}

	return resp, err
}

// DoWithContext executes an HTTP request with a custom context
func (c *HTTPClient) DoWithContext(ctx context.Context, req *http.Request) (*http.Response, error) {
	return c.Do(req.WithContext(ctx))
}

// shouldRetry determines if a request should be retried based on status code
func shouldRetry(statusCode int) bool {
	// Retry on server errors and specific client errors
	switch statusCode {
	case http.StatusTooManyRequests, // 429
		http.StatusInternalServerError, // 500
		http.StatusBadGateway,          // 502
		http.StatusServiceUnavailable,  // 503
		http.StatusGatewayTimeout:      // 504
		return true
	}
	return false
}

// Get performs a GET request with retry logic
func (c *HTTPClient) Get(url string) (*http.Response, error) {
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create GET request: %w", err)
	}
	return c.Do(req)
}

// GetWithContext performs a GET request with context
func (c *HTTPClient) GetWithContext(ctx context.Context, url string) (*http.Response, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create GET request: %w", err)
	}
	return c.Do(req)
}
