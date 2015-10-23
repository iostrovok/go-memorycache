package memorycache

import (
	"time"
)

var Singleton *MemoryCache = nil

// NewMemoryCache returns new initializated instance of MemoryCache
func NewSingleton(shards int, totalLimit int) {
	Singleton = New(shards, totalLimit)
}

func TTLClean(t time.Duration) {
	Singleton.TTLClean(t)
}

// Close mCache memory channels
func Close() {
	Singleton.Close()
	Singleton = nil
}

// Get returns data by key
func Get(k string) (data interface{}, ok bool) {
	return Singleton.Get(k)
}

// Put puts new data in mCache _WITHOUT_ TTL
func Put(data interface{}, k string, tags ...string) {
	Singleton.Put(data, k, tags...)
}

// PutTTL puts new data in mCache _WITH_ TTL
func PutTTL(data interface{}, k string, TTL time.Duration, tags ...string) {
	Singleton.PutTTL(data, k, TTL, tags...)
}

// Remove remove data in cache by key
func Remove(k string) {
	Singleton.Remove(k)
}

func RemoveTag(tag string) {
	Singleton.RemoveTag(tag)
}

// Flush removes all entries from cache and returns number of flushed entries
func Flush() {
	Singleton.Flush()
}
