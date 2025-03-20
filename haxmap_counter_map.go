package countermap

import (
	"github.com/alphadose/haxmap"
)

// haxMapStorage adapts haxmap.Map to CounterStorage for string keys
type haxMapStorage[V counter] struct {
	m *haxmap.Map[string, V]
}

func newHaxMapStorage[V counter]() counterStorage[string, V] {
	return &haxMapStorage[V]{
		m: haxmap.New[string, V](),
	}
}

func (s *haxMapStorage[V]) Load(key string) (V, bool) {
	return s.m.Get(key)
}

func (s *haxMapStorage[V]) Store(key string, value V) {
	s.m.Set(key, value)
}

func (s *haxMapStorage[V]) LoadOrStore(key string, value V) (V, bool) {
	return s.m.GetOrSet(key, value)
}

func (s *haxMapStorage[V]) LoadAndDelete(key string) (V, bool) {
	return s.m.GetAndDel(key)
}

func (s *haxMapStorage[V]) Range(f func(key string, value V) bool) {
	s.m.ForEach(func(key string, value V) bool {
		return f(key, value)
	})
}

func (s *haxMapStorage[V]) Clear() {
	s.m.Clear()
}

// HaxMapCounterMap implements CounterMap using haxmap.
type HaxMapCounterMap struct {
	*counterMap[string, counter]
}

// NewHaxMapCounterMap creates a counter map backed by haxmap.
func NewHaxMapCounterMap() CounterMap {
	// Create function for the counterInt64 type.
	newCounter := func() counter {
		return &counterInt64{}
	}

	return &HaxMapCounterMap{
		counterMap: newCounterMap[string, counter](
			func() counterStorage[string, counter] {
				return newHaxMapStorage[counter]()
			},
			newCounter,
		),
	}
}
