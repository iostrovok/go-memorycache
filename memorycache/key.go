package memorycache

import (
	"hash/fnv"
)

const sep byte = '_'

// Key defines cache key for some record
type Key struct {
	ID   string
	Tags []string
}

func NewKey(id string, tags ...string) Key {
	return Key{
		ID:   id,
		Tags: tags,
	}
}

func (key *Key) ShardID(maxId int) int {
	hash := fnv.New32()
	hash.Write([]byte(key.ID))
	return int(hash.Sum32()) % maxId
}
