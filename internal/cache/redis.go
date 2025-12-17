package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"sync/atomic"
	"time"

	"go-admin/config"
	"go-admin/internal/logger"

	"github.com/go-redis/redis/v8"
	"go.uber.org/zap"
)

// RedisCache represents a Redis-based cache
type RedisCache struct {
	client *redis.Client
	ctx    context.Context

	// Statistics
	hits   uint64
	misses uint64
}

// NewRedisCache creates a new Redis cache instance
func NewRedisCache(cfg config.RedisConfig) *RedisCache {
	rdb := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", cfg.Host, cfg.Port),
		Password: cfg.Password,
		DB:       cfg.DB,
		PoolSize: cfg.PoolSize,
	})

	ctx := context.Background()

	// Test connection
	_, err := rdb.Ping(ctx).Result()
	if err != nil {
		panic(fmt.Sprintf("Failed to connect to Redis: %v", err))
	}

	return &RedisCache{
		client: rdb,
		ctx:    ctx,
	}
}

// Set stores a value in the cache with an expiration time
func (r *RedisCache) Set(key string, value interface{}, expire time.Duration) error {
	// Serialize value
	data, err := json.Marshal(value)
	if err != nil {
		return fmt.Errorf("failed to serialize value: %w", err)
	}

	// Store in Redis
	if expire > 0 {
		return r.client.Set(r.ctx, key, data, expire).Err()
	}
	return r.client.Set(r.ctx, key, data, 0).Err()
}

// Get retrieves a value from the cache
func (r *RedisCache) Get(key string) (interface{}, bool) {
	// Get from Redis
	result, err := r.client.Get(r.ctx, key).Result()
	if err != nil {
		if err == redis.Nil {
			atomic.AddUint64(&r.misses, 1)
			return nil, false
		}
		// Log error but continue
		logger.Error("Redis get error", zap.Error(err), zap.String("key", key))
		atomic.AddUint64(&r.misses, 1)
		return nil, false
	}

	// Deserialize value
	var value interface{}
	err = json.Unmarshal([]byte(result), &value)
	if err != nil {
		// Log error but continue
		logger.Error("Failed to deserialize cached value", zap.Error(err), zap.String("key", key))
		atomic.AddUint64(&r.misses, 1)
		return nil, false
	}

	atomic.AddUint64(&r.hits, 1)
	return value, true
}

// Delete removes a value from the cache
func (r *RedisCache) Delete(key string) error {
	return r.client.Del(r.ctx, key).Err()
}

// Clear removes all values from the cache
func (r *RedisCache) Clear() error {
	return r.client.FlushDB(r.ctx).Err()
}

// Size returns the number of items in the cache
func (r *RedisCache) Size() int {
	size, err := r.client.DBSize(r.ctx).Result()
	if err != nil {
		return 0
	}
	return int(size)
}

// Stats returns cache statistics
func (r *RedisCache) Stats() (hits, misses uint64) {
	hits = atomic.LoadUint64(&r.hits)
	misses = atomic.LoadUint64(&r.misses)
	return hits, misses
}

// HitRate returns the cache hit rate
func (r *RedisCache) HitRate() float64 {
	hits, misses := r.Stats()
	total := hits + misses
	if total == 0 {
		return 0
	}
	return float64(hits) / float64(total)
}

// ResetStats resets the cache statistics
func (r *RedisCache) ResetStats() {
	atomic.StoreUint64(&r.hits, 0)
	atomic.StoreUint64(&r.misses, 0)
}

// Close closes the cache and cleans up resources
func (r *RedisCache) Close() error {
	return r.client.Close()
}

// Health checks the Redis connection health
func (r *RedisCache) Health() error {
	_, err := r.client.Ping(r.ctx).Result()
	return err
}