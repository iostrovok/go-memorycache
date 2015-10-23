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
	Body        []byte
}

func EmptyEntry() *Entry {
	return &Entry{
		Key:        NewKey(""),
		CreateDate: time.Now(),
		Data:       nil,
		Tags:       map[string]bool{},
		Body:       []byte{},
	}
}

func _press(data interface{}, PF map[string]Press, tag string) (interface{}, bool) {
	f, ok := PF[tag]
	if !ok {
		return nil, false
	}

	out, err := f(data)
	if err != nil {
		return nil, false
	}

	return out, true
}

// CreateEntry returns new instance of Entry
func CreateEntry(k Key, data interface{}, comp Compress, tags []string, TTL time.Duration, PF map[string]Press) *Entry {
	out := &Entry{
		Key:         k,
		CreateDate:  time.Now(),
		Tags:        map[string]bool{},
		CompressTag: comp,
	}

	if TTL != 0 {
		out.CheckTime = true
		out.EndDate = time.Now().Add(TTL)
	}

	pressed := false
	for _, s := range tags {
		out.Tags[s] = true
		if !pressed {
			out.Data, pressed = _press(data, PF, s)
		}
	}

	// Get press function for default value (if it exists)
	if !pressed {
		out.Data, pressed = _press(data, PF, "")
	}

	if !pressed {
		out.Data = data
	}

	return out
}

func (e *Entry) Valid() bool {
	if !e.CheckTime {
		return true
	}
	return time.Now().Before(e.EndDate)
}
