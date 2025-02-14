package countermap

import "sync"

func NewMutexCounterMap() *MutexCounterMap {
	return &MutexCounterMap{counts: map[string]int64{}}
}

type MutexCounterMap struct {
	lock   sync.Mutex
	counts map[string]int64
}

func (cm *MutexCounterMap) Inc(key string) {
	cm.lock.Lock()
	cm.counts[key]++
	cm.lock.Unlock()
}

func (cm *MutexCounterMap) GetAndReset() map[string]int64 {
	cm.lock.Lock()
	ret := cm.counts
	cm.counts = map[string]int64{}
	cm.lock.Unlock()
	return ret
}
