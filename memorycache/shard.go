package memorycache

import (
	"sync"
	"time"
)

const (
	TTLClean = 500 * time.Millisecond
)

// MemoryCache storage memory
type Shard struct {
	sync.RWMutex

	Len     int
	entries map[string]*Entry
	inChan  chan *Request

	oldest []string
}

func NewShard(l int) *Shard {
	s := &Shard{
		Len:     l,
		entries: map[string]*Entry{},
		inChan:  make(chan *Request, 100),
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

	ticker := time.Tick(TTLClean)

	for {
		select {
		case <-ticker:
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
	case TypeFlush:
		s.flush()
	}
}

func (s *Shard) get(mes *Request) {
	s.RLock()
	entry, ok := s.entries[mes.Key.ID]
	s.RUnlock()

	out := NewRes()
	if ok {
		out.Entry = entry
		out.Ok = true
	}

	mes.ResultChan <- out
}

// func (s *Shard) shardCount(shardID int, stChan chan []int) (mes Request) {
// 	stChan <- []int{shardID, len(storageMemory.entries[shardID])}
// 	return
// }

func (s *Shard) put(mes *Request) {

	e := CreateEntry(mes.Key, mes.Data, mes.Compress)

	s.Lock()
	defer s.Unlock()

	s.entries[mes.Key.ID] = e
	s.oldest = append(s.oldest, mes.Key.ID)

	return
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
		if ent, ok := s.entries[key]; ok {
			list[key] = ent
		}
	}
	s.entries = list
}
