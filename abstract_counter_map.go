package countermap

import (
	"sync"
	"sync/atomic"
)

// counterStorage defines the interface that any underlying map implementation must provide.
type counterStorage[K comparable, V counter] interface {
	// Load retrieves a counter from the map.
	Load(key K) (V, bool)
	// Store adds a counter to the map.
	Store(key K, value V)
	// LoadOrStore atomically loads or initializes a counter.
	LoadOrStore(key K, value V) (actual V, loaded bool)
	// LoadAndDelete atomically loads and deletes a counter.
	LoadAndDelete(key K) (value V, loaded bool)
	// Range iterates over all entries in the map.
	Range(f func(key K, value V) bool)
	// Clear removes all entries from the map.
	Clear()
}

// state encapsulates a map and its associated in-flight operation tracking.
type state[K comparable, V counter] struct {
	counters    counterStorage[K, V]
	inFlightOps *atomic.Int64
	cond        *sync.Cond
}

// newState creates a new state with the given counter storage.
func newState[K comparable, V counter](storage counterStorage[K, V]) *state[K, V] {
	return &state[K, V]{
		counters:    storage,
		inFlightOps: &atomic.Int64{},
		cond:        sync.NewCond(&sync.Mutex{}),
	}
}

// counterMap provides a generic implementation of a counter map
type counterMap[K comparable, V counter] struct {
	// Pointer to the active state - allows atomic swapping.
	activeState atomic.Pointer[state[K, V]]
	// Prevents concurrent GetAndReset operations.
	snapshotLock *sync.RWMutex

	// Function to create a counter value.
	newCounterFn func() V
	// Function to create new storage.
	newStorageFn func() counterStorage[K, V]
}

// newCounterMap creates a new abstract counter map with the specified functions.
func newCounterMap[K comparable, V counter](
	newStorageFn func() counterStorage[K, counter],
	newCounterFn func() counter,
) *counterMap[K, counter] {
	cm := &counterMap[K, counter]{
		snapshotLock: &sync.RWMutex{},
		newCounterFn: newCounterFn,
		newStorageFn: newStorageFn,
	}
	// Initialize with an empty state.
	storage := newStorageFn()
	cm.activeState.Store(newState(storage))
	return cm
}

// Inc increments the counter for the given key.
func (cm *counterMap[K, V]) Inc(key K) {
	// Get the current active state.
	state := cm.activeState.Load()

	// Mark operation as in-flight on this state.
	state.inFlightOps.Add(1)
	defer func() {
		// If this was the last in-flight operation, signal waiters
		if state.inFlightOps.Add(-1) == 0 {
			state.cond.L.Lock()
			state.cond.Signal()
			state.cond.L.Unlock()
		}
	}()

	// Try to load existing counter.
	counter, loaded := state.counters.Load(key)
	if !loaded {
		// Create a new counter.
		counter = cm.newCounterFn()
		actual, loaded := state.counters.LoadOrStore(key, counter)
		if loaded {
			// Another goroutine beat us to it, use their counter.
			counter = actual
		}
	}

	// Increment the counter.
	counter.Add(1)
}

// GetAndReset returns the current counts and resets all counters.
func (cm *counterMap[K, V]) GetAndReset() map[K]int64 {
	// Ensure only one GetAndReset operation runs at a time.
	cm.snapshotLock.Lock()
	defer cm.snapshotLock.Unlock()

	// Create a new empty state for new increments.
	newState := newState(cm.newStorageFn())

	// Atomically swap the states.
	oldState := cm.activeState.Swap(newState)

	// Wait for all in-flight operations on the old map to complete.
	oldState.cond.L.Lock()
	for oldState.inFlightOps.Load() > 0 {
		oldState.cond.Wait() // Releases the lock while waiting.
	}
	oldState.cond.L.Unlock()

	// Process the old state - no more in-flight operations on it.
	result := make(map[K]int64)
	oldState.counters.Range(func(key K, _ V) bool {
		value, loaded := oldState.counters.LoadAndDelete(key)
		if loaded {
			result[key] = value.Value()
		}
		return true
	})

	return result
}
