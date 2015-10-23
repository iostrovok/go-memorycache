package memorycache

import (
	"sync"
	"time"
)

const (
	TTLCleanDefault = 1 * time.Minute
)

// MemoryCache mCache memory
type Shard struct {
	sync.RWMutex

	Len     int
	entries map[string]*Entry
	inChan  chan *Request

	oldest []string
	ticker *time.Ticker
}

func NewShard(l int) *Shard {
	s := &Shard{
		Len:     l,
		entries: map[string]*Entry{},
		inChan:  make(chan *Request, 100),
		ticker:  time.NewTicker(TTLCleanDefault),
	}

	go s._work()

	return s
}

func (s *Shard) Act(mes *Request) {
	s.inChan <- mes
}

func (s *Shard) Stop() {
	close(s.inChan)
}

func (s *Shard) _work() {

	for {
		select {
		case <-s.ticker.C:
			s.trim()
		case mes, ok := <-s.inChan:
			if !ok {
				return
			}
			s._make(mes)
		}
	}
}

func (s *Shard) _make(mes *Request) {

	switch mes.Action {
	case TypeGet:
		s.get(mes)
	case TypePut:
		s.put(mes)
	case TypeRemove:
		s.remove(mes)
	case TypeRemTag:
		s.removeTag(mes)
	case TypeFlush:
		s.flush()
	case TypeSetTTL:
		s.setTTL(mes)
	}
}

func (s *Shard) setTTL(mes *Request) {
	s.ticker.Stop()
	s.ticker = time.NewTicker(mes.TTL)
}

func (s *Shard) get(mes *Request) {
	s.RLock()
	entry, ok := s.entries[mes.Key.ID]
	s.RUnlock()

	out := NewRes()
	if ok && entry.Valid() {
		out.Data = entry.Data
		out.Ok = true
	}

	//mes.ResultChan <- out
	mes.ResultChan <- out
}

// func (s *Shard) shardCount(shardID int, stChan chan []int) (mes Request) {
// 	stChan <- []int{shardID, len(mCacheMemory.entries[shardID])}
// 	return
// }

func (s *Shard) put(mes *Request) {

	e := CreateEntry(mes.Key, mes.Data, mes.Compress, mes.Tags, mes.TTL)

	if !e.Valid() {
		return
	}

	s.Lock()
	defer s.Unlock()

	s.entries[mes.Key.ID] = e
	s.oldest = append(s.oldest, mes.Key.ID)

	return
}

func (s *Shard) removeTag(mes *Request) {
	s.Lock()
	defer s.Unlock()

	list := map[string]*Entry{}
	for key, entry := range s.entries {
		if !entry.Valid() {
			continue
		}

		valid := true
		for _, tag := range mes.Tags {
			if entry.Tags[tag] {
				valid = false
				break
			}
		}

		if valid {
			list[key] = entry
		}
	}
	s.entries = list

	// Update oldest
	oldest := []string{}
	for _, key := range s.oldest {
		if _, ok := s.entries[key]; ok {
			oldest = append(oldest, key)
		}
	}
	s.oldest = oldest
}

func (s *Shard) remove(mes *Request) {

	s.Lock()
	defer s.Unlock()

	if _, ok := s.entries[mes.Key.ID]; ok {
		s.entries[mes.Key.ID] = nil
		delete(s.entries, mes.Key.ID)
	}
}

func (s *Shard) flush() {
	s.entries = map[string]*Entry{}
}

func (s *Shard) trim() {

	forDelete := len(s.entries) - s.Len
	if forDelete < 0 {
		return
	}

	s.Lock()
	defer s.Unlock()

	s.oldest = s.oldest[forDelete:]
	list := map[string]*Entry{}
	for _, key := range s.oldest {
		if entry, ok := s.entries[key]; ok && entry.Valid() {
			list[key] = entry
		}
	}
	s.entries = list
}
