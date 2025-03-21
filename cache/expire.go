package cache

import (
	"container/list"
	"sync"
	"time"
)

type expireCache[T any] struct {
	name            string
	mu              sync.RWMutex
	cache           map[string]*list.Element
	ll              *list.List
	cacheSize       int
	keyFunc         KeyFunc
	ttl             time.Duration
	expireExtension bool
	metrics         *cacheMetrics
}

type expireCacheItem[T any] struct {
	key        string
	value      T
	expiration time.Time
}

func NewExpireCache[T any](name string, capacity int, keyFunc KeyFunc, ttl time.Duration, expireExtension bool) Cache[T] {
	return &expireCache[T]{
		name:            name,
		mu:              sync.RWMutex{},
		cache:           make(map[string]*list.Element),
		ll:              list.New(),
		cacheSize:       capacity,
		keyFunc:         keyFunc,
		ttl:             ttl,
		expireExtension: expireExtension,
		metrics:         &cacheMetrics{},
	}
}

func (ex *expireCache[T]) Name() string {
	return ex.name
}

func (ex *expireCache[T]) Get(keys ...any) (T, bool) {
	key := ex.keyFunc(keys...)
	return ex.get(key)
}

func (ex *expireCache[T]) Put(keys []any, value T) {
	key := ex.keyFunc(keys...)
	ex.put(key, value)
}

func (ex *expireCache[T]) Delete(keys ...any) {
	key := ex.keyFunc(keys...)
	ex.delete(key)
}

func (ex *expireCache[T]) get(key string) (T, bool) {
	ex.mu.RLock()
	var zeroValue T
	if elem, found := ex.cache[key]; found {
		ex.ll.MoveToFront(elem)
		item := elem.Value.(*expireCacheItem[T])
		if item.isExpired() {
			ex.mu.RUnlock()
			ex.mu.Lock()
			defer ex.mu.Unlock()
			ex.ll.Remove(elem)
			delete(ex.cache, key)
			ex.metrics.Miss()
			return zeroValue, false
		}
		ex.mu.RUnlock()
		if ex.expireExtension {
			ex.mu.Lock()
			defer ex.mu.Unlock()
			item.expiration = time.Now().Add(ex.ttl)
		}
		ex.metrics.Hit()
		return item.value, true
	}
	ex.metrics.Miss()
	defer ex.mu.RUnlock()
	return zeroValue, false
}

func (ex *expireCache[T]) put(key string, value T) {
	ex.mu.Lock()
	defer ex.mu.Unlock()
	expiration := time.Now().Add(ex.ttl)
	if elem, found := ex.cache[key]; found {
		ex.ll.MoveToFront(elem)
		item := elem.Value.(*expireCacheItem[T])
		item.value = value
		item.expiration = expiration
		return
	}

	elem := ex.ll.PushFront(&expireCacheItem[T]{key, value, expiration})
	ex.cache[key] = elem

	if ex.ll.Len() > ex.cacheSize {
		last := ex.ll.Back()
		if last != nil {
			ex.ll.Remove(last)
			delete(ex.cache, last.Value.(*expireCacheItem[T]).key)
		}
	}
}

func (ex *expireCache[T]) delete(key string) {
	ex.mu.Lock()
	defer ex.mu.Unlock()

	if elem, found := ex.cache[key]; found {
		ex.ll.Remove(elem)
		delete(ex.cache, key)
	}
}

func (ei *expireCacheItem[T]) isExpired() bool {
	return time.Now().After(ei.expiration)
}

func (ei *expireCacheItem[T]) isNotExpired() bool {
	return !time.Now().After(ei.expiration)
}
