package core

func evict() {
	for k := range store {
		delete(store, k)
		return
	}
}
