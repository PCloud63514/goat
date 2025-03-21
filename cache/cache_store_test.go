package cache

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

const testCapacity = 100

type TestItem[T any] struct {
	key   string
	value T
}

func (t *TestItem[T]) Key() string {
	return t.key
}
func (t *TestItem[T]) Value() T {
	return t.value
}

func Test_CacheStore(t *testing.T) {
	t.Run("should store put and get value", func(t *testing.T) {
		store := newCacheStore[string, *TestItem[string]](testCapacity)
		item := &TestItem[string]{"key", "Test1"}
		store.Put(item)
		v, ok := store.Get("key")
		assert.True(t, ok)
		assert.Equal(t, "Test1", v.Value())
	})

	t.Run("should return false if key not found", func(t *testing.T) {
		store := newCacheStore[string, *TestItem[string]](testCapacity)
		_, ok := store.Get("key")
		assert.False(t, ok)
	})

	t.Run("should delete storeItem", func(t *testing.T) {
		store := newCacheStore[string, *TestItem[string]](testCapacity)
		item := &TestItem[string]{"key", "Test1"}
		store.Put(item)
		store.Delete("key")
		_, ok := store.Get("key")
		assert.False(t, ok)
	})

	t.Run("should return size", func(t *testing.T) {
		store := newCacheStore[string, *TestItem[string]](testCapacity)
		item := &TestItem[string]{"key", "Test1"}
		store.Put(item)
		assert.Equal(t, 1, store.Size())
	})

	t.Run("should return capacity", func(t *testing.T) {
		store := newCacheStore[string, *TestItem[string]](testCapacity)
		assert.Equal(t, testCapacity, store.Capacity())
	})

	t.Run("should return false if capacity is full", func(t *testing.T) {
		store := newCacheStore[string, *TestItem[string]](1)
		item1 := &TestItem[string]{"key", "Test1"}
		store.Put(item1)
		item2 := &TestItem[string]{"key2", "Test2"}
		store.Put(item2)
		_, ok := store.Get("key")
		assert.False(t, ok)
		v, ok := store.Get("key2")
		assert.True(t, ok)
		assert.Equal(t, item2.Value(), v.Value())
	})

	t.Run("should update storeItem", func(t *testing.T) {
		store := newCacheStore[string, *TestItem[string]](testCapacity)
		item := &TestItem[string]{"key", "Test1"}
		store.Put(item)
		item = &TestItem[string]{"key", "Test2"}
		store.Put(item)
		v, ok := store.Get("key")
		assert.True(t, ok)
		assert.Equal(t, "Test2", v.Value())
	})

	t.Run("should return false if key not found", func(t *testing.T) {
		store := newCacheStore[string, *TestItem[string]](testCapacity)
		store.Put(&TestItem[string]{"key", "Test1"})
		store.Delete("key")
		_, ok := store.Get("key")
		assert.False(t, ok)
	})
}
