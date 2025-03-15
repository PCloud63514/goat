package cache

import (
	"context"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

type TestActiveExpireCacheItem struct{}

func TestActiveExpireCache_New(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cache := NewActiveExpireCache[*TestActiveExpireCacheItem](ctx, "Name", 10, defaultKeyFunc, time.Second*2, false, time.Microsecond*100, 25, 100)
	assert.NotNil(t, cache, "Failed to create an instance.")
	cancel()
}

func TestActiveExpireCache_Name(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cache := NewActiveExpireCache[*TestActiveExpireCacheItem](ctx, "Name", 10, defaultKeyFunc, time.Second*2, false, time.Microsecond*100, 25, 100)
	assert.Equal(t, cache.Name(), "Name", "Name() = %s; want Name", cache.Name())
	cancel()
}

func TestActiveExpireCache_Get(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cache := NewActiveExpireCache[*TestActiveExpireCacheItem](ctx, "Name", 10, defaultKeyFunc, time.Second*2, false, time.Microsecond*100, 25, 100)

	t.Run("should return false if there is no data", func(t *testing.T) {
		_, ok := cache.Get("key")
		assert.Equal(t, ok, false, "Get() = (_, false); want (_, true)")
	})

	t.Run("should return true if there is data", func(t *testing.T) {
		cache.Put([]any{"key"}, &TestActiveExpireCacheItem{})
		_, ok := cache.Get("key")
		assert.Equal(t, ok, true, "Get() = (_, false); want (_, true)")
	})

	t.Run("should return false if data is expired", func(t *testing.T) {
		cache.Put([]any{"key"}, &TestActiveExpireCacheItem{})
		time.Sleep(time.Second * 4)
		_, ok := cache.Get("key")
		assert.Equal(t, ok, false, "Get() = (_, false); want (_, true)")
	})
	cancel()
}

func TestActiveExpireCache_Put(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cache := NewActiveExpireCache[*TestActiveExpireCacheItem](ctx, "Name", 10, defaultKeyFunc, time.Second*2, false, time.Microsecond*100, 25, 100)
	cache.Put([]any{"key"}, &TestActiveExpireCacheItem{})
	_, ok := cache.Get("key")
	assert.Equal(t, ok, true, "Put() = (_, false); want (_, true)")
	cancel()
}

func TestActiveExpireCache_Delete(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cache := NewActiveExpireCache[*TestActiveExpireCacheItem](ctx, "Name", 10, defaultKeyFunc, time.Second*2, false, time.Microsecond*100, 25, 100)
	cache.Put([]any{"key"}, &TestActiveExpireCacheItem{})
	cache.Delete("key")
	_, ok := cache.Get("key")
	assert.Equal(t, ok, false, "Delete() = (_, false); want (_, true)")
	cancel()
}
