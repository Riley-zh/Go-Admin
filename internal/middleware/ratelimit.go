package middleware

import (
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

// RateLimiter represents a rate limiter
type RateLimiter struct {
	limit    int
	duration time.Duration
	clients  map[string]*client
	mutex    sync.Mutex
}

// client represents a client's rate limit info
type client struct {
	requests []time.Time
}

// NewRateLimiter creates a new rate limiter
func NewRateLimiter(limit int, duration time.Duration) *RateLimiter {
	return &RateLimiter{
		limit:    limit,
		duration: duration,
		clients:  make(map[string]*client),
	}
}

// Allow checks if a client is allowed to make a request
func (rl *RateLimiter) Allow(clientIP string) bool {
	rl.mutex.Lock()
	defer rl.mutex.Unlock()

	now := time.Now()

	// Get or create client info
	c, exists := rl.clients[clientIP]
	if !exists {
		c = &client{requests: []time.Time{}}
		rl.clients[clientIP] = c
	}

	// Remove old requests
	validRequests := []time.Time{}
	for _, reqTime := range c.requests {
		if now.Sub(reqTime) < rl.duration {
			validRequests = append(validRequests, reqTime)
		}
	}
	c.requests = validRequests

	// Check if limit is exceeded
	if len(c.requests) >= rl.limit {
		return false
	}

	// Add current request
	c.requests = append(c.requests, now)
	return true
}

// CleanUp removes old client data periodically
func (rl *RateLimiter) CleanUp() {
	ticker := time.NewTicker(1 * time.Minute)
	defer ticker.Stop()

	for range ticker.C {
		rl.mutex.Lock()
		now := time.Now()

		for clientIP, c := range rl.clients {
			// Remove clients with no recent requests
			if len(c.requests) == 0 || now.Sub(c.requests[len(c.requests)-1]) > rl.duration {
				delete(rl.clients, clientIP)
			}
		}

		rl.mutex.Unlock()
	}
}

// Limit is the middleware function for rate limiting
func (rl *RateLimiter) Limit() gin.HandlerFunc {
	// Start cleanup goroutine
	go rl.CleanUp()

	return func(c *gin.Context) {
		clientIP := c.ClientIP()

		if !rl.Allow(clientIP) {
			c.JSON(http.StatusTooManyRequests, gin.H{
				"error":       "Rate limit exceeded",
				"retry_after": rl.duration.Seconds(),
			})
			c.Abort()
			return
		}

		c.Next()
	}
}
