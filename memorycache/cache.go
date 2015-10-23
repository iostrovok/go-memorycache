package memorycache

import (
	"sync"
	"time"
)

// MemoryCache mCache memory
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
	mCache := &MemoryCache{}
	mCache.maxShards = shards
	mCache.maxEntries = totalLimit
	mCache.Shards = make([]*Shard, mCache.maxShards)
	mCache.maxEntriesPerShard = 1 + totalLimit/shards

	mCache.ttlClean = TTLCleanDefault

	for i := 0; i < shards; i++ {
		mCache.Shards[i] = NewShard(mCache.maxEntriesPerShard)
	}

	return mCache
}

/*
	>>>>>>>>> Config function
*/
func (mCache *MemoryCache) TTLClean(t time.Duration) {
	mCache.RLock()
	defer mCache.RUnlock()

	mCache.ttlClean = t

	for _, one := range mCache.Shards {
		mes := NewRequest(TypeSetTTL, NewKey(""))
		mes.TTL = t
		one.Act(mes)
	}
}

/*
	Config function <<<<<<<<<
*/

func (mCache *MemoryCache) Close() {
	mCache.Lock()
	defer mCache.Unlock()

	for _, one := range mCache.Shards {
		one.Stop()
	}
}

// Get returns data by key
func (mCache *MemoryCache) _sendToShard(mes *Request) {

	keyShard := mes.Key.ShardID(mCache.maxShards)

	mCache.RLock()
	mCache.Shards[keyShard].Act(mes)
	mCache.RUnlock()

}

// Get returns data by key
func (mCache *MemoryCache) Get(k string) (interface{}, bool) {

	mes := NewRequest(TypeGet, NewKey(k))

	mCache._sendToShard(mes)

	out, ok := <-mes.ResultChan

	if ok {
		//return out.Entry.Data, true
		return out.Data, true
	}

	return nil, false
}

// // Get returns data by key
// func (mCache *MemoryCache) GetTags(tags ...string) ([]interface{}, bool) {

// 	if len(tags) == 0 {
// 		return nil, false
// 	}

// 	mes := NewRequest(TypeGetTag, NewKey(""))

// 	mCache._sendToShard(mes)

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

// Put puts new data in mCache _WITHOUT_ TTL
func (mCache *MemoryCache) Put(data interface{}, k string, tags ...string) {
	mCache.PutTTL(data, k, time.Duration(0), tags...)
}

// Put puts new data in mCache
func (mCache *MemoryCache) PutTTL(data interface{}, k string, TTL time.Duration, tags ...string) {

	mes := NewRequest(TypePut, NewKey(k))
	mes.Data = data
	mes.Tags = tags
	mes.TTL = TTL

	mCache._sendToShard(mes)
}

// Remove remove data in cache by key
func (mCache *MemoryCache) Remove(k string) {
	mes := NewRequest(TypeRemove, NewKey(k))
	mCache._sendToShard(mes)
}

func (mCache *MemoryCache) RemoveTag(tag string) {
	mCache.RLock()
	defer mCache.RUnlock()

	for _, one := range mCache.Shards {
		mes := NewRequest(TypeRemTag, NewKey(""))
		one.Act(mes)
	}
}

// Flush removes all entries from cache and returns number of flushed entries
func (mCache *MemoryCache) Flush() {

	mCache.RLock()
	defer mCache.RUnlock()

	for _, one := range mCache.Shards {
		mes := NewRequest(TypeFlush, NewKey(""))
		one.Act(mes)
	}

}
