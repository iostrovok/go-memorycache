package memorycache

var Singleton *MemoryCache = nil

// NewMemoryCache returns new initializated instance of MemoryCache
func NewSingleton(shards int, maxEntriesPerShard int) {
	Singleton = New(shards, maxEntriesPerShard)
}

// Close storage memory channels
func Close() {
	Singleton.Close()
	Singleton = nil
}

// Get returns data by key
func Get(k string) (data interface{}, ok bool) {
	return Singleton.Get(k)
}

// Put puts new data in storage
func Put(data interface{}, k string, tags ...string) {
	Singleton.Put(data, k, tags...)
}

// Remove remove data in cache by key
func Remove(k string) {
	Singleton.Remove(k)
}

// Flush removes all entries from cache and returns number of flushed entries
func Flush() {
	Singleton.Flush()
}
