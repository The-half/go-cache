package lfu

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestLfu_Del(t *testing.T) {
	assert := assert.New(t)

	cache := New(24, nil)
	cache.DelOldest()
	cache.Set("k1",1)
	v := cache.Get("k1")
	assert.Equal(v, 1)

	cache.Del("k1")
	assert.Equal(cache.Len(), 0)
}

func TestOnEvicted(t *testing.T) {
	assert := assert.New(t)

	keys := make([]string, 0, 8)
	onEvicted := func(key string, value interface{}) {
		keys = append(keys, key)
	}

	cache := New(32, onEvicted)

	cache.Set("k1",1)
	cache.Set("k2",2)
	cache.Get("k1")
	cache.Get("k1")
	cache.Get("k2")
	cache.Set("k3",3)
	cache.Set("k4",4)

	expected := []string{"k1", "k2"}
	assert.Equal(expected, keys)
	assert.Equal(cache.Len(), 1)
}