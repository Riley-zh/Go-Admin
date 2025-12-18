package httpclient

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"go-admin/internal/logger"
	"go-admin/pkg/errors"

	"go.uber.org/zap"
)

// Client represents a reusable HTTP client with optimized configurations
type Client struct {
	httpClient *http.Client
	baseURL    string
	headers    map[string]string
	timeout    time.Duration
	retries    int
}

// Config holds configuration for the HTTP client
type Config struct {
	BaseURL        string
	DefaultHeaders map[string]string
	Timeout        time.Duration
	MaxIdleConns   int
	RequestTimeout time.Duration
	MaxRetries     int
	EnableLogging  bool
	EnableMetrics  bool
}

// DefaultConfig returns a default configuration for the HTTP client
func DefaultConfig() *Config {
	return &Config{
		BaseURL: "",
		DefaultHeaders: map[string]string{
			"Content-Type": "application/json",
			"Accept":       "application/json",
		},
		Timeout:        30 * time.Second,
		MaxIdleConns:   100,
		RequestTimeout: 30 * time.Second,
		MaxRetries:     3,
		EnableLogging:  true,
		EnableMetrics:  true,
	}
}

// NewClient creates a new HTTP client with the given configuration
func NewClient(config *Config) *Client {
	if config == nil {
		config = DefaultConfig()
	}

	// Create HTTP client with optimized transport
	transport := &http.Transport{
		MaxIdleConns:        config.MaxIdleConns,
		MaxIdleConnsPerHost: 10,
		IdleConnTimeout:     90 * time.Second,
		DisableCompression:  false,
	}

	client := &Client{
		httpClient: &http.Client{
			Transport: transport,
			Timeout:   config.Timeout,
		},
		baseURL: config.BaseURL,
		headers: config.DefaultHeaders,
		timeout: config.RequestTimeout,
		retries: config.MaxRetries,
	}

	return client
}

// SetHeader sets a header for all requests
func (c *Client) SetHeader(key, value string) {
	if c.headers == nil {
		c.headers = make(map[string]string)
	}
	c.headers[key] = value
}

// Get performs a GET request
func (c *Client) Get(ctx context.Context, path string, result interface{}) error {
	return c.doRequest(ctx, http.MethodGet, path, nil, result)
}

// Post performs a POST request
func (c *Client) Post(ctx context.Context, path string, body, result interface{}) error {
	return c.doRequest(ctx, http.MethodPost, path, body, result)
}

// Put performs a PUT request
func (c *Client) Put(ctx context.Context, path string, body, result interface{}) error {
	return c.doRequest(ctx, http.MethodPut, path, body, result)
}

// Delete performs a DELETE request
func (c *Client) Delete(ctx context.Context, path string, result interface{}) error {
	return c.doRequest(ctx, http.MethodDelete, path, nil, result)
}

// doRequest performs the actual HTTP request with retry logic
func (c *Client) doRequest(ctx context.Context, method, path string, body, result interface{}) error {
	var lastErr error

	// Build full URL
	url := c.baseURL + path

	// Marshal body if provided
	var reqBody io.Reader
	if body != nil {
		jsonData, err := json.Marshal(body)
		if err != nil {
			return fmt.Errorf("failed to marshal request body: %w", err)
		}
		reqBody = bytes.NewBuffer(jsonData)
	}

	// Retry logic
	for attempt := 0; attempt <= c.retries; attempt++ {
		if attempt > 0 {
			// Exponential backoff
			backoffTime := time.Duration(1<<uint(attempt-1)) * time.Second
			select {
			case <-time.After(backoffTime):
			case <-ctx.Done():
				return ctx.Err()
			}
		}

		// Create request
		req, err := http.NewRequestWithContext(ctx, method, url, reqBody)
		if err != nil {
			lastErr = err
			continue
		}

		// Set headers
		for key, value := range c.headers {
			req.Header.Set(key, value)
		}

		// Execute request
		resp, err := c.httpClient.Do(req)
		if err != nil {
			lastErr = err
			continue
		}

		// Process response
		defer resp.Body.Close()

		// Read response body
		respBody, err := io.ReadAll(resp.Body)
		if err != nil {
			lastErr = err
			continue
		}

		// Check status code
		if resp.StatusCode >= 400 {
			lastErr = c.handleErrorResponse(resp.StatusCode, respBody)
			continue
		}

		// Unmarshal response if result is provided
		if result != nil && len(respBody) > 0 {
			if err := json.Unmarshal(respBody, result); err != nil {
				lastErr = fmt.Errorf("failed to unmarshal response: %w", err)
				continue
			}
		}

		// Success
		return nil
	}

	return fmt.Errorf("request failed after %d attempts: %w", c.retries+1, lastErr)
}

