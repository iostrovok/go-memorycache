package memorycache

import (
	"time"
)

type Compress uint

const (
	// States
	Nothing Compress = iota
	JSON    Compress = iota
)

// Entry contains data
type Entry struct {
	Key         Key
	CreateDate  time.Time
	EndDate     time.Time
	CompressTag Compress
	Data        interface{}
	Tags        map[string]bool
	CheckTime   bool
}

func EmptyEntry() *Entry {
	return &Entry{
		Key:        NewKey(""),
		CreateDate: time.Now(),
		Data:       nil,
		Tags:       map[string]bool{},
	}
}

// CreateEntry returns new instance of Entry
func CreateEntry(k Key, data interface{}, comp Compress, tags []string, TTL time.Duration) *Entry {
	out := &Entry{
		Key:         k,
		CreateDate:  time.Now(),
		Data:        data,
		Tags:        map[string]bool{},
		CompressTag: comp,
	}

	if TTL != 0 {
		out.CheckTime = true
		out.EndDate = time.Now().Add(TTL)
	}

	for _, s := range tags {
		out.Tags[s] = true
	}

	return out
}

func (e *Entry) Valid() bool {
	if !e.CheckTime {
		return true
	}
	return time.Now().Before(e.EndDate)
}
