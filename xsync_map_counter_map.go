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
	counter, _ := cm.counts.LoadOrStore(key, xsync.NewCounter())
	counter.Add(1)
}

func (cm *XSyncMapCounterMap) GetAndReset() map[string]int64 {
	ret := map[string]int64{}
	cm.counts.Range(func(key string, counter *xsync.Counter) bool {
		ret[key] = counter.Value()
		cm.counts.Delete(key)
		return true
	})
	return ret
}
