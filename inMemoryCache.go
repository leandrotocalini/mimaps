package mimaps

import (
	"sync"
	"time"
)

type InMemoryCache[K comparable, V any] struct {
	mu              *sync.RWMutex
	data            map[K]V
	expireInSeconds int64
	expirationKeys  map[int64][]K
}

func (e *InMemoryCache[K, V]) Put(key K, value V) error {
	e.mu.Lock()
	defer e.mu.Unlock()
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

func (e *InMemoryCache[K, V]) Delete(key K) error {
	e.mu.Lock()
	defer e.mu.Unlock()
	if _, ok := e.data[key]; ok {
		delete(e.data, key)

		return nil
	}
	return ErrInvalidKey
}

func (e *InMemoryCache[K, V]) Get(key K) (V, error) {
	e.mu.RLock()
	defer e.mu.RUnlock()
	if val, ok := e.data[key]; ok {
		return val, nil
	}
	var result V // making a zero value of the type V
	return result, ErrInvalidKey
}

func (e *InMemoryCache[K, V]) run() {
	for {
		now := time.Now() // current local time
		sec := now.Unix()
		e.mu.Lock()
		if val, ok := e.expirationKeys[sec]; ok {

			for _, key := range val {
				delete(e.data, key)

			}
			delete(e.expirationKeys, sec)
		}
		e.mu.Unlock()
		time.Sleep(1 * time.Second)
	}
}

func NewInMemoryCache[K comparable, V any](expire int64) InMemoryCache[K, V] {
	e := InMemoryCache[K, V]{
		mu:              &sync.RWMutex{},
		data:            make(map[K]V),
		expireInSeconds: expire,
		expirationKeys:  make(map[int64][]K),
	}
	if expire > 0 {
		go e.run()
	}
	return e
}
