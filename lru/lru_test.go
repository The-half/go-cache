package lru

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestSetGet(t *testing.T) {
	assert := assert.New(t)

	keys := make([]string, 0, 8)
	onEvicted := func(key string, value interface{}) {
		keys = append(keys , key)
	}

	cache := New(8, onEvicted)
	cache.Set("k1", 1)
	cache.Set("k2", 2)
	cache.Get("k1")
	cache.Set("k3", 3)
	cache.Get("k1")
	cache.Set("k4", 4)

	expected := []string{"k2", "k3"}
	assert.Equal(expected, keys)
	assert.Equal(cache.Len(), 2)
}