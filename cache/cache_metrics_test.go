package cache

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_CacheMetrics(t *testing.T) {
	t.Run("should initialize with zero counts", func(t *testing.T) {
		cm := &cacheMetrics{}

		assert.Equal(t, uint64(0), cm.HitCount())
		assert.Equal(t, uint64(0), cm.MissCount())
		assert.Equal(t, uint64(0), cm.Total())
		assert.Equal(t, float64(0), cm.HitRate())
	})

	t.Run("should increment hit count", func(t *testing.T) {
		cm := &cacheMetrics{}
		cm.Hit()
		cm.Hit()

		assert.Equal(t, uint64(2), cm.HitCount())
		assert.Equal(t, uint64(0), cm.MissCount())
		assert.Equal(t, uint64(2), cm.Total())
		assert.Equal(t, float64(100), cm.HitRate())
	})

	t.Run("should increment miss count", func(t *testing.T) {
		cm := &cacheMetrics{}
		cm.Miss()
		cm.Miss()
		cm.Miss()

		assert.Equal(t, uint64(0), cm.HitCount())
		assert.Equal(t, uint64(3), cm.MissCount())
		assert.Equal(t, uint64(3), cm.Total())
		assert.Equal(t, float64(0), cm.HitRate())
	})

	t.Run("should calculate hit rate correctly", func(t *testing.T) {
		cm := &cacheMetrics{}
		cm.Hit()
		cm.Miss() // 1 hit, 1 miss

		assert.Equal(t, float64(50), cm.HitRate())
	})

	t.Run("should reset all counters", func(t *testing.T) {
		cm := &cacheMetrics{}
		cm.Hit()
		cm.Miss()
		cm.Reset()

		assert.Equal(t, uint64(0), cm.HitCount())
		assert.Equal(t, uint64(0), cm.MissCount())
		assert.Equal(t, float64(0), cm.HitRate())
	})
}
