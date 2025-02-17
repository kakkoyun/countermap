package countermap

import (
	xsync "github.com/puzpuzpuz/xsync/v2"
)

func NewXSyncV2MapCounterMap() *XSyncV2MapCounterMap {
	return &XSyncV2MapCounterMap{counts: xsync.NewMap()}
}

type XSyncV2MapCounterMap struct {
	counts *xsync.Map
}

func (cm *XSyncV2MapCounterMap) Inc(key string) {
	val, _ := cm.counts.LoadOrStore(key, xsync.NewCounter())
	val.(*xsync.Counter).Add(1)
}

func (cm *XSyncV2MapCounterMap) GetAndReset() map[string]int64 {
	ret := map[string]int64{}
	cm.counts.Range(func(key string, val any) bool {
		ret[key] = val.(*xsync.Counter).Value()
		cm.counts.Delete(key)
		return true
	})
	return ret
}
