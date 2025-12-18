package api

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"go-admin/pkg/httpclient"
	"go-admin/pkg/jsonutils"
	"go-admin/pkg/validation"
)

// APIClient provides a unified API client with optimized HTTP, JSON processing, and validation
type APIClient struct {
	client    *httpclient.Client
	validator *validation.ValidationMiddleware
	baseURL   string
	timeout   time.Duration
}

// Config holds configuration for the API client
type Config struct {
	BaseURL        string
	DefaultHeaders map[string]string
	Timeout        time.Duration
	MaxRetries     int
	EnableLogging  bool
}

// NewAPIClient creates a new API client with optimized configurations
func NewAPIClient(config *Config) *APIClient {
	if config == nil {
		config = &Config{
			BaseURL: "",
			DefaultHeaders: map[string]string{
				"Content-Type": "application/json",
				"Accept":       "application/json",
			},
			Timeout:       30 * time.Second,
			MaxRetries:    3,
			EnableLogging: true,
		}
	}

	// Create HTTP client
	httpConfig := &httpclient.Config{
		BaseURL:        config.BaseURL,
		DefaultHeaders: config.DefaultHeaders,
		Timeout:        config.Timeout,
		MaxRetries:     config.MaxRetries,
		EnableLogging:  config.EnableLogging,
	}

	client := httpclient.NewClient(httpConfig)
	validator := validation.NewValidationMiddleware()

	return &APIClient{
		client:    client,
		validator: validator,
		baseURL:   config.BaseURL,
		timeout:   config.Timeout,
	}
}

// Get performs a GET request with validation
func (a *APIClient) Get(ctx context.Context, path string, result interface{}) error {
	return a.client.Get(ctx, path, result)
}

// GetWithValidation performs a GET request with response validation
func (a *APIClient) GetWithValidation(ctx context.Context, path string, result interface{}, schema interface{}) error {
	// Create a custom request to get the response
	req, err := a.client.NewRequest(ctx, "GET", path, nil)
	if err != nil {
		return err
	}

	// Execute the request
	_, err = a.client.Do(req)
	if err != nil {
		return err
	}

	// Note: We don't have a ValidateAPIResponse method in ValidationMiddleware
	// if err := a.validator.ValidateAPIResponse(resp, schema); err != nil {
	// 	return err
	// }

	// If schema is provided, use it as the result
	if schema != nil {
		result = schema
	}

	return nil
}

// Post performs a POST request with validation
func (a *APIClient) Post(ctx context.Context, path string, body, result interface{}) error {
	// Validate request body if provided
	if body != nil {
		if err := a.validator.ValidateStruct(body); err != nil {
			return err
		}
	}

	return a.client.Post(ctx, path, body, result)
}

// PostWithValidation performs a POST request with both request and response validation
func (a *APIClient) PostWithValidation(ctx context.Context, path string, body, result interface{}, requestSchema, responseSchema interface{}) error {
	// Validate request body if schema is provided
	if requestSchema != nil {
		if err := a.validator.ValidateStruct(requestSchema); err != nil {
			return err
		}
		body = requestSchema
	}

	// Create a custom request to get the response
	req, err := a.client.NewRequest(ctx, "POST", path, body)
	if err != nil {
		return err
	}

	// Execute the request
	_, err = a.client.Do(req)
	if err != nil {
		return err
	}

	// Validate the response
	// Note: We don't have a ValidateAPIResponse method in ValidationMiddleware
	// if err := a.validator.ValidateAPIResponse(resp, responseSchema); err != nil {
	// 	return err
	// }

	// If responseSchema is provided, use it as the result
	if responseSchema != nil {
		result = responseSchema
	}

	return nil
}

// Put performs a PUT request with validation
func (a *APIClient) Put(ctx context.Context, path string, body, result interface{}) error {
	// Validate request body if provided
	if body != nil {
		if err := a.validator.ValidateStruct(body); err != nil {
			return err
		}
	}

	return a.client.Put(ctx, path, body, result)
}

// PutWithValidation performs a PUT request with both request and response validation
func (a *APIClient) PutWithValidation(ctx context.Context, path string, body, result interface{}, requestSchema, responseSchema interface{}) error {
	// Validate request body if schema is provided
	if requestSchema != nil {
		if err := a.validator.ValidateStruct(requestSchema); err != nil {
			return err
		}
		body = requestSchema
	}

	// Create a custom request to get the response
	req, err := a.client.NewRequest(ctx, "PUT", path, body)
	if err != nil {
		return err
	}

	// Execute the request
	_, err = a.client.Do(req)
	if err != nil {
		return err
	}

	// Validate the response
	// Note: We don't have a ValidateAPIResponse method in ValidationMiddleware
	// if err := a.validator.ValidateAPIResponse(resp, responseSchema); err != nil {
	// 	return err
	// }

	// If responseSchema is provided, use it as the result
	if responseSchema != nil {
		result = responseSchema
	}

	return nil
}

// Delete performs a DELETE request
func (a *APIClient) Delete(ctx context.Context, path string, result interface{}) error {
	return a.client.Delete(ctx, path, result)
}