// handleErrorResponse processes error responses and returns appropriate errors
func (c *Client) handleErrorResponse(statusCode int, body []byte) error {
	// Try to parse error response
	var errorResp struct {
		Error   string `json:"error"`
		Details string `json:"details"`
		Message string `json:"message"`
	}

	if err := json.Unmarshal(body, &errorResp); err == nil {
		// Use structured error if available
		errorMsg := errorResp.Error
		if errorMsg == "" {
			errorMsg = errorResp.Message
		}
		if errorMsg == "" {
			errorMsg = "Unknown error"
		}

		details := errorResp.Details
		if details == "" {
			details = string(body)
		}

		return errors.New(statusCode, errorMsg, details)
	}

	// Fallback to generic error with body content
	return errors.New(statusCode, http.StatusText(statusCode), string(body))
}

// GetWithHeaders performs a GET request with custom headers
func (c *Client) GetWithHeaders(ctx context.Context, path string, headers map[string]string, result interface{}) error {
	return c.doRequestWithHeaders(ctx, http.MethodGet, path, nil, result, headers)
}

// PostWithHeaders performs a POST request with custom headers
func (c *Client) PostWithHeaders(ctx context.Context, path string, body, result interface{}, headers map[string]string) error {
	return c.doRequestWithHeaders(ctx, http.MethodPost, path, body, result, headers)
}

// doRequestWithHeaders performs a request with custom headers
func (c *Client) doRequestWithHeaders(ctx context.Context, method, path string, body, result interface{}, customHeaders map[string]string) error {
	// Create a copy of the client with custom headers
	tempClient := &Client{
		httpClient: c.httpClient,
		baseURL:    c.baseURL,
		headers:    make(map[string]string),
		timeout:    c.timeout,
		retries:    c.retries,
	}

	// Copy default headers
	for k, v := range c.headers {
		tempClient.headers[k] = v
	}

	// Add custom headers
	for k, v := range customHeaders {
		tempClient.headers[k] = v
	}

	return tempClient.doRequest(ctx, method, path, body, result)
}

// NewRequest creates a new HTTP request with context and common headers
func (c *Client) NewRequest(ctx context.Context, method, path string, body interface{}) (*http.Request, error) {
	url := c.baseURL + path

	var reqBody io.Reader
	if body != nil {
		jsonData, err := json.Marshal(body)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal request body: %w", err)
		}
		reqBody = bytes.NewBuffer(jsonData)
	}

	req, err := http.NewRequestWithContext(ctx, method, url, reqBody)
	if err != nil {
		return nil, err
	}

	// Set headers
	for key, value := range c.headers {
		req.Header.Set(key, value)
	}

	return req, nil
}

// Do executes an HTTP request and returns the response
func (c *Client) Do(req *http.Request) (*http.Response, error) {
	return c.httpClient.Do(req)
}

// CloseIdleConnections closes any idle connections
func (c *Client) CloseIdleConnections() {
	if transport, ok := c.httpClient.Transport.(*http.Transport); ok {
		transport.CloseIdleConnections()
	}
}

// LogRequest logs HTTP request details
func (c *Client) LogRequest(req *http.Request) {
	logger.Info("HTTP Request",
		zap.String("method", req.Method),
		zap.String("url", req.URL.String()),
		zap.String("user_agent", req.UserAgent()),
	)
}

// LogResponse logs HTTP response details
func (c *Client) LogResponse(resp *http.Response) {
	logger.Info("HTTP Response",
		zap.Int("status_code", resp.StatusCode),
		zap.String("status", resp.Status),
	)
}

// ValidateURL checks if a URL is valid
func ValidateURL(url string) error {
	if url == "" {
		return errors.New(400, "URL cannot be empty", "")
	}

	if !strings.HasPrefix(url, "http://") && !strings.HasPrefix(url, "https://") {
		return errors.New(400, "URL must start with http:// or https://", "")
	}

	return nil
}
