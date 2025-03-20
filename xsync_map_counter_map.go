package countermap

import "github.com/puzpuzpuz/xsync/v3"

// XSyncMapCounterMap implements CounterMap using xsync.MapOf
type XSyncMapCounterMap struct {
	*counterMap[string, counter]
}

// NewXSyncMapCounterMap creates a new counter map backed by xsync.MapOf of xsync.Counter.
func NewXSyncMapCounterMap() CounterMap {
	// Create function for the counterInt64 type.
	newCounter := func() counter {
		return xsync.NewCounter()
	}

	return &XSyncMapCounterMap{
		counterMap: newCounterMap[string, counter](
			func() counterStorage[string, counter] {
				return NewXSyncMapStorage[string, counter]()
			},
			newCounter,
		),
	}
}
