package middleware

import (
	"bytes"
	"fmt"
	"net/http"
	"strings"
	"time"

	"go-admin/internal/cache"
	"go-admin/internal/logger"

	"github.com/gin-gonic/gin"
)

// CacheConfig defines the configuration for the response cache middleware
type CacheConfig struct {
	// CacheDuration specifies how long responses should be cached
	CacheDuration time.Duration

	// CacheableStatusCodes defines which HTTP status codes should be cached
	CacheableStatusCodes []int

	// SkipCacheFunc allows custom logic to skip caching
	SkipCacheFunc func(*gin.Context) bool

	// KeyGenerator allows custom cache key generation
	KeyGenerator func(*gin.Context) string

	// CacheStore specifies the cache store to use
	CacheStore cache.CacheInterface
}

// DefaultCacheConfig returns a default configuration for the response cache middleware
func DefaultCacheConfig() CacheConfig {
	return CacheConfig{
		CacheDuration: 5 * time.Minute,
		CacheableStatusCodes: []int{
			http.StatusOK,
			http.StatusNotModified,
			http.StatusNoContent,
		},
		SkipCacheFunc: func(c *gin.Context) bool {
			// Skip caching for non-GET requests
			if c.Request.Method != "GET" {
				return true
			}

			// Skip caching if authorization header is present
			if _, exists := c.Request.Header["Authorization"]; exists {
				return true
			}

			// Skip caching if request contains cookies
			if len(c.Request.Cookies()) > 0 {
				return true
			}

			return false
		},
		KeyGenerator: func(c *gin.Context) string {
			// Create a cache key based on the full request URL
			return fmt.Sprintf("response_cache:%s", c.Request.URL.String())
		},
		CacheStore: cache.GetInstance(),
	}
}

// ResponseCacheMiddleware creates a new response cache middleware
func ResponseCacheMiddleware(config CacheConfig) gin.HandlerFunc {
	// Use default config if not provided
	if config.CacheStore == nil {
		config = DefaultCacheConfig()
	}

	return func(c *gin.Context) {
		// Check if we should skip caching
		if config.SkipCacheFunc != nil && config.SkipCacheFunc(c) {
			c.Next()
			return
		}

		// Generate cache key
		cacheKey := config.KeyGenerator(c)

		// Try to get response from cache
		if cachedResponse, found := config.CacheStore.Get(cacheKey); found {
			// Type assert to our cached response structure
			if response, ok := cachedResponse.(*CachedResponse); ok {
				// Set headers from cached response
				for key, values := range response.Headers {
					for _, value := range values {
						c.Header(key, value)
					}
				}

				// Set cache status header
				c.Header("X-Cache-Status", "HIT")

				// Write cached response
				c.Data(response.StatusCode, response.ContentType, response.Body)
				c.Abort()

				// Log cache hit
				logger.DefaultStructuredLogger().
					WithField("cache_key", cacheKey).
					WithField("status", "HIT").
					Info("Response cache hit")

				return
			}
		}

		// Cache miss - proceed with request
		c.Header("X-Cache-Status", "MISS")

		// Capture response
		writer := &cacheResponseBodyWriter{
			ResponseWriter: c.Writer,
			buffer:         &bytes.Buffer{},
		}
		c.Writer = writer

		// Process request
		c.Next()

		// Cache response if it's cacheable
		if shouldCacheResponse(c.Writer.Status(), config.CacheableStatusCodes) {
			// Get response body
			body := writer.buffer.Bytes()

			// Create cached response
			cachedResponse := &CachedResponse{
				StatusCode:  c.Writer.Status(),
				Headers:     make(map[string][]string),
				ContentType: c.Writer.Header().Get("Content-Type"),
				Body:        body,
				CachedAt:    time.Now(),
			}

			// Copy headers
			for key, values := range c.Writer.Header() {
				// Skip certain headers that shouldn't be cached
				if !shouldSkipHeader(key) {
					cachedResponse.Headers[key] = values
				}
			}

			// Store in cache
			if err := config.CacheStore.Set(cacheKey, cachedResponse, config.CacheDuration); err != nil {
				// Log error but don't fail the request
				logger.DefaultStructuredLogger().
					WithField("cache_key", cacheKey).
					WithField("error", err).
					Error("Failed to cache response")
			} else {
				// Log successful caching
				logger.DefaultStructuredLogger().
					WithField("cache_key", cacheKey).
					WithField("status", "MISS").
					WithField("cache_duration", config.CacheDuration).
					Info("Response cached")
			}
		}
	}
}

// CachedResponse represents a cached HTTP response
type CachedResponse struct {
	StatusCode  int
	Headers     map[string][]string
	ContentType string
	Body        []byte
	CachedAt    time.Time
}

// cacheResponseBodyWriter is a wrapper around gin.ResponseWriter that captures the response body
type cacheResponseBodyWriter struct {
	gin.ResponseWriter
	buffer *bytes.Buffer
}

// Write captures the response body
func (w *cacheResponseBodyWriter) Write(b []byte) (int, error) {
	// Write to buffer
	n, err := w.buffer.Write(b)
	if err != nil {
		return n, err
	}

	// Write to original writer
	return w.ResponseWriter.Write(b)
}

// shouldCacheResponse checks if the response status code is cacheable
func shouldCacheResponse(statusCode int, cacheableStatusCodes []int) bool {
	for _, code := range cacheableStatusCodes {
		if statusCode == code {
			return true
		}
	}
	return false
}

// shouldSkipHeader checks if a header should be skipped when caching
func shouldSkipHeader(header string) bool {
	skipHeaders := []string{
		"Set-Cookie",
		"X-Cache-Status",
		"Date",
		"Connection",
		"Transfer-Encoding",
		"Vary",
	}

	for _, h := range skipHeaders {
		if strings.EqualFold(header, h) {
			return true
		}
	}

	return false
}

// InvalidateCacheFunc returns a function that can be used to invalidate cache entries
func InvalidateCacheFunc(cacheStore cache.CacheInterface, keyPattern string) func(c *gin.Context) {
	return func(c *gin.Context) {
		// This is a simplified implementation
		// In a real-world scenario, you might want to implement pattern-based cache invalidation
		// or maintain a list of cache keys for each resource type

		// For now, we'll just clear the entire cache
		if err := cacheStore.Clear(); err != nil {
			logger.DefaultStructuredLogger().
				WithField("error", err).
				Error("Failed to clear cache")
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to clear cache"})
		} else {
			logger.DefaultStructuredLogger().
				Info("Cache cleared successfully")
			c.JSON(http.StatusOK, gin.H{"message": "Cache cleared successfully"})
		}
	}
}