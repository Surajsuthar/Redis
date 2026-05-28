// Eviction helpers used when the in-memory key store reaches its limit.
package core

import "github.com/Surajsuthar/go-redis/config"

// evictSimpleFirst removes the first key found in the map iteration order.
func evictSimpleFirst() {
	for k := range store {
		delete(store, k)
		return
	}
}

// evictAllKeyRandom removes a configured fraction of keys from the store.
func evictAllKeyRandom() {
	evictCount := int64(config.EvicationRatio * float64(config.KeyLimit))

	for key := range store {
		DEL(key)
		evictCount--
		if evictCount <= 0 {
			break
		}
	}
}

// evict selects the configured eviction policy and applies it.
func evict() {
	switch config.EcitType {
	case "simple-first":
		evictSimpleFirst()
	case "allkey-random":
		evictAllKeyRandom()
	default:

	}
}
