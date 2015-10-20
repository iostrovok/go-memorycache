package memorycache

import "time"

type Act uint

// Storage types
const (
	TypeGet    Act = iota
	TypeGetTag Act = iota
	TypePut    Act = iota
	TypeRemove Act = iota
	TypeRemTag Act = iota
	TypeFlush  Act = iota
	TypeCount  Act = iota
	TypeSetTTL Act = iota
)

type Res struct {
	Entry *Entry
	Ok    bool
	Count int
}

type Request struct {
	Action     Act
	Key        Key
	Data       interface{}
	Compress   Compress
	ResultChan chan *Res
	Tags       []string
	TTL        time.Duration
}

func NewRes() *Res {
	return &Res{
		Entry: EmptyEntry(),
		Ok:    false,
		Count: 0,
	}
}

func NewRequest(act Act, k Key, TTL ...time.Duration) *Request {
	return &Request{
		Action:     act,
		Key:        k,
		ResultChan: make(chan *Res, 1),
	}
}
