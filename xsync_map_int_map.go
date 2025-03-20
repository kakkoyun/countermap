package countermap

import (
	xsync "github.com/puzpuzpuz/xsync/v3"
)

// xsyncMapStorage adapts xsync.MapOf to CounterStorage.
type xsyncMapStorage[K comparable, V counter] struct {
	m *xsync.MapOf[K, V]
}

func NewXSyncMapStorage[K comparable, V counter]() counterStorage[K, V] {
	return &xsyncMapStorage[K, V]{
		m: xsync.NewMapOf[K, V](),
	}
}

func (s *xsyncMapStorage[K, V]) Load(key K) (V, bool) {
	return s.m.Load(key)
}

func (s *xsyncMapStorage[K, V]) Store(key K, value V) {
	s.m.Store(key, value)
}

func (s *xsyncMapStorage[K, V]) LoadOrStore(key K, value V) (V, bool) {
	return s.m.LoadOrStore(key, value)
}

func (s *xsyncMapStorage[K, V]) LoadAndDelete(key K) (V, bool) {
	return s.m.LoadAndDelete(key)
}

func (s *xsyncMapStorage[K, V]) Range(f func(key K, value V) bool) {
	s.m.Range(f)
}

func (s *xsyncMapStorage[K, V]) Clear() {
	s.m.Clear()
}

// XSyncMapIntMap implements CounterMap using xsync.MapOf of ints.
type XSyncMapIntMap struct {
	*counterMap[string, counter]
}

// NewXSyncMapIntMap creates a new counter map backed by xsync.MapOf of ints.
// This implementation is identical to XSyncMapCounterMap but preserved for backward compatibility
func NewXSyncMapIntMap() CounterMap {
	// Create function for the counterInt64 type
	newCounter := func() counter {
		return &counterInt64{}
	}

	return &XSyncMapIntMap{
		counterMap: newCounterMap[string, counter](
			func() counterStorage[string, counter] {
				return NewXSyncMapStorage[string, counter]()
			},
			newCounter,
		),
	}
}
