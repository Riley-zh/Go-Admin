package handler

import (
	"go-admin/internal/cache"
	"go-admin/pkg/response"

	"github.com/gin-gonic/gin"
)

// CacheHandler handles cache-related HTTP requests
type CacheHandler struct{}

// NewCacheHandler creates a new cache handler
func NewCacheHandler() *CacheHandler {
	return &CacheHandler{}
}

// GetCacheStats handles requests to get cache statistics
func (h *CacheHandler) GetCacheStats(c *gin.Context) {
	cacheInstance := cache.GetInstance()

	hits, misses := cacheInstance.Stats()
	hitRate := cacheInstance.HitRate()
	size := cacheInstance.Size()

	response.Success(c, "Cache statistics retrieved successfully", gin.H{
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

	response.Success(c, "Cache statistics reset successfully", nil)
}

// ClearCache handles requests to clear all cache entries
func (h *CacheHandler) ClearCache(c *gin.Context) {
	cacheInstance := cache.GetInstance()
	cacheInstance.Clear()

	response.Success(c, "Cache cleared successfully", nil)
}
