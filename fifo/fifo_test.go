package fifo

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestFifo(t *testing.T) {
	assert := assert.New(t)

	cache := New(24, nil)
	cache.DelOldest()
	cache.Set("k1", 1)
	assert.Equal(cache.Len(), 1)
	v := cache.Get("k1")
	assert.Equal(1, v)

	cache.Del("k1")
	assert.Equal(0, cache.Len())
}


func TestFifo_OnEvicted(t *testing.T) {
	assert := assert.New(t)

	keys := make([]string, 0, 0)
	onEvicted := func(key string, value interface{}) {
		keys = append(keys, key)
	}
	//如果maxBytes为０，则表示不限制内存
	cache := New(0, onEvicted)

	cache.Set("k1", 1)
	cache.Set("k2", 2)
	//cache.Get("k1")
	cache.Set("k3", 3)
	cache.Set("k4", 4)
	//cache.Get("k2")
	cache.Set("k5", 5)
	cache.Set("k6", 6)

	expected := []string{"k1", "k2"}
	assert.Equal(expected, keys)
	assert.Equal(5, cache.Len())

}