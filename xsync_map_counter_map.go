package countermap

import (
	xsync "github.com/puzpuzpuz/xsync/v3"
)

func NewXSyncMapCounterMap() *XSyncMapCounterMap {
	return &XSyncMapCounterMap{counts: xsync.NewMapOf[string, *xsync.Counter]()}
}

type XSyncMapCounterMap struct {
	counts *xsync.MapOf[string, *xsync.Counter]
}

func (cm *XSyncMapCounterMap) Inc(key string) {
	// NOTICE: LoadOrCompute locks the whole map for the duration of the function,
	// it is not suitable for high-contention scenarios.
	val, ok := cm.counts.Load(key)
	if !ok {
		val, _ = cm.counts.LoadOrStore(key, xsync.NewCounter())
	}
	val.Inc()
}

func (cm *XSyncMapCounterMap) GetAndReset() map[string]int64 {
	ret := map[string]int64{}
	cm.counts.Range(func(key string, counter *xsync.Counter) bool {
		val, loaded := cm.counts.LoadAndDelete(key)
		if loaded {
			ret[key] = val.Value()
		}
		return true
	})
	return ret
}
