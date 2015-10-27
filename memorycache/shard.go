package memorycache

import (
	"sort"
	"sync"
	"time"
)

const (
	// 7.333 min
	TTLCleanDefault   = 60 * 7333 * time.Millisecond
	PercentBadDefault = 30
)

// MemoryCache mCache memory
type Shard struct {
	sync.RWMutex

	Len     int
	entries map[string]*Entry
	inChan  chan *Request

	// Counter of not-valid records
	countBad   int
	percentBad int

	oldest []string
	ticker *time.Ticker
	ttl    time.Duration

	pressFuncs   map[string]Press
	upPressFuncs map[string]UnPress
}

func NewShard(l int) *Shard {

	s := &Shard{
		Len:        l,
		entries:    map[string]*Entry{},
		inChan:     make(chan *Request, 100),
		ticker:     time.NewTicker(TTLCleanDefault),
		percentBad: PercentBadDefault,
		ttl:        TTLCleanDefault,
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
	case TypeSetPercentBad:
		s.setPercentBad(mes)
	case TypeGetTag:
		s.getTag(mes)
	}
}

func (s *Shard) setTTL(mes *Request) {
	s.ttl = mes.TTL
	s.ticker.Stop()
	s.ticker = time.NewTicker(mes.TTL)
}

func (s *Shard) setPercentBad(mes *Request) {
	i, ok := mes.Data.(int)
	if ok && i > -1 {
		s.percentBad = i
	}
}

func (s *Shard) get(mes *Request) {
	s.RLock()
	entry, ok := s.entries[mes.Key.ID]
	s.RUnlock()

	out := NewRes()
	if ok && entry.Valid() {
		out.Data = entry.Data
		out.Ok = true
	} else {
		s.countBad++
	}

	mes.ResultChan <- out
	s.checkCounterTrim()
}

func (s *Shard) getTag(mes *Request) {

	t := mes.Tags[0]

	s.RLock()
	list := []interface{}{}
	for _, entry := range s.entries {
		if !entry.Tags[t] {
			continue
		}

		if entry.Valid() {
			list = append(list, entry.Data)
		} else {
			s.countBad++
		}
	}
	s.RUnlock()

	out := NewRes()
	out.Data = list
	if len(list) > 0 {
		out.Ok = true
	}

	mes.ResultChan <- out
	s.checkCounterTrim()
}

// func (s *Shard) shardCount(shardID int, stChan chan []int) (mes Request) {
// 	stChan <- []int{shardID, len(mCacheMemory.entries[shardID])}
// 	return
// }

func (s *Shard) put(mes *Request) {

	e := CreateEntry(mes, s.pressFuncs)

	if !e.Valid() {
		return
	}

	s.Lock()
	defer s.Unlock()

	s.entries[mes.Key.ID] = e

	// if b, ok := data.(Body); ok {
	// 	mes.Body = b.Press()
	// 	mes.IsPressed = true
	// }

	s.oldest = append(s.oldest, mes.Key.ID)

	return
}

func (s *Shard) removeTag(mes *Request) {
	s.Lock()
	defer s.Unlock()

	noMustClean := true

	list := map[string]*Entry{}
	for key, entry := range s.entries {
		if !entry.Valid() {
			continue
		}

		valid := true
		for _, tag := range mes.Tags {
			if entry.Tags[tag] {
				valid = false
				noMustClean = false
				break
			}
		}

		if valid {
			list[key] = entry
		}
	}

	// Nothing was removed
	if noMustClean {
		return
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

func (s *Shard) checkCounterTrim() {

	if (100*s.countBad)/s.Len < s.percentBad {
		return
	}

	s.ticker.Stop()
	s.trim()
	s.ticker = time.NewTicker(s.ttl)
}

func (s *Shard) trim() {

	if (100*len(s.entries))/s.Len < s.percentBad {
		return
	}

	s.Lock()
	defer s.Unlock()

	good := 0
	news := []string{}
	list := map[string]*Entry{}
	for i := len(s.oldest) - 1; i >= 0; i-- {
		key := s.oldest[i]

		entry, ok := s.entries[key]
		if !ok || !entry.Valid() {
			continue
		}

		_, ok = list[key]
		if ok {
			continue
		}

		list[key] = entry
		news = append(news, key)
		good++

		if good >= s.Len {
			break
		}
	}

	sort.Reverse(sort.StringSlice(news))

	s.entries = list
	s.oldest = news
	s.countBad = 0
}
