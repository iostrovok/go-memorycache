package memorycache

import (
	"sync"
	"time"
)

// MemoryCache storage memory
type MemoryCache struct {
	sync.RWMutex

	maxShards          int
	maxEntries         int
	maxEntriesPerShard int

	ttlClean time.Duration
	Shards   []*Shard
}

// NewMemoryCache returns new initializated instance of MemoryCache
func New(shards int, totalLimit int) *MemoryCache {
	storage := &MemoryCache{}
	storage.maxShards = shards
	storage.maxEntries = totalLimit
	storage.Shards = make([]*Shard, storage.maxShards)
	storage.maxEntriesPerShard = 1 + totalLimit/shards

	storage.ttlClean = TTLCleanDefault

	for i := 0; i < shards; i++ {
		storage.Shards[i] = NewShard(storage.maxEntriesPerShard)
	}

	return storage
}

/*
	>>>>>>>>> Config function
*/
func (storage *MemoryCache) TTLClean(t time.Duration) {
	storage.RLock()
	defer storage.RUnlock()

	storage.ttlClean = t

	for _, one := range storage.Shards {
		mes := NewRequest(TypeSetTTL, NewKey(""))
		mes.TTL = t
		one.Act(mes)
	}
}

/*
	Config function <<<<<<<<<
*/

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
func (storage *MemoryCache) Get(k string) (interface{}, bool) {

	mes := NewRequest(TypeGet, NewKey(k))

	storage._sendToShard(mes)

	out, ok := <-mes.ResultChan

	if ok && out.Entry.Valid() {
		return out.Entry.Data, true
	}

	return nil, false
}

// // Get returns data by key
// func (storage *MemoryCache) GetTags(tags ...string) ([]interface{}, bool) {

// 	if len(tags) == 0 {
// 		return nil, false
// 	}

// 	mes := NewRequest(TypeGetTag, NewKey(""))

// 	storage._sendToShard(mes)

// 	out, ok := <-mes.ResultChan
// 	if !ok {
// 		return nil, false
// 	}

// 	data, ok := out.Entry.Data.([]interface{})
// 	if !ok {
// 		return nil, false
// 	}

// 	return data, true
// }

// Put puts new data in storage _WITHOUT_ TTL
func (storage *MemoryCache) Put(data interface{}, k string, tags ...string) {
	storage.PutTTL(data, k, time.Duration(0), tags...)
}

// Put puts new data in storage
func (storage *MemoryCache) PutTTL(data interface{}, k string, TTL time.Duration, tags ...string) {

	mes := NewRequest(TypePut, NewKey(k))
	mes.Data = data
	mes.Tags = tags
	mes.TTL = TTL

	storage._sendToShard(mes)
}

// Remove remove data in cache by key
func (storage *MemoryCache) Remove(k string) {
	mes := NewRequest(TypeRemove, NewKey(k))
	storage._sendToShard(mes)
}

func (storage *MemoryCache) RemoveTag(tag string) {
	storage.RLock()
	defer storage.RUnlock()

	for _, one := range storage.Shards {
		mes := NewRequest(TypeRemTag, NewKey(""))
		one.Act(mes)
	}
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
