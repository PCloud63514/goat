package cache

import (
	"context"
	"time"
)

type ActiveExpireCache[T any] struct {
	name            string
	store           *cacheStore[T, *activeExpireCacheItem[T]]
	keyGen          *keyGenerator
	metrics         *cacheMetrics
	ttl             time.Duration
	expireExtension bool
	ctx             context.Context
	samplingDelay   time.Duration
	samplingRatio   int
	samplingSize    int
}

type activeExpireCacheItem[T any] struct {
	key        string
	value      T
	expiration time.Time
}

func NewActiveExpireCache[T any](name string, capacity int, ttl time.Duration, expireExtension bool, ctx context.Context, samplingDelay time.Duration, samplingRatio int, samplingSize int) Cache[T] {
	if ctx == nil {
		panic("[ActiveExpireCache] The context must not be nil.")
	}
	cache := &ActiveExpireCache[T]{
		name:            name,
		store:           newCacheStore[T, *activeExpireCacheItem[T]](capacity),
		keyGen:          &keyGenerator{},
		metrics:         &cacheMetrics{},
		ttl:             ttl,
		ctx:             ctx,
		expireExtension: expireExtension,
		samplingDelay:   samplingDelay,
		samplingRatio:   samplingRatio,
		samplingSize:    samplingSize,
	}
	go cache.backgroundScan()
	return cache
}

func (c *ActiveExpireCache[T]) Name() string {
	return c.name
}

func (c *ActiveExpireCache[T]) Get(keys ...any) (T, bool) {
	key := c.keyGen.Generate(keys)
	return c.get(key)
}

func (c *ActiveExpireCache[T]) Put(keys []any, value T) {
	key := c.keyGen.Generate(keys)
	c.put(key, value)
}

func (c *ActiveExpireCache[T]) Delete(keys ...any) {
	key := c.keyGen.Generate(keys)
	c.delete(key)
}

func (c *ActiveExpireCache[T]) Stat() *CacheStat {
	return &CacheStat{
		Name:        c.name,
		MaxEntries:  c.store.Capacity(),
		CurrentSize: c.store.Size(),
		HitCount:    c.metrics.HitCount(),
		MissCount:   c.metrics.MissCount(),
		HitRate:     c.metrics.HitRate(),
	}
}

func (c *ActiveExpireCache[T]) get(key string) (T, bool) {
	if it, ok := c.store.Get(key); ok && it.IsNotExpired() {
		if c.expireExtension {
			it.ExtendExpiration(c.ttl)
		}
		c.metrics.Hit()
		return it.Value(), true
	}
	c.metrics.Miss()
	var zeroValue T
	return zeroValue, false
}

func (c *ActiveExpireCache[T]) put(key string, value T) {
	c.store.Put(&activeExpireCacheItem[T]{
		key:        key,
		value:      value,
		expiration: time.Now().Add(c.ttl),
	})
}

func (c *ActiveExpireCache[T]) delete(key string) {
	c.store.Delete(key)
}

func (ai *activeExpireCacheItem[T]) Key() string {
	return ai.key
}

func (ai *activeExpireCacheItem[T]) Value() T {
	return ai.value
}

func (ai *activeExpireCacheItem[T]) Expired() bool {
	return time.Now().After(ai.expiration)
}

func (ai *activeExpireCacheItem[T]) IsNotExpired() bool {
	return !time.Now().After(ai.expiration)
}

func (ai *activeExpireCacheItem[T]) ExtendExpiration(ttl time.Duration) {
	ai.expiration = time.Now().Add(ttl)
}

func (c *ActiveExpireCache[T]) backgroundScan() {
	ticker := time.NewTicker(c.samplingDelay)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			for {
				if ratio := c.cleanupExpiredItems(); ratio < c.samplingRatio {
					break
				}
			}
		case <-c.ctx.Done():
			return
		}
	}
}

func (c *ActiveExpireCache[T]) cleanupExpiredItems() int {
	sampledKeys := c.store.SampleKeys(c.samplingSize)
	expireCount := 0
	itemSize := 0
	for _, key := range sampledKeys {
		if item, ok := c.store.Get(key); ok {
			itemSize += 1
			if item.Expired() {
				c.store.Delete(key)
				expireCount += 1
			}
		}
	}
	if itemSize == 0 {
		return 0
	}
	return expireCount / itemSize
}
