package service

import (
	"runtime"
	"sync"
	"time"

	"go-admin/internal/cache"
	"go-admin/internal/database"
)

// MetricsService collects system metrics
type MetricsService struct {
	mu sync.RWMutex
}

// Metrics represents system metrics
type Metrics struct {
	Timestamp       time.Time `json:"timestamp"`
	CPUUsage        float64   `json:"cpu_usage"`
	MemoryAllocated uint64    `json:"memory_allocated"`
	MemorySystem    uint64    `json:"memory_system"`
	Goroutines      int       `json:"goroutines"`
	CacheHits       uint64    `json:"cache_hits"`
	CacheMisses     uint64    `json:"cache_misses"`
	CacheHitRate    float64   `json:"cache_hit_rate"`
	DBMaxOpenConns  int       `json:"db_max_open_conns"`
	DBOpenConns     int       `json:"db_open_conns"`
	DBInUseConns    int       `json:"db_in_use_conns"`
	DBIdleConns     int       `json:"db_idle_conns"`
}

// NewMetricsService creates a new metrics service
func NewMetricsService() *MetricsService {
	return &MetricsService{}
}

// CollectMetrics collects current system metrics
func (s *MetricsService) CollectMetrics() *Metrics {
	s.mu.Lock()
	defer s.mu.Unlock()

	// Get memory stats
	var memStats runtime.MemStats
	runtime.ReadMemStats(&memStats)

	// Get cache stats
	cacheInstance := cache.GetInstance()
	cacheHits, cacheMisses := cacheInstance.Stats()
	cacheHitRate := cacheInstance.HitRate()

	// Get database stats
	dbStats := database.Stats()

	// Create metrics object
	metrics := &Metrics{
		Timestamp:       time.Now(),
		CPUUsage:        0, // In a real implementation, you would collect actual CPU usage
		MemoryAllocated: memStats.Alloc,
		MemorySystem:    memStats.Sys,
		Goroutines:      runtime.NumGoroutine(),
		CacheHits:       cacheHits,
		CacheMisses:     cacheMisses,
		CacheHitRate:    cacheHitRate,
		DBMaxOpenConns:  dbStats.MaxOpenConnections,
		DBOpenConns:     dbStats.OpenConnections,
		DBInUseConns:    dbStats.InUse,
		DBIdleConns:     dbStats.Idle,
	}

	return metrics
}

// GetHealthStatus returns the health status of the application
func (s *MetricsService) GetHealthStatus() map[string]interface{} {
	// Check database connectivity
	dbHealthy := true
	dbStats := database.Stats()
	if dbStats.OpenConnections == 0 {
		dbHealthy = false
	}

	// Check cache availability
	cacheInstance := cache.GetInstance()
	cacheHealthy := cacheInstance != nil

	// Get uptime
	// In a real implementation, you would track application start time

	return map[string]interface{}{
		"status":    "healthy",
		"timestamp": time.Now(),
		"components": map[string]interface{}{
			"database": map[string]interface{}{
				"status":  dbHealthy,
				"details": dbStats,
			},
			"cache": map[string]interface{}{
				"status": cacheHealthy,
			},
			"system": map[string]interface{}{
				"goroutines": runtime.NumGoroutine(),
				"memory": map[string]interface{}{
					"allocated": "N/A", // In a real implementation, you would format bytes appropriately
					"system":    "N/A",
				},
			},
		},
	}
}
