package memorycache

import (
	"hash/crc32"
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
	return int(crc32.Checksum([]byte(key.ID), crc32.MakeTable(crc32.Koopman))) % maxId
}
