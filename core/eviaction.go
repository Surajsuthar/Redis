package core

import "github.com/Surajsuthar/go-redis/config"

func evictSimpleFirst() {
	for k := range store {
		delete(store, k)
		return
	}
}

func evictAllKeyRandom() {
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
