package cache

import "testing"

type TestItem struct{}

func TestLRUCache_New(t *testing.T) {
	cache := newLRUCache[*TestItem]("Name", 10, defaultKeyFUnc)
	if cache == nil {
		t.Fatal("Failed to create an instance.")
	}
}

func TestLRUCache_Name(t *testing.T) {
	cache := newLRUCache[*TestItem]("Name", 10, defaultKeyFUnc)
	if cache.Name() != "Name" {
		t.Errorf("Name() = %s; want Name", cache.Name())
	}
}

func TestLRUCache_Get(t *testing.T) {
	cache := newLRUCache[*TestItem]("Name", 10, defaultKeyFUnc)
	cache.Put([]any{"key"}, &TestItem{})
	if _, ok := cache.Get("key"); !ok {
		t.Errorf("Get() = (_, false); want (_, true)")
	}
}

func TestLRUCache_Put(t *testing.T) {
	cache := newLRUCache[*TestItem]("Name", 10, defaultKeyFUnc)
	cache.Put([]any{"key"}, &TestItem{})
	if _, ok := cache.Get("key"); !ok {
		t.Errorf("Put() = (_, false); want (_, true)")
	}
}

func TestLRUCache_Delete(t *testing.T) {
	cache := newLRUCache[*TestItem]("Name", 10, defaultKeyFUnc)
	cache.Put([]any{"key"}, &TestItem{})
	cache.Delete("key")
	if _, ok := cache.Get("key"); ok {
		t.Errorf("Delete() = (_, true); want (_, false)")
	}
}
