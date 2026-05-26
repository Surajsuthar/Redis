package core

import "time"

var store map[string]*Obj

type Obj struct {
	Value    interface{}
	ExpireAt int64
}

func init() {
	store = make(map[string]*Obj)
}

func NewObj(value interface{}, durationMs int64) *Obj {
	var expireAt int64 = -1
	if durationMs > 0 {
		expireAt = durationMs + time.Now().UnixMilli()
	}

	return &Obj{
		Value:    value,
		ExpireAt: expireAt,
	}
}

func PUT(key string, obj *Obj) {
	store[key] = obj
}

func GET(key string) *Obj {
	v := store[key]
	if v != nil {
		if v.ExpireAt <= time.Now().UnixMilli() {
			delete(store, key)
			return nil
		}
	}
	return v
}

func DEL(key string) bool {
	_, exits := store[key]
	if exits {
		delete(store, key)
		return true
	}
	return false
}
