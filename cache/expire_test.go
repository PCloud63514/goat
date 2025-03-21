package cache

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

type TestExpireCacheItem struct{}

func TestExpireCache_New(t *testing.T) {
	cache := NewExpireCache[*TestExpireCacheItem]("Name", 10, time.Second*3, false)
	assert.NotNil(t, cache, "Failed to create an instance.")
}

func TestExpireCache_Name(t *testing.T) {
	cache := NewExpireCache[*TestExpireCacheItem]("Name", 10, time.Second*3, false)
	assert.Equal(t, cache.Name(), "Name", "Name() = %s; want Name", cache.Name())
}

func TestExpireCache_Get(t *testing.T) {
	cache := NewExpireCache[*TestExpireCacheItem]("Name", 10, time.Second*1, false)

	t.Run("should return false if there is no data", func(t *testing.T) {
		_, ok := cache.Get("key")
		assert.Equal(t, ok, false, "Get() = (_, false); want (_, true)")
	})

	t.Run("should return true if there is data", func(t *testing.T) {
		cache.Put([]any{"key"}, &TestExpireCacheItem{})
		_, ok := cache.Get("key")
		assert.Equal(t, ok, true, "Get() = (_, false); want (_, true)")
	})

	t.Run("should return false if data is expired", func(t *testing.T) {
		cache.Put([]any{"key"}, &TestExpireCacheItem{})
		time.Sleep(time.Second * 2)
		_, ok := cache.Get("key")
		assert.Equal(t, ok, false, "Get() = (_, false); want (_, true)")
	})
}

func TestExpireCache_Put(t *testing.T) {
	cache := NewExpireCache[*TestExpireCacheItem]("Name", 10, time.Second*1, false)
	cache.Put([]any{"key"}, &TestExpireCacheItem{})
	_, ok := cache.Get("key")
	assert.Equal(t, ok, true, "Put() = (_, false); want (_, true)")
}
