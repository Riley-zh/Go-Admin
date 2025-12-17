package handler

import (
	"go-admin/internal/cache"

	"github.com/gin-gonic/gin"
)

// CacheHandler handles cache-related HTTP requests
type CacheHandler struct {
	*BaseHandler
}

// NewCacheHandler creates a new cache handler
func NewCacheHandler() *CacheHandler {
	return &CacheHandler{
		BaseHandler: NewBaseHandler(),
	}
}

// GetCacheStats handles requests to get cache statistics
func (h *CacheHandler) GetCacheStats(c *gin.Context) {
	cacheInstance := cache.GetInstance()

	hits, misses := cacheInstance.Stats()
	hitRate := cacheInstance.HitRate()
	size := cacheInstance.Size()

	h.HandleSuccess(c, gin.H{
		"hits":     hits,
		"misses":   misses,
		"hit_rate": hitRate,
		"size":     size,
	})
}

// ResetCacheStats handles requests to reset cache statistics
func (h *CacheHandler) ResetCacheStats(c *gin.Context) {
	cacheInstance := cache.GetInstance()
	cacheInstance.ResetStats()

	h.HandleSuccess(c, gin.H{"message": "Cache statistics reset successfully"})
}

// ClearCache handles requests to clear all cache entries
func (h *CacheHandler) ClearCache(c *gin.Context) {
	cacheInstance := cache.GetInstance()
	cacheInstance.Clear()

	h.HandleSuccess(c, gin.H{"message": "Cache cleared successfully"})
}
