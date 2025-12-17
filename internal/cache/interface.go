package cache

import "time"

// CacheInterface defines the cache interface
type CacheInterface interface {
	Set(key string, value interface{}, expire time.Duration) error
	Get(key string) (interface{}, bool)
	Delete(key string) error
	Clear() error
	Size() int
	Stats() (hits, misses uint64)
	HitRate() float64
	ResetStats()
	Close() error
}