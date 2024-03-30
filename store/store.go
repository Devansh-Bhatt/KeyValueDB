package store

import (
	"errors"
	"sync"
	"time"
)

type Db struct {
	mtx     sync.Mutex
	storage map[string]DbValue
}

type DbValue struct {
	Value  []byte
	expiry *time.Time
}

func NewDb() *Db {
	return &Db{
		storage: make(map[string]DbValue),
	}
}

func (db *Db) Set(key string, val []byte, expiry int64) {
	newDbVal := DbValue{
		Value: val,
	}
	if expiry != -1 {
		ValExpiry := time.Now().Add(time.Duration(expiry) * time.Millisecond)
		newDbVal.expiry = &ValExpiry
	}
	db.mtx.Lock()
	db.storage[key] = newDbVal
	db.mtx.Unlock()

}

func (db *Db) Get(key string) ([]byte, error) {
	db.mtx.Lock()
	defer db.mtx.Unlock()

	if value, exist := db.storage[key]; exist {
		if value.expiry != nil && time.Now().After(*value.expiry) {
			delete(db.storage, key)
			return nil, errors.New("Key Expired")
		}
		return value.Value, nil
	}
	return nil, errors.New("Key Not Found")
}
