package mimaps

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestExpirationMapWithInt(t *testing.T) {
	key := "test"
	val := 123
	miMap := NewInMemoryCache[string, int](int64(3))

	err := miMap.Put(key, val)
	assert.Equal(t, err, nil)

	rval, err := miMap.Get(key)
	assert.Equal(t, rval, val)
	assert.Equal(t, err, nil)
	time.Sleep(4 * time.Second)
	rval, err = miMap.Get(key)
	assert.Equal(t, rval, 0)
	assert.Equal(t, err, ErrInvalidKey)
}

func TestExpirationMapWithIntDelete(t *testing.T) {
	key := "test"
	val := 123
	miMap := NewInMemoryCache[string, int](int64(3))

	err := miMap.Put(key, val)
	assert.Equal(t, err, nil)
	rval, err := miMap.Get(key)
	assert.Equal(t, rval, val)
	assert.Equal(t, err, nil)
	err = miMap.Delete(key)
	assert.Equal(t, err, nil)
	err = miMap.Delete(key)
	assert.Equal(t, err, ErrInvalidKey)
	rval, err = miMap.Get(key)
	assert.Equal(t, rval, 0)
	assert.Equal(t, err, ErrInvalidKey)
}

func TestExpirationMapWithString(t *testing.T) {
	key := "test"
	val := "testOut"
	miMap := NewInMemoryCache[string, string](int64(3))

	err := miMap.Put(key, val)
	assert.Equal(t, err, nil)
	rval, err := miMap.Get(key)
	assert.Equal(t, rval, val)
	assert.Equal(t, err, nil)
	time.Sleep(4 * time.Second)
	rval, err = miMap.Get(key)
	assert.Equal(t, rval, "")
	assert.Equal(t, err, ErrInvalidKey)
}

func TestExpirationMapWithoutExpiration(t *testing.T) {
	key := "test"
	val := "testOut"
	miMap := NewInMemoryCache[string, string](int64(0))
	err := miMap.Put(key, val)
	assert.Equal(t, err, nil)
	time.Sleep(20 * time.Second)
	rval, err := miMap.Get(key)
	assert.Equal(t, rval, val)
	assert.Equal(t, err, nil)
}
