package cache

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestCache_SetAndGet(t *testing.T) {
	// Create a new cache instance
	cache := &Cache{}

	// Test setting and getting a value
	key := "test_key"
	value := "test_value"
	cache.Set(key, value, 0)

	retrievedValue, exists := cache.Get(key)
	assert.True(t, exists)
	assert.Equal(t, value, retrievedValue)

	// Test getting a non-existent key
	nonExistentValue, exists := cache.Get("non_existent_key")
	assert.False(t, exists)
	assert.Nil(t, nonExistentValue)
}

func TestCache_Delete(t *testing.T) {
	// Create a new cache instance
	cache := &Cache{}

	// Set a value
	key := "test_key"
	value := "test_value"
	cache.Set(key, value, 0)

	// Delete the value
	cache.Delete(key)

	// Try to get the deleted value
	_, exists := cache.Get(key)
	assert.False(t, exists)
}

func TestCache_Expiration(t *testing.T) {
	// Create a new cache instance
	cache := &Cache{}

	// Set a value with short expiration
	key := "test_key"
	value := "test_value"
	cache.Set(key, value, 10*time.Millisecond)

	// Value should exist immediately
	retrievedValue, exists := cache.Get(key)
	assert.True(t, exists)
	assert.Equal(t, value, retrievedValue)

	// Wait for expiration
	time.Sleep(20 * time.Millisecond)

	// Value should no longer exist
	_, exists = cache.Get(key)
	assert.False(t, exists)
}

func TestCache_Size(t *testing.T) {
	// Create a new cache instance
	cache := &Cache{}

	// Initially cache should be empty
	assert.Equal(t, 0, cache.Size())

	// Add some items
	cache.Set("key1", "value1", 0)
	cache.Set("key2", "value2", 0)
	cache.Set("key3", "value3", 0)

	// Check size
	assert.Equal(t, 3, cache.Size())

	// Delete an item
	cache.Delete("key2")

	// Check size again
	assert.Equal(t, 2, cache.Size())
}

func TestCache_Clear(t *testing.T) {
	// Create a new cache instance
	cache := &Cache{}

	// Add some items
	cache.Set("key1", "value1", 0)
	cache.Set("key2", "value2", 0)
	cache.Set("key3", "value3", 0)

	// Clear cache
	cache.Clear()

	// Cache should be empty
	assert.Equal(t, 0, cache.Size())

	// All keys should be gone
	_, exists := cache.Get("key1")
	assert.False(t, exists)
	_, exists = cache.Get("key2")
	assert.False(t, exists)
	_, exists = cache.Get("key3")
	assert.False(t, exists)
}

func TestCache_Stats(t *testing.T) {
	// Create a new cache instance
	cache := &Cache{}

	// Initially stats should be zero
	hits, misses := cache.Stats()
	assert.Equal(t, uint64(0), hits)
	assert.Equal(t, uint64(0), misses)
	assert.Equal(t, float64(0), cache.HitRate())

	// Set a value
	cache.Set("key", "value", 0)

	// Get existing value (hit)
	cache.Get("key")
	hits, misses = cache.Stats()
	assert.Equal(t, uint64(1), hits)
	assert.Equal(t, uint64(0), misses)
	assert.Equal(t, float64(1), cache.HitRate())

	// Get non-existing value (miss)
	cache.Get("nonexistent")
	hits, misses = cache.Stats()
	assert.Equal(t, uint64(1), hits)
	assert.Equal(t, uint64(1), misses)
	assert.Equal(t, float64(0.5), cache.HitRate())

	// Reset stats
	cache.ResetStats()
	hits, misses = cache.Stats()
	assert.Equal(t, uint64(0), hits)
	assert.Equal(t, uint64(0), misses)
	assert.Equal(t, float64(0), cache.HitRate())
}
