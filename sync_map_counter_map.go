package countermap

import (
	"sync"
)

// syncMapStorage adapts sync.Map to CounterStorage
type syncMapStorage[K comparable, V counter] struct {
	m sync.Map
}

func newSyncMapStorage[K comparable, V counter]() counterStorage[K, V] {
	return &syncMapStorage[K, V]{}
}

func (s *syncMapStorage[K, V]) Load(key K) (V, bool) {
	var zero V
	value, ok := s.m.Load(key)
	if !ok {
		return zero, false
	}
	return value.(V), true
}

func (s *syncMapStorage[K, V]) Store(key K, value V) {
	s.m.Store(key, value)
}

func (s *syncMapStorage[K, V]) LoadOrStore(key K, value V) (V, bool) {
	actual, loaded := s.m.LoadOrStore(key, value)
	if !loaded {
		return value, false
	}
	return actual.(V), true
}

func (s *syncMapStorage[K, V]) LoadAndDelete(key K) (V, bool) {
	actual, loaded := s.m.LoadAndDelete(key)
	return actual.(V), loaded
}

func (s *syncMapStorage[K, V]) Range(f func(key K, value V) bool) {
	s.m.Range(func(key, value any) bool {
		return f(key.(K), value.(V))
	})
}

func (s *syncMapStorage[K, V]) Clear() {
	s.m.Clear()
}

// SyncMapCounterMap implements CounterMap using sync.Map.
type SyncMapCounterMap struct {
	*counterMap[string, counter]
}

// NewSyncMapCounterMap creates a new counter map backed by sync.Map.
func NewSyncMapCounterMap() *SyncMapCounterMap {
	// Create function for the counterInt64 type
	newCounter := func() counter {
		return &counterInt64{}
	}

	cm := &SyncMapCounterMap{
		counterMap: newCounterMap[string, counter](
			func() counterStorage[string, counter] {
				return newSyncMapStorage[string, counter]()
			},
			newCounter,
		),
	}
	return cm
}
