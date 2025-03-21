package cache

import "sync/atomic"

type cacheMetrics struct {
	hitCount  uint64
	missCount uint64
}

func (cm *cacheMetrics) Hit() {
	atomic.AddUint64(&cm.hitCount, 1)
}

func (cm *cacheMetrics) Miss() {
	atomic.AddUint64(&cm.missCount, 1)
}

func (cm *cacheMetrics) HitCount() uint64 {
	return atomic.LoadUint64(&cm.hitCount)
}

func (cm *cacheMetrics) MissCount() uint64 {
	return atomic.LoadUint64(&cm.missCount)
}

func (cm *cacheMetrics) Total() uint64 {
	return cm.HitCount() + cm.MissCount()
}

func (cm *cacheMetrics) HitRate() float64 {
	total := cm.Total()
	if total == 0 {
		return 0
	}
	return (float64(cm.HitCount()) / float64(total)) * 100
}

func (cm *cacheMetrics) Reset() {
	atomic.StoreUint64(&cm.hitCount, 0)
	atomic.StoreUint64(&cm.missCount, 0)
}
