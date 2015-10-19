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
	CompressTag Compress
	Data        interface{}
}

func EmptyEntry() *Entry {
	return &Entry{
		Key:        NewKey(""),
		CreateDate: time.Now(),
		Data:       nil,
	}
}

// CreateEntry returns new instance of Entry
func CreateEntry(k Key, data interface{}, CompressTag ...Compress) *Entry {
	out := &Entry{
		Key:        k,
		CreateDate: time.Now(),
		Data:       data,
	}

	if len(CompressTag) > 0 {
		if CompressTag[0] != Nothing {
			out.CompressTag = CompressTag[0]
		}
	}

	return out
}
