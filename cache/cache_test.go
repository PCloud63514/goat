package cache

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

type TestItem struct{}

func TestLRUCache_New(t *testing.T) {
	cache := NewLRUCache[*TestItem]("Name", 10, defaultKeyFUnc)
	assert.NotNil(t, cache, "Failed to create an instance.")
}

func TestLRUCache_Name(t *testing.T) {
	cache := NewLRUCache[*TestItem]("Name", 10, defaultKeyFUnc)
	assert.Equal(t, cache.Name(), "Name", "Name() = %s; want Name", cache.Name())
}

func TestLRUCache_Get(t *testing.T) {
	cache := NewLRUCache[*TestItem]("Name", 10, defaultKeyFUnc)
	t.Run("should return false if there is no data", func(t *testing.T) {
		_, ok := cache.Get("key")
		assert.Equal(t, ok, false, "Get() = (value, ok); want (_, true)")
	})

	t.Run("should return true if there is data", func(t *testing.T) {
		givenItem := &TestItem{}
		cache.Put([]any{"key"}, &TestItem{})
		v, ok := cache.Get("key")
		assert.Equal(t, ok, true, "Get() = (_, false); want (_, true)")
		assert.Equal(t, v, givenItem, "Get() = (_, false); want (_, true)")
	})
}

func TestLRUCache_Put(t *testing.T) {
	givenItem := &TestItem{}
	cache := NewLRUCache[*TestItem]("Name", 10, defaultKeyFUnc)
	cache.Put([]any{"key"}, &TestItem{})
	v, ok := cache.Get("key")
	assert.Equal(t, ok, true, "Put() = (_, false); want (_, true)")
	assert.Equal(t, v, givenItem, "Put() = (_, false); want (_, true)")
}

func TestLRUCache_Delete(t *testing.T) {
	cache := NewLRUCache[*TestItem]("Name", 10, defaultKeyFUnc)
	cache.Put([]any{"key"}, &TestItem{})
	cache.Delete("key")
	_, ok := cache.Get("key")
	assert.Equal(t, ok, false, "Delete() = (_, false); want (_, true)")
}
