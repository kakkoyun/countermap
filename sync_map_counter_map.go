package countermap

import (
	"sync"
	"sync/atomic"
)

func NewSyncMapCounterMap() *SyncMapCounterMap {
	return &SyncMapCounterMap{}
}

type SyncMapCounterMap struct {
	counts sync.Map
}

func (cm *SyncMapCounterMap) Inc(key string) {
	val, ok := cm.counts.Load(key)
	if !ok {
		val, _ = cm.counts.LoadOrStore(key, &atomic.Int64{})
	}
	val.(*atomic.Int64).Add(1)
}

func (cm *SyncMapCounterMap) GetAndReset() map[string]int64 {
	ret := make(map[string]int64)
	cm.counts.Range(func(key, value any) bool {
		k := key.(string)
		if val, ok := cm.counts.Swap(k, &atomic.Int64{}); ok {
			if v := val.(*atomic.Int64).Load(); v > 0 {
				ret[k] = v
			}
		}
		return true
	})
	return ret
}
