package cache

import (
	"context"
	"fmt"
	"reflect"
	"time"
)

const (
	LRU CacheType = iota
	EXPIRE_TYPE_LAZY_DELETION
	EXPIRE_TYPE_ACTIVE_EXPIRATION
)

var (
	defaultCacheSize       = 1000
	defaultKeyFUnc         = defaultKeyFunc
	defaultTTL             = time.Minute * 10
	defaultExpireExtension = false
	defaultSamplingDelay   = time.Millisecond * 100
	defaultSamplingRatio   = 25
	defaultSamplingSize    = 100
)

type CacheType int

type Option struct {
	Type            CacheType
	Name            string
	CacheSize       int
	KeyFunc         KeyFunc
	TTL             time.Duration
	ExpireExtension bool
	Sampling        struct {
		BackgroundCtx context.Context
		Delay         time.Duration
		Ratio         int
		Size          int
	}
}

func New[T any](opts ...Option) Cache[T] {
	tp, name, cacheSize, keyFunc, ttl, expireExtension, backgroundCtx, sampleDelay, sampleRatio, sampleSize := split[T](opts...)
	switch tp {
	case LRU:
		return NewLRUCache[T](name, cacheSize, keyFunc)
	case EXPIRE_TYPE_LAZY_DELETION:
		return NewExpireCache[T](name, cacheSize, keyFunc, ttl, expireExtension)
	case EXPIRE_TYPE_ACTIVE_EXPIRATION:
		return NewActiveExpireCache[T](backgroundCtx, name, cacheSize, keyFunc, ttl, expireExtension, sampleDelay, sampleRatio, sampleSize)
	default:
		return NewLRUCache[T](name, cacheSize, keyFunc)
	}
}

func split[T any](opts ...Option) (cacheType CacheType, name string, cacheSize int, keyFunc KeyFunc, ttl time.Duration, expireExtension bool, backgroundCtx context.Context, samplingDelay time.Duration, samplingRatio int, samplingSize int) {
	var zeroValue T
	name = fmt.Sprintf("Cache-%v", reflect.TypeOf(zeroValue).Name())
	cacheSize = defaultCacheSize
	keyFunc = defaultKeyFUnc
	ttl = defaultTTL
	expireExtension = defaultExpireExtension
	samplingDelay = defaultSamplingDelay
	samplingRatio = defaultSamplingRatio
	samplingSize = defaultSamplingSize

	if len(opts) > 0 {
		cacheType = opts[0].Type
		if opts[0].Name != "" {
			name = opts[0].Name
		}
		if opts[0].CacheSize > 0 {
			cacheSize = opts[0].CacheSize
		}
		if opts[0].KeyFunc != nil {
			keyFunc = opts[0].KeyFunc
		}
		if opts[0].TTL > 0 {
			ttl = opts[0].TTL
		}
		expireExtension = opts[0].ExpireExtension
		backgroundCtx = opts[0].Sampling.BackgroundCtx
		if opts[0].Sampling.Delay > 0 {
			samplingDelay = opts[0].Sampling.Delay
		}
		if opts[0].Sampling.Ratio > 0 {
			samplingRatio = opts[0].Sampling.Ratio
		}
		if opts[0].Sampling.Size > 0 {
			samplingSize = opts[0].Sampling.Size
		}
	}
	return
}
