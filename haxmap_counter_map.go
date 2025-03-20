package countermap

import (
	"sync/atomic"

	"github.com/alphadose/haxmap"
)

type HaxMapCounterMap struct {
	counts *haxmap.Map[string, *atomic.Int64]
}

func NewHaxMapCounterMap() CounterMap {
	return &HaxMapCounterMap{
		counts: haxmap.New[string, *atomic.Int64](),
	}
}

func (cm *HaxMapCounterMap) Inc(key string) {
	counter, loaded := cm.counts.GetOrCompute(key, func() *atomic.Int64 {
		var i atomic.Int64
		i.Store(1)
		return &i
	})
	if loaded {
		counter.Add(1)
	}
}

func (cm *HaxMapCounterMap) GetAndReset() map[string]int64 {
	result := make(map[string]int64)
	cm.counts.ForEach(func(key string, counter *atomic.Int64) bool {
		count, loaded := cm.counts.GetAndDel(key)
		if loaded {
			result[key] = count.Load()
		}
		return true
	})
	return result
}
