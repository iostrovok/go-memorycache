package memorycache

import (
	"sync"
)

// MemoryCache storage memory
type MemoryCache struct {
	sync.RWMutex

	maxShards          int
	maxEntries         int
	maxEntriesPerShard int

	Shards []*Shard
}

// NewMemoryCache returns new initializated instance of MemoryCache
func New(shards int, totalLimit int) *MemoryCache {
	storage := &MemoryCache{}
	storage.maxShards = shards
	storage.maxEntries = totalLimit

	storage.Shards = make([]*Shard, storage.maxShards)

	storage.maxEntriesPerShard = 1 + totalLimit/shards

	for i := 0; i < shards; i++ {
		storage.Shards[i] = NewShard(storage.maxEntriesPerShard)
	}

	return storage
}

func (storage *MemoryCache) Close() {
	storage.Lock()
	defer storage.Unlock()

	for _, one := range storage.Shards {
		one.Stop()
	}
}

// Get returns data by key
func (storage *MemoryCache) _sendToShard(mes *Request) {

	keyShard := mes.Key.ShardID(storage.maxShards)

	storage.RLock()
	storage.Shards[keyShard].Act(mes)
	storage.RUnlock()

}

// Get returns data by key
func (storage *MemoryCache) Get(k string) (data interface{}, ok bool) {

	mes := NewRequest(TypeGet, NewKey(k))

	storage._sendToShard(mes)

	out, ok := <-mes.ResultChan
	if !ok {
		return nil, false
	}

	return out.Entry.Data, out.Ok
}

// Put puts new data in storage
func (storage *MemoryCache) Put(data interface{}, k string, tags ...string) {

	mes := NewRequest(TypePut, NewKey(k))
	mes.Data = data

	storage._sendToShard(mes)
}

// Remove remove data in cache by key
func (storage *MemoryCache) Remove(k string) {
	mes := NewRequest(TypeRemove, NewKey(k))
	storage._sendToShard(mes)
}

// Flush removes all entries from cache and returns number of flushed entries
func (storage *MemoryCache) Flush() {

	storage.RLock()
	defer storage.RUnlock()

	for _, one := range storage.Shards {
		mes := NewRequest(TypeFlush, NewKey(""))
		one.Act(mes)
	}

}