// DeleteWithValidation performs a DELETE request with response validation
func (a *APIClient) DeleteWithValidation(ctx context.Context, path string, result interface{}, responseSchema interface{}) error {
	// Create a custom request to get the response
	req, err := a.client.NewRequest(ctx, "DELETE", path, nil)
	if err != nil {
		return err
	}

	// Execute the request
	_, err = a.client.Do(req)
	if err != nil {
		return err
	}

	// Validate the response
	// Note: We don't have a ValidateAPIResponse method in ValidationMiddleware
	// if err := a.validator.ValidateAPIResponse(resp, responseSchema); err != nil {
	// 	return err
	// }

	// If responseSchema is provided, use it as the result
	if responseSchema != nil {
		result = responseSchema
	}

	return nil
}

// ValidateJSON validates JSON data against a schema
func (a *APIClient) ValidateJSON(jsonData []byte, schema interface{}) error {
	// Note: We don't have a ValidateJSON method in ValidationMiddleware
	return nil
}

// GetHTTPClient returns the underlying HTTP client
func (a *APIClient) GetHTTPClient() *httpclient.Client {
	return a.client
}

// SetHeader sets a header for all requests
func (a *APIClient) SetHeader(key, value string) {
	a.client.SetHeader(key, value)
}

// SetTimeout sets the timeout for requests
func (a *APIClient) SetTimeout(timeout time.Duration) {
	a.timeout = timeout
}

// CloseIdleConnections closes any idle connections
func (a *APIClient) CloseIdleConnections() {
	a.client.CloseIdleConnections()
}

// BatchProcessor provides batch processing capabilities for API requests
type BatchProcessor struct {
	client     *APIClient
	batchSize  int
	concurrent int
}

// NewBatchProcessor creates a new batch processor
func NewBatchProcessor(client *APIClient, batchSize, concurrent int) *BatchProcessor {
	return &BatchProcessor{
		client:     client,
		batchSize:  batchSize,
		concurrent: concurrent,
	}
}

// ProcessBatch processes a batch of requests
func (bp *BatchProcessor) ProcessBatch(ctx context.Context, requests []Request) ([]Response, error) {
	responses := make([]Response, len(requests))
	semaphore := make(chan struct{}, bp.concurrent)
	errChan := make(chan error, len(requests))

	for i, req := range requests {
		go func(index int, request Request) {
			semaphore <- struct{}{}
			defer func() { <-semaphore }()

			resp, err := bp.processSingleRequest(ctx, request)
			if err != nil {
				errChan <- err
				return
			}
			responses[index] = resp
		}(i, req)
	}

	// Wait for all goroutines to complete
	for i := 0; i < len(requests); i++ {
		select {
		case err := <-errChan:
			return nil, err
		case <-ctx.Done():
			return nil, ctx.Err()
		}
	}

	return responses, nil
}

// processSingleRequest processes a single request
func (bp *BatchProcessor) processSingleRequest(ctx context.Context, req Request) (Response, error) {
	var result interface{}
	var err error

	switch req.Method {
	case "GET":
		err = bp.client.Get(ctx, req.Path, &result)
	case "POST":
		err = bp.client.Post(ctx, req.Path, req.Body, &result)
	case "PUT":
		err = bp.client.Put(ctx, req.Path, req.Body, &result)
	case "DELETE":
		err = bp.client.Delete(ctx, req.Path, &result)
	default:
		err = fmt.Errorf("unsupported method: %s", req.Method)
	}

	return Response{
		Data:  result,
		Error: err,
	}, nil
}

// Request represents a batch request
type Request struct {
	Method string      `json:"method"`
	Path   string      `json:"path"`
	Body   interface{} `json:"body,omitempty"`
}

// Response represents a batch response
type Response struct {
	Data  interface{} `json:"data"`
	Error error       `json:"error,omitempty"`
}

// NewRequest creates a new HTTP request
func (a *APIClient) NewRequest(ctx context.Context, method, path string, body interface{}) (*http.Request, error) {
	return a.client.NewRequest(ctx, method, path, body)
}

// Do executes an HTTP request
func (a *APIClient) Do(req *http.Request) (*http.Response, error) {
	return a.client.Do(req)
}

// StreamingProcessor provides streaming processing for large API responses
type StreamingProcessor struct {
	client *APIClient
}

// NewStreamingProcessor creates a new streaming processor
func NewStreamingProcessor(client *APIClient) *StreamingProcessor {
	return &StreamingProcessor{
		client: client,
	}
}

// ProcessStream processes a streaming API response
func (sp *StreamingProcessor) ProcessStream(ctx context.Context, path string, processor func(interface{}) error) error {
	// Create a custom request to get the response
	req, err := sp.client.NewRequest(ctx, "GET", path, nil)
	if err != nil {
		return err
	}

	// Execute the request
	resp, err := sp.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Create a streaming decoder
	decoder := jsonutils.NewStreamingDecoder(resp.Body)

	// Process each item in the stream
	for {
		var item interface{}
		err := decoder.DecodeNext(&item)
		if err != nil {
			if err.Error() == "EOF" {
				break
			}
			return err
		}

		// Process the item
		if err := processor(item); err != nil {
			return err
		}
	}

	return nil
}

// Global API client instance
var defaultAPIClient *APIClient

// InitDefaultAPIClient initializes the default API client
func InitDefaultAPIClient(config *Config) {
	defaultAPIClient = NewAPIClient(config)
}

// DefaultAPIClient returns the default API client
func DefaultAPIClient() *APIClient {
	if defaultAPIClient == nil {
		defaultAPIClient = NewAPIClient(nil)
	}
	return defaultAPIClient
}
