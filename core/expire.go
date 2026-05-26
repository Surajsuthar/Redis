package core

import (
	"log"
	"time"
)

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

func DeletionExpireKey() {
	for {
		frq := ExpireSmple()
		if frq < 0.25 {
			break
		}
	}

	log.Println("Deleted Exprie key")
}
