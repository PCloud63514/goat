package cache

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func Test_ExpireCache(t *testing.T) {
	t.Run("should store put and get value", func(t *testing.T) {
		cache := NewExpireCache[string]("Name", 10, time.Second*10, false)
		cache.Put([]any{"key"}, "Test1")
		v, ok := cache.Get("key")
		assert.True(t, ok)
		assert.Equal(t, "Test1", v)
	})

	t.Run("should return false if key not found", func(t *testing.T) {
		cache := NewExpireCache[string]("Name", 10, time.Second*10, false)
		_, ok := cache.Get("key")
		assert.False(t, ok)
	})

	t.Run("should delete storeItem", func(t *testing.T) {
		cache := NewExpireCache[string]("Name", 10, time.Second*10, false)
		cache.Put([]any{"key"}, "Test1")
		cache.Delete("key")
		_, ok := cache.Get("key")
		assert.False(t, ok)
	})

	t.Run("should return stat", func(t *testing.T) {
		cache := NewExpireCache[string]("Name", 10, time.Second*10, false)
		cache.Put([]any{"key"}, "Test1")
		cache.Put([]any{"key2"}, "Test2")
		cache.Get("key")
		cache.Get("key3")

		stat := cache.Stat()

		assert.Equal(t, cache.Name(), stat.Name)
		assert.Equal(t, 10, stat.MaxEntries)
		assert.Equal(t, 2, stat.CurrentSize)
		assert.Equal(t, uint64(1), stat.HitCount)
		assert.Equal(t, uint64(1), stat.MissCount)
		assert.Equal(t, float64(50), stat.HitRate)
	})

	t.Run("should return false if capacity is full", func(t *testing.T) {
		cache := NewExpireCache[string]("Name", 1, time.Second*10, false)
		cache.Put([]any{"key"}, "Test1")
		cache.Put([]any{"key2"}, "Test2")
		_, ok := cache.Get("key")
		assert.False(t, ok)
	})

	t.Run("should return false if key is expired", func(t *testing.T) {
		cache := NewExpireCache[string]("Name", 10, time.Second*1, false)
		cache.Put([]any{"key"}, "Test1")
		time.Sleep(time.Second * 2)
		_, ok := cache.Get("key")
		assert.False(t, ok)
	})
}
