package cache

import (
	"container/list"
	"context"
	"sync"
	"time"
)

type activeExpireCache[T any] struct {
	ctx             context.Context
	name            string
	mu              sync.RWMutex
	cache           map[string]*list.Element
	ll              *list.List
	cacheSize       int
	keyFunc         KeyFunc
	ttl             time.Duration
	expireExtension bool
	samplingDelay   time.Duration
	samplingRatio   int
	samplingSize    int
}

type activeExpireCacheItem[T any] struct {
	key        string
	value      T
	expiration time.Time
}

func NewActiveExpireCache[T any](ctx context.Context, name string, capacity int, keyFunc KeyFunc,
	ttl time.Duration, expireExtension bool,
	samplingDelay time.Duration, samplingRatio int, samplingSize int) Cache[T] {
	if ctx == nil {
		panic("[ActiveExpireCache] The context must not be nil.")
	}
	cache := &activeExpireCache[T]{
		ctx:             ctx,
		name:            name,
		mu:              sync.RWMutex{},
		cache:           make(map[string]*list.Element),
		ll:              list.New(),
		cacheSize:       capacity,
		keyFunc:         keyFunc,
		ttl:             ttl,
		expireExtension: expireExtension,
		samplingDelay:   samplingDelay,
		samplingRatio:   samplingRatio,
		samplingSize:    samplingSize,
	}
	go cache.backgroundScan()
	return cache
}

func (ae *activeExpireCache[T]) Name() string {
	return ae.name
}

func (ae *activeExpireCache[T]) Get(keys ...any) (T, bool) {
	key := ae.keyFunc(keys...)
	return ae.get(key)
}

func (ae *activeExpireCache[T]) Put(keys []any, value T) {
	key := ae.keyFunc(keys...)
	ae.put(key, value)
}

func (ae *activeExpireCache[T]) Delete(keys ...any) {
	key := ae.keyFunc(keys...)
	ae.delete(key)
}

func (ae *activeExpireCache[T]) get(key string) (T, bool) {
	ae.mu.RLock()
	var zeroValue T
	if elem, found := ae.cache[key]; found {
		ae.ll.MoveToFront(elem)
		item := elem.Value.(*activeExpireCacheItem[T])
		if item.isExpired() {
			ae.mu.RUnlock()
			return zeroValue, false
		}
		ae.mu.RUnlock()
		if ae.expireExtension {
			ae.mu.Lock()
			defer ae.mu.Unlock()
			item.expiration = time.Now().Add(ae.ttl)
		}
		return item.value, true
	}
	defer ae.mu.RUnlock()
	return zeroValue, false
}

func (ae *activeExpireCache[T]) put(key string, value T) {
	ae.mu.Lock()
	defer ae.mu.Unlock()
	expiration := time.Now().Add(ae.ttl)
	if elem, found := ae.cache[key]; found {
		ae.ll.MoveToFront(elem)
		item := elem.Value.(*activeExpireCacheItem[T])
		item.value = value
		item.expiration = expiration
		return
	}

	elem := ae.ll.PushFront(&activeExpireCacheItem[T]{key, value, expiration})
	ae.cache[key] = elem

	if ae.ll.Len() > ae.cacheSize {
		last := ae.ll.Back()
		if last != nil {
			ae.ll.Remove(last)
			delete(ae.cache, last.Value.(*activeExpireCacheItem[T]).key)
		}
	}
}

func (ae *activeExpireCache[T]) delete(key string) {
	ae.mu.Lock()
	defer ae.mu.Unlock()

	if elem, found := ae.cache[key]; found {
		ae.ll.Remove(elem)
		delete(ae.cache, key)
	}
}

func (ai *activeExpireCacheItem[T]) isExpired() bool {
	return time.Now().After(ai.expiration)
}

func (ai *activeExpireCacheItem[T]) isNotExpired() bool {
	return !time.Now().After(ai.expiration)
}

func (ae *activeExpireCache[T]) backgroundScan() {
	ticker := time.NewTicker(time.Millisecond * ae.samplingDelay)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			if ratio := ae.cleanupExpiredItems(); ae.samplingRatio <= ratio {
				continue
			}
		case <-ae.ctx.Done():
			return
		}
	}
}

func (ae *activeExpireCache[T]) cleanupExpiredItems() int {
	sampledKeys := ae.sampleKeys()
	expireCount := 0
	itemSize := 0
	ae.mu.Lock()
	defer ae.mu.Unlock()
	now := time.Now()
	for _, key := range sampledKeys {
		if item, exists := ae.cache[key]; exists {
			itemSize += 1
			if now.After(item.Value.(*activeExpireCacheItem[T]).expiration) {
				ae.ll.Remove(item)
				delete(ae.cache, key)
				expireCount += 1
			}
		}
	}
	return expireCount / itemSize
}

func (ae *activeExpireCache[T]) sampleKeys() []string {
	ae.mu.RLock()
	defer ae.mu.RUnlock()

	sampledKeys := make([]string, 0, ae.samplingSize)
	sampled := 0

	for elem := ae.ll.Front(); elem != nil && sampled < ae.samplingSize; elem = elem.Next() {
		item := elem.Value.(*activeExpireCacheItem[T])
		sampledKeys = append(sampledKeys, item.key)
		sampled++
	}
	return sampledKeys
}
