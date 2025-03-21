package cache

type LRUCache[T any] struct {
	name    string
	store   *cacheStore[T, *lruCacheItem[T]]
	keyGen  *keyGenerator
	metrics *cacheMetrics
}

type lruCacheItem[T any] struct {
	key   string
	value T
}

func NewLRUCache[T any](name string, capacity int) Cache[T] {
	return &LRUCache[T]{
		name:    name,
		store:   newCacheStore[T, *lruCacheItem[T]](capacity),
		keyGen:  &keyGenerator{},
		metrics: &cacheMetrics{},
	}
}

func (c *LRUCache[T]) Name() string {
	return c.name
}

func (c *LRUCache[T]) Get(keys ...any) (T, bool) {
	key := c.keyGen.Generate(keys)
	return c.get(key)
}

func (c *LRUCache[T]) Put(keys []any, value T) {
	key := c.keyGen.Generate(keys)
	c.put(key, value)
}

func (c *LRUCache[T]) Delete(keys ...any) {
	key := c.keyGen.Generate(keys)
	c.delete(key)
}

func (c *LRUCache[T]) Stat() *CacheStat {
	return &CacheStat{
		Name:        c.name,
		MaxEntries:  c.store.Capacity(),
		CurrentSize: c.store.Size(),
		HitCount:    c.metrics.HitCount(),
		MissCount:   c.metrics.MissCount(),
		HitRate:     c.metrics.HitRate(),
	}
}

func (i *lruCacheItem[T]) Key() string {
	return i.key
}

func (i *lruCacheItem[T]) Value() T {
	return i.value
}

func (c *LRUCache[T]) get(key string) (T, bool) {
	if it, ok := c.store.Get(key); ok {
		c.metrics.Hit()
		return it.Value(), true
	}
	c.metrics.Miss()
	var zeroValue T
	return zeroValue, false
}

func (c *LRUCache[T]) put(key string, value T) {
	c.store.Put(&lruCacheItem[T]{key: key, value: value})
}

func (c *LRUCache[T]) delete(key string) {
	c.store.Delete(key)
}