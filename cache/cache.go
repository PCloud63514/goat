package cache

import (
	"container/list"
	"sync"
)

type Cache[T any] interface {
	Name() string
	Get(keys ...any) (T, bool)
	Put(keys []any, value T)
	Delete(keys ...any)
}

type Item[T any] interface {
	Key() string
	Value() T
}

type lruCache[T any] struct {
	name      string
	mu        sync.RWMutex
	cache     map[string]*list.Element
	ll        *list.List
	cacheSize int
	keyFunc   KeyFunc
}

type cacheItem[T any] struct {
	key   string
	value T
}

func newLRUCache[T any](cacheName string, cacheSize int, keyFunc KeyFunc) *lruCache[T] {

	return &lruCache[T]{
		name:      cacheName,
		mu:        sync.RWMutex{},
		cache:     make(map[string]*list.Element),
		ll:        list.New(),
		cacheSize: cacheSize,
		keyFunc:   keyFunc,
	}
}

func (c *lruCache[T]) Name() string {
	return c.name
}

func (c *lruCache[T]) Get(keys ...any) (T, bool) {
	key := c.keyFunc(keys...)
	return c.get(key)
}

func (c *lruCache[T]) Put(keys []any, value T) {
	key := c.keyFunc(keys...)
	c.put(key, value)
}

func (c *lruCache[T]) Delete(keys ...any) {
	key := c.keyFunc(keys...)
	c.delete(key)
}

func (c *lruCache[T]) get(key string) (T, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	var zeroValue T
	if elem, found := c.cache[key]; found {
		c.ll.MoveToFront(elem)
		return elem.Value.(*cacheItem[T]).value, true
	}
	return zeroValue, false
}

func (c *lruCache[T]) put(key string, value T) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if elem, found := c.cache[key]; found {
		c.ll.MoveToFront(elem)
		item := elem.Value.(*cacheItem[T])
		item.value = value
		return
	}

	elem := c.ll.PushFront(&cacheItem[T]{key, value})
	c.cache[key] = elem

	if c.ll.Len() > c.cacheSize {
		last := c.ll.Back()
		if last != nil {
			c.ll.Remove(last)
			delete(c.cache, last.Value.(*cacheItem[T]).key)
		}
	}
}

func (c *lruCache[T]) delete(key string) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if elem, found := c.cache[key]; found {
		c.ll.Remove(elem)
		delete(c.cache, key)
	}
}
