package cache

import "time"

type ExpireCache[T any] struct {
	name            string
	store           *cacheStore[T, *expireCacheItem[T]]
	keyGen          *keyGenerator
	metrics         *cacheMetrics
	ttl             time.Duration
	expireExtension bool
}

type expireCacheItem[T any] struct {
	key        string
	value      T
	expiration time.Time
}

func NewExpireCache[T any](name string, capacity int, ttl time.Duration, expireExtension bool) Cache[T] {
	return &ExpireCache[T]{
		name:            name,
		store:           newCacheStore[T, *expireCacheItem[T]](capacity),
		keyGen:          &keyGenerator{},
		metrics:         &cacheMetrics{},
		ttl:             ttl,
		expireExtension: expireExtension,
	}
}

func (c *ExpireCache[T]) Name() string {
	return c.name
}

func (c *ExpireCache[T]) Get(keys ...any) (T, bool) {
	key := c.keyGen.Generate(keys)
	return c.get(key)
}

func (c *ExpireCache[T]) Put(keys []any, value T) {
	key := c.keyGen.Generate(keys)
	c.put(key, value)
}

func (c *ExpireCache[T]) Delete(keys ...any) {
	key := c.keyGen.Generate(keys)
	c.delete(key)
}

func (c *ExpireCache[T]) Stat() *CacheStat {
	return &CacheStat{
		Name:        c.name,
		MaxEntries:  c.store.Capacity(),
		CurrentSize: c.store.Size(),
		HitCount:    c.metrics.HitCount(),
		MissCount:   c.metrics.MissCount(),
		HitRate:     c.metrics.HitRate(),
	}
}

func (c *ExpireCache[T]) get(key string) (T, bool) {
	if it, ok := c.store.Get(key); ok {
		if it.Expired() {
			c.store.Delete(key)
			c.metrics.Miss()
			var zeroValue T
			return zeroValue, false
		}
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

func (c *ExpireCache[T]) put(key string, value T) {
	c.store.Put(&expireCacheItem[T]{key: key, value: value, expiration: time.Now().Add(c.ttl)})
}

func (c *ExpireCache[T]) delete(key string) {
	c.store.Delete(key)
}

func (ei *expireCacheItem[T]) Key() string {
	return ei.key
}

func (ei *expireCacheItem[T]) Value() T {
	return ei.value
}

func (ei *expireCacheItem[T]) Expired() bool {
	return time.Now().After(ei.expiration)
}

func (ei *expireCacheItem[T]) IsNotExpired() bool {
	return !time.Now().After(ei.expiration)
}

func (ei *expireCacheItem[T]) ExtendExpiration(ttl time.Duration) {
	ei.expiration = time.Now().Add(ttl)
}
