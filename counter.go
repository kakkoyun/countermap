package countermap

import "sync/atomic"

// counter is an interface for a counter value.
type counter interface {
	Add(int64)
	Value() int64
}

// counterInt64 adapts atomic.Int64 to the counter interface
type counterInt64 struct {
	value atomic.Int64
}

func (c *counterInt64) Add(val int64) {
	c.value.Add(val)
}

func (c *counterInt64) Value() int64 {
	return c.value.Load()
}
