package core

import "github.com/Surajsuthar/go-redis/config"

func evictSimpleFirst() {
	for k := range store {
		delete(store, k)
		return
	}
}

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

func evict() {
	switch config.EcitType {
	case "simple-first":
		evictSimpleFirst()
	case "allkey-random":
		evictAllKeyRandom()
	default:

	}
}
