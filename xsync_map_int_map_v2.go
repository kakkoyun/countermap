package countermap

import (
	"sync/atomic"

	xsync "github.com/puzpuzpuz/xsync/v2"
)

func NewXSyncV2MapIntMap() *XSyncV2MapIntMap {
	return &XSyncV2MapIntMap{counts: xsync.NewMap()}
}

type XSyncV2MapIntMap struct {
	counts *xsync.Map
}

func (cm *XSyncV2MapIntMap) Inc(key string) {
	val, ok := cm.counts.Load(key)
	if !ok {
		val, _ = cm.counts.LoadOrStore(key, &atomic.Int64{})
	}
	val.(*atomic.Int64).Add(1)
}

func (cm *XSyncV2MapIntMap) GetAndReset() map[string]int64 {
	ret := map[string]int64{}
	cm.counts.Range(func(key string, val any) bool {
		ret[key] = val.(*atomic.Int64).Load()
		cm.counts.Delete(key)
		return true
	})
	return ret
}
