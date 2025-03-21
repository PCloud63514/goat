package cache

import (
	"container/list"
	"sync"
)

type cacheStore[T any, IT storeItem[T]] struct {
	mu       sync.RWMutex
	buff     map[string]*list.Element
	ll       *list.List
	capacity int
}

type storeItem[T any] interface {
	Key() string
	Value() T
}

func newCacheStore[T any, IT storeItem[T]](capacity int) *cacheStore[T, IT] {
	return &cacheStore[T, IT]{
		mu:       sync.RWMutex{},
		buff:     make(map[string]*list.Element),
		ll:       list.New(),
		capacity: capacity,
	}
}

func (store *cacheStore[T, IT]) Get(key string) (IT, bool) {
	store.mu.RLock()
	defer store.mu.RUnlock()
	var zeroValue IT
	if elem, found := store.buff[key]; found {
		store.ll.MoveToFront(elem)
		return elem.Value.(IT), true
	}
	return zeroValue, false
}

func (store *cacheStore[T, IT]) Put(item IT) {
	store.mu.Lock()
	defer store.mu.Unlock()

	if elem, found := store.buff[item.Key()]; found {
		store.ll.MoveToFront(elem)
		elem.Value = item
		return
	}

	elem := store.ll.PushFront(item)
	store.buff[item.Key()] = elem

	if store.ll.Len() > store.capacity {
		last := store.ll.Back()
		if last != nil {
			store.ll.Remove(last)
			delete(store.buff, last.Value.(IT).Key())
		}
	}
}

func (store *cacheStore[T, IT]) Delete(key string) {
	store.mu.Lock()
	defer store.mu.Unlock()
	if elem, found := store.buff[key]; found {
		store.ll.Remove(elem)
		delete(store.buff, key)
	}
}

func (store *cacheStore[T, IT]) Contains(key string) bool {
	store.mu.RLock()
	defer store.mu.RUnlock()
	_, found := store.buff[key]
	return found
}

func (store *cacheStore[T, IT]) Keys() []string {
	store.mu.RLock()
	defer store.mu.RUnlock()
	keys := make([]string, 0, len(store.buff))
	for key := range store.buff {
		keys = append(keys, key)
	}
	return keys
}

func (store *cacheStore[T, IT]) SampleKeys(samplingSize int) []string {
	store.mu.RLock()
	defer store.mu.RUnlock()

	sampledKeys := make([]string, 0, samplingSize)
	sampled := 0

	for elem := store.ll.Front(); elem != nil && sampled < samplingSize; elem = elem.Next() {
		item := elem.Value.(*activeExpireCacheItem[T])
		sampledKeys = append(sampledKeys, item.key)
		sampled++
	}
	return sampledKeys
}

func (store *cacheStore[T, IT]) Clear() {
	store.mu.Lock()
	defer store.mu.Unlock()
	store.buff = make(map[string]*list.Element)
	store.ll.Init()
}

func (store *cacheStore[T, IT]) Size() int {
	store.mu.RLock()
	defer store.mu.RUnlock()
	return store.ll.Len()
}

func (store *cacheStore[T, IT]) Capacity() int {
	return store.capacity
}
