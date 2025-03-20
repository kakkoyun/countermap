package countermap

import (
	"sync/atomic"

	"github.com/puzpuzpuz/xsync/v3"
)

type XSyncMapIntMap struct {
	counts *xsync.MapOf[string, *atomic.Int64]
}

func NewXSyncMapIntMap() *XSyncMapIntMap {
	return &XSyncMapIntMap{
		counts: xsync.NewMapOf[string, *atomic.Int64](),
	}
}

func (cm *XSyncMapIntMap) Inc(key string) {
	// NOTICE: LoadOrCompute locks the whole map for the duration of the function,
	// it is not suitable for high-contention scenarios.
	val, ok := cm.counts.Load(key)
	if !ok {
		val, _ = cm.counts.LoadOrStore(key, &atomic.Int64{})
	}
	val.Add(1)
}

func (cm *XSyncMapIntMap) GetAndReset() map[string]int64 {
	ret := map[string]int64{}
	cm.counts.Range(func(key string, value *atomic.Int64) bool {
		val, loaded := cm.counts.LoadAndDelete(key)
		if loaded {
			ret[key] = val.Load()
		}
		return true
	})
	return ret
}
