package client

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestNewHTTPClient(t *testing.T) {
	client := NewHTTPClient()

	if client == nil {
		t.Fatal("NewHTTPClient() returned nil")
	}

	if client.maxRetries != 3 {
		t.Errorf("maxRetries = %d, want 3", client.maxRetries)
	}

	if client.timeout != 30*time.Second {
		t.Errorf("timeout = %v, want 30s", client.timeout)
	}
}

func TestWithTimeout(t *testing.T) {
	client := NewHTTPClient().WithTimeout(10 * time.Second)

	if client.timeout != 10*time.Second {
		t.Errorf("timeout = %v, want 10s", client.timeout)
	}

	if client.client.Timeout != 10*time.Second {
		t.Errorf("client.Timeout = %v, want 10s", client.client.Timeout)
	}
}

func TestWithRetries(t *testing.T) {
	client := NewHTTPClient().WithRetries(5, 2*time.Second)

	if client.maxRetries != 5 {
		t.Errorf("maxRetries = %d, want 5", client.maxRetries)
	}

	if client.retryDelay != 2*time.Second {
		t.Errorf("retryDelay = %v, want 2s", client.retryDelay)
	}
}

func TestShouldRetry(t *testing.T) {
	tests := []struct {
		name       string
		statusCode int
		want       bool
	}{
		{"429 Too Many Requests", http.StatusTooManyRequests, true},
		{"500 Internal Server Error", http.StatusInternalServerError, true},
		{"502 Bad Gateway", http.StatusBadGateway, true},
		{"503 Service Unavailable", http.StatusServiceUnavailable, true},
		{"504 Gateway Timeout", http.StatusGatewayTimeout, true},
		{"200 OK", http.StatusOK, false},
		{"400 Bad Request", http.StatusBadRequest, false},
		{"401 Unauthorized", http.StatusUnauthorized, false},
		{"404 Not Found", http.StatusNotFound, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := shouldRetry(tt.statusCode)
			if got != tt.want {
				t.Errorf("shouldRetry(%d) = %v, want %v", tt.statusCode, got, tt.want)
			}
		})
	}
}

func TestDoSuccessfulRequest(t *testing.T) {
	// Create test server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("success"))
	}))
	defer server.Close()

	client := NewHTTPClient()
	req, err := http.NewRequest(http.MethodGet, server.URL, nil)
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}

	resp, err := client.Do(req)
	if err != nil {
		t.Fatalf("Do() error = %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("StatusCode = %d, want %d", resp.StatusCode, http.StatusOK)
	}
}

func TestDoWithRetry(t *testing.T) {
	attempts := 0
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		attempts++
		if attempts < 3 {
			w.WriteHeader(http.StatusServiceUnavailable)
			return
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	client := NewHTTPClient().WithRetries(3, 10*time.Millisecond)
	req, err := http.NewRequest(http.MethodGet, server.URL, nil)
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}

	resp, err := client.Do(req)
	if err != nil {
		t.Fatalf("Do() error = %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("StatusCode = %d, want %d", resp.StatusCode, http.StatusOK)
	}

	if attempts != 3 {
		t.Errorf("attempts = %d, want 3", attempts)
	}
}

func TestDoWithContext(t *testing.T) {
	// Create test server with delay
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(100 * time.Millisecond)
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	client := NewHTTPClient()
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Millisecond)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, server.URL, nil)
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}

	_, err = client.DoWithContext(ctx, req)
	if err == nil {
		t.Error("Expected context deadline exceeded error, got nil")
	}
}

func TestGet(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("Method = %s, want GET", r.Method)
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	client := NewHTTPClient()
	resp, err := client.Get(server.URL)
	if err != nil {
		t.Fatalf("Get() error = %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("StatusCode = %d, want %d", resp.StatusCode, http.StatusOK)
	}
}
