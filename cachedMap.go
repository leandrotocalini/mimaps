package mimaps

import (
	"time"
)

type CachedMap[K comparable, V any] struct {
	data            map[K]V
	expireInSeconds int64
	expirationKeys  map[int64][]K
}

func (e *CachedMap[K, V]) Put(key K, value V) error {
	e.data[key] = value
	if e.expireInSeconds > 0 {
		now := time.Now() // current local time
		sec := now.Unix() + e.expireInSeconds
		if _, ok := e.expirationKeys[sec]; !ok {
			e.expirationKeys[sec] = []K{}
		}
		e.expirationKeys[sec] = append(e.expirationKeys[sec], key)
	}
	return nil
}

func (e *CachedMap[K, V]) Delete(key K) error {
	if _, err := e.Get(key); err != nil {
		return err
	}
	delete(e.data, key)
	return nil
}

func (e *CachedMap[K, V]) Get(key K) (V, error) {
	if val, ok := e.data[key]; ok {
		return val, nil
	}
	var result V // making a zero value of the type V
	return result, ErrInvalidKey
}

func (e *CachedMap[K, V]) run() {
	for {
		now := time.Now() // current local time
		sec := now.Unix()
		if val, ok := e.expirationKeys[sec]; ok {
			for _, key := range val {
				delete(e.data, key)
			}
			delete(e.expirationKeys, sec)
		}
		time.Sleep(1 * time.Second)
	}
}

func CreateCachedMap[K comparable, V any](expire int64) CachedMap[K, V] {
	e := CachedMap[K, V]{
		data:            make(map[K]V),
		expireInSeconds: expire,
		expirationKeys:  make(map[int64][]K),
	}
	if expire > 0 {
		go e.run()
	}
	return e
}
