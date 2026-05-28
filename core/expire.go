// Expiration utilities for deleting keys whose TTL has passed.
package core

import (
	"log"
	"time"
)

// ExpireSmple samples a small group of expiring keys and removes stale ones.
func ExpireSmple() float32 {
	var limit int = 20
	var expriedCount int = 0

	for key, obj := range store {
		if obj.ExpireAt != -1 {
			limit--

			if obj.ExpireAt <= time.Now().UnixMilli() {
				delete(store, key)
				expriedCount++
			}
		}

		if limit == 0 {
			break
		}
	}

	return float32(expriedCount) / float32(20.0)
}

// DeletionExpireKey repeats expiration sampling while many sampled keys expire.
func DeletionExpireKey() {
	for {
		frq := ExpireSmple()
		if frq < 0.25 {
			break
		}
	}

	log.Println("Deleted Exprie key")
}
