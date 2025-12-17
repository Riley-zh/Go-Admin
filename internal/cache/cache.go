package cache

import (
	"sync"
	"sync/atomic"
	"time"

	"go-admin/config"
)

// Cache represents a thread-safe cache with expiration support
type Cache struct {
	data        sync.Map
	expireTimes sync.Map
	maxSize     int

	// Statistics
	hits   uint64
	misses uint64
}

var (
	instance *Cache
	once     sync.Once
)

// Init initializes the cache with the given configuration
func Init(cfg config.CacheConfig) {
	once.Do(func() {
		instance = &Cache{
			maxSize: cfg.MaxSize,
		}

		// Start garbage collection goroutine
		go instance.startGC(cfg.GCInterval)
	})
}

// GetInstance returns the singleton cache instance
func GetInstance() *Cache {
	return instance
}

// Set stores a value in the cache with an expiration time
func (c *Cache) Set(key string, value interface{}, expire time.Duration) {
	c.data.Store(key, value)

	if expire > 0 {
		expireTime := time.Now().Add(expire)
		c.expireTimes.Store(key, expireTime)
	} else {
		c.expireTimes.Delete(key)
	}
}

// Get retrieves a value from the cache
func (c *Cache) Get(key string) (interface{}, bool) {
	value, ok := c.data.Load(key)
	if !ok {
		atomic.AddUint64(&c.misses, 1)
		return nil, false
	}

	// Check if the item has expired
	if expireTime, ok := c.expireTimes.Load(key); ok {
		if time.Now().After(expireTime.(time.Time)) {
			c.Delete(key)
			atomic.AddUint64(&c.misses, 1)
			return nil, false
		}
	}

	atomic.AddUint64(&c.hits, 1)
	return value, true
}

// Delete removes a value from the cache
func (c *Cache) Delete(key string) {
	c.data.Delete(key)
	c.expireTimes.Delete(key)
}

// Clear removes all values from the cache
func (c *Cache) Clear() {
	c.data.Range(func(key, _ interface{}) bool {
		c.data.Delete(key)
		return true
	})

	c.expireTimes.Range(func(key, _ interface{}) bool {
		c.expireTimes.Delete(key)
		return true
	})
}

// Size returns the number of items in the cache
func (c *Cache) Size() int {
	size := 0
	c.data.Range(func(_, _ interface{}) bool {
		size++
		return true
	})
	return size
}

// Stats returns cache statistics
func (c *Cache) Stats() (hits, misses uint64) {
	hits = atomic.LoadUint64(&c.hits)
	misses = atomic.LoadUint64(&c.misses)
	return hits, misses
}

// HitRate returns the cache hit rate
func (c *Cache) HitRate() float64 {
	hits, misses := c.Stats()
	total := hits + misses
	if total == 0 {
		return 0
	}
	return float64(hits) / float64(total)
}

// ResetStats resets the cache statistics
func (c *Cache) ResetStats() {
	atomic.StoreUint64(&c.hits, 0)
	atomic.StoreUint64(&c.misses, 0)
}

// startGC starts the garbage collection process
func (c *Cache) startGC(interval time.Duration) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for range ticker.C {
		c.GC()
	}
}

// GC removes expired items from the cache
func (c *Cache) GC() {
	now := time.Now()

	c.expireTimes.Range(func(key, expireTime interface{}) bool {
		if now.After(expireTime.(time.Time)) {
			c.Delete(key.(string))
		}
		return true
	})
}
