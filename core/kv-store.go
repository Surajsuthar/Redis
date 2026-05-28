// In-memory key store operations used by command evaluation.
package core

import (
	"time"

	"github.com/Surajsuthar/go-redis/config"
)

var store map[string]*Obj

func init() {
	store = make(map[string]*Obj)
}

// NewObj creates a store object and translates a duration into an expiry time.
func NewObj(value interface{}, durationMs int64, te uint8, e uint8) *Obj {
	var expireAt int64 = -1
	if durationMs > 0 {
		expireAt = durationMs + time.Now().UnixMilli()
	}

	return &Obj{
		TypeEncoding: te | e,
		value:        value,
		ExpireAt:     expireAt,
	}
}

// PUT stores or replaces a key and updates keyspace statistics.
func PUT(key string, obj *Obj) {
	if len(store) >= config.MaxStoreSize {
		evict()
	}
	store[key] = obj
	if keyspaceStat[0] == nil {
		keyspaceStat[0] = make(map[string]int)
	}
	keyspaceStat[0]["key"]++
}

// GET returns a key if it exists and has not expired.
func GET(key string) *Obj {
	v := store[key]
	if v != nil {
		if v.ExpireAt != -1 && v.ExpireAt <= time.Now().UnixMilli() {
			delete(store, key)
			return nil
		}
	}
	return v
}

// DEL removes a key and reports whether it existed.
func DEL(key string) bool {
	_, exits := store[key]
	if exits {
		delete(store, key)
		if keyspaceStat[0] != nil {
			keyspaceStat[0]["key"]--
		}
		return true
	}
	return false
}
