package cache

type Cache[T any] interface {
	Name() string
	Get(keys ...any) (T, bool)
	Put(keys []any, value T)
	Delete(keys ...any)

	Stat() *CacheStat
}

type Item[T any] interface {
	Key() string
	Value() T
}

type CacheStat struct {
	Name        string
	MaxEntries  int
	CurrentSize int
	HitCount    uint64
	MissCount   uint64
	HitRate     float64
}
