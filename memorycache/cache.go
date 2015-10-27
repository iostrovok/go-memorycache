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

	percentBad int
	ttlClean   time.Duration
	Shards     []*Shard
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
	>>>>>>>>> Config functions START
*/

// PercentBad sets the TTL for each shards.
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

// PercentBad sets the percentage of overload for each shards.
func (mCache *MemoryCache) PercentBad(i int) {
	mCache.RLock()
	defer mCache.RUnlock()

	mCache.percentBad = i

	for _, one := range mCache.Shards {
		mes := NewRequest(TypeSetPercentBad, NewKey(""))
		mes.Data = i
		one.Act(mes)
	}
}

/*
	Config functions FINISH <<<<<<<<<
*/

func (mCache *MemoryCache) Close() {
	mCache.Lock()
	defer mCache.Unlock()

	for _, one := range mCache.Shards {
		one.Stop()
	}
}

// _sendToShard send one message to the shard
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

// GetTag returns data by tag
func (mCache *MemoryCache) GetTag(tag string) ([]interface{}, bool) {

	out := []interface{}{}

	if tag == "" {
		return out, false
	}

	chRes := make(chan *Res, 1)

	count := 0
	mCache.RLock()
	for _, one := range mCache.Shards {
		count++
		mes := NewRequest(TypeGetTag, NewKey(""))
		mes.Tags = []string{tag}
		mes.ResultChan = chRes
		one.Act(mes)
	}
	mCache.RUnlock()

	for count > 0 {
		count--

		res, ok := <-chRes
		if !ok {
			continue
		}

		d, ok := res.Data.([]interface{})
		if ok {
			out = append(out, d)
		}
	}

	if len(out) > 0 {
		return out, true
	}
	return out, false
}

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
