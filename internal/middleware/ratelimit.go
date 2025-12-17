package middleware

import (
	"fmt"
	"net/http"
	"time"

	"go-admin/internal/cache"
	"go-admin/internal/logger"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// RateLimitConfig defines the configuration for rate limiting
type RateLimitConfig struct {
	// Number of requests allowed within the window
	Requests int
	// Time window for the rate limit (e.g., 1*time.Minute)
	Window time.Duration
	// Key generator function for identifying clients
	KeyGenerator func(*gin.Context) string
	// Custom response function when limit is exceeded
	OnLimitReached func(*gin.Context, time.Duration)
	// Whether to use distributed cache (Redis) for rate limiting
	UseDistributedCache bool
}

// DefaultRateLimitConfig returns a default configuration
func DefaultRateLimitConfig() RateLimitConfig {
	return RateLimitConfig{
		Requests:            100,
		Window:              time.Minute,
		KeyGenerator:        DefaultKeyGenerator,
		OnLimitReached:      DefaultOnLimitReached,
		UseDistributedCache: true,
	}
}

// DefaultKeyGenerator generates a key based on client IP
func DefaultKeyGenerator(c *gin.Context) string {
	return fmt.Sprintf("rate_limit:%s", c.ClientIP())
}

// DefaultOnLimitReached returns the default response when rate limit is exceeded
func DefaultOnLimitReached(c *gin.Context, retryAfter time.Duration) {
	c.JSON(http.StatusTooManyRequests, gin.H{
		"error":       "Rate limit exceeded",
		"retry_after": retryAfter.Seconds(),
	})
	c.Abort()
}

// RateLimiter represents a rate limiter
type RateLimiter struct {
	config RateLimitConfig
}

// NewRateLimiter creates a new rate limiter with the given configuration
func NewRateLimiter(config RateLimitConfig) *RateLimiter {
	return &RateLimiter{
		config: config,
	}
}

// NewDefaultRateLimiter creates a rate limiter with default configuration
func NewDefaultRateLimiter() *RateLimiter {
	return NewRateLimiter(DefaultRateLimitConfig())
}

// Allow checks if a client is allowed to make a request
func (rl *RateLimiter) Allow(c *gin.Context) bool {
	key := rl.config.KeyGenerator(c)
	
	// Try to use distributed cache first if enabled
	if rl.config.UseDistributedCache {
		return rl.allowWithDistributedCache(key)
	}
	
	// Fall back to in-memory rate limiting
	return rl.allowWithMemory(key)
}

// allowWithDistributedCache checks rate limit using distributed cache
func (rl *RateLimiter) allowWithDistributedCache(key string) bool {
	cacheInstance := cache.GetInstance()
	
	// Get current request count
	currentCount, exists := cacheInstance.Get(key)
	if !exists {
		// First request in the window
		if err := cacheInstance.Set(key, 1, rl.config.Window); err != nil {
			logger.Error("Failed to set rate limit counter", zap.Error(err), zap.String("key", key))
			// If cache fails, allow the request to avoid blocking legitimate traffic
			return true
		}
		return true
	}
	
	// Convert count to int
	count, ok := currentCount.(int)
	if !ok {
		// Invalid type in cache, reset it
		if err := cacheInstance.Set(key, 1, rl.config.Window); err != nil {
			logger.Error("Failed to reset rate limit counter", zap.Error(err), zap.String("key", key))
			return true
		}
		return true
	}
	
	// Check if limit is exceeded
	if count >= rl.config.Requests {
		return false
	}
	
	// Increment counter
	if err := cacheInstance.Set(key, count+1, rl.config.Window); err != nil {
		logger.Error("Failed to increment rate limit counter", zap.Error(err), zap.String("key", key))
		// If cache fails, allow the request to avoid blocking legitimate traffic
		return true
	}
	
	return true
}

// allowWithMemory checks rate limit using in-memory storage
func (rl *RateLimiter) allowWithMemory(key string) bool {
	// This is a simplified in-memory implementation
	// In a production environment, you might want a more sophisticated approach
	// with proper cleanup and memory management
	
	// For now, we'll use the cache interface which defaults to memory cache
	return rl.allowWithDistributedCache(key)
}

// Limit is the middleware function for rate limiting
func (rl *RateLimiter) Limit() gin.HandlerFunc {
	return func(c *gin.Context) {
		if !rl.Allow(c) {
			rl.config.OnLimitReached(c, rl.config.Window)
			return
		}
		
		c.Next()
	}
}

// UserBasedKeyGenerator generates a key based on user ID if available, otherwise falls back to IP
func UserBasedKeyGenerator(c *gin.Context) string {
	// Try to get user ID from context
	if userID, exists := c.Get("userID"); exists {
		return fmt.Sprintf("rate_limit:user:%v", userID)
	}
	
	// Fall back to IP-based key
	return DefaultKeyGenerator(c)
}

// APIBasedKeyGenerator generates a key based on API endpoint and client IP
func APIBasedKeyGenerator(c *gin.Context) string {
	return fmt.Sprintf("rate_limit:api:%s:%s", c.Request.Method, c.ClientIP())
}

// UserAPIBasedKeyGenerator generates a key based on user ID and API endpoint
func UserAPIBasedKeyGenerator(c *gin.Context) string {
	if userID, exists := c.Get("userID"); exists {
		return fmt.Sprintf("rate_limit:user:%v:api:%s", userID, c.Request.Method)
	}
	
	// Fall back to API-based key
	return APIBasedKeyGenerator(c)
}

// RateLimit creates a rate limiting middleware with default configuration
func RateLimit() gin.HandlerFunc {
	limiter := NewDefaultRateLimiter()
	return limiter.Limit()
}

// RateLimitWithConfig creates a rate limiting middleware with custom configuration
func RateLimitWithConfig(config RateLimitConfig) gin.HandlerFunc {
	limiter := NewRateLimiter(config)
	return limiter.Limit()
}

// RateLimitForAPI creates a rate limiting middleware specifically for API endpoints
func RateLimitForAPI(requests int, window time.Duration) gin.HandlerFunc {
	config := DefaultRateLimitConfig()
	config.Requests = requests
	config.Window = window
	config.KeyGenerator = APIBasedKeyGenerator
	
	limiter := NewRateLimiter(config)
	return limiter.Limit()
}

// RateLimitForUser creates a rate limiting middleware specifically for authenticated users
func RateLimitForUser(requests int, window time.Duration) gin.HandlerFunc {
	config := DefaultRateLimitConfig()
	config.Requests = requests
	config.Window = window
	config.KeyGenerator = UserBasedKeyGenerator
	
	limiter := NewRateLimiter(config)
	return limiter.Limit()
}

// RateLimitForUserAPI creates a rate limiting middleware for authenticated users on specific APIs
func RateLimitForUserAPI(requests int, window time.Duration) gin.HandlerFunc {
	config := DefaultRateLimitConfig()
	config.Requests = requests
	config.Window = window
	config.KeyGenerator = UserAPIBasedKeyGenerator
	
	limiter := NewRateLimiter(config)
	return limiter.Limit()
}
