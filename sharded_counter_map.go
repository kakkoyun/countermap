package countermap

import (
	"hash/fnv"
	"sync"
)

const numberOfShards = 256

type ShardedCounterMap struct {
	shards [numberOfShards]shard
}

type shard struct {
	lock   sync.Mutex
	counts map[string]int64
}

func NewShardedCounterMap() *ShardedCounterMap {
	cm := &ShardedCounterMap{}
	for i := range cm.shards {
		cm.shards[i].counts = make(map[string]int64)
	}
	return cm
}

func (c *ShardedCounterMap) shard(key string) *shard {
	h := fnv.New64a()
	h.Write([]byte(key))
	return &c.shards[h.Sum64()%numberOfShards]
}

func (c *ShardedCounterMap) Inc(key string) {
	shard := c.shard(key)
	shard.lock.Lock()
	shard.counts[key]++
	shard.lock.Unlock()
}

func (c *ShardedCounterMap) GetAndReset() map[string]int64 {
	ret := make(map[string]int64)

	// Lock all shards in order to prevent deadlocks.
	for i := range c.shards {
		c.shards[i].lock.Lock()
	}

	// Collect all non-zero values.
	for i := range c.shards {
		for k, v := range c.shards[i].counts {
			if v > 0 {
				ret[k] = v
			}
		}
		c.shards[i].counts = make(map[string]int64)
	}

	// Unlock all shards.
	for i := range c.shards {
		c.shards[i].lock.Unlock()
	}

	return ret
}
