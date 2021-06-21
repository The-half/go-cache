package go_cache_test

import (
	"github.com/stretchr/testify/assert"
	go_cache "go-cache"
	"go-cache/lru"
	"log"
	"sync"
	"testing"
)

func TestTourCacheSet(t *testing.T) {
	assert := assert.New(t)
	
	db := map[string]string{
		"key1":"val1",
		"key2":"val2",
		"key3":"val3",
		"key4":"val4",
		//"key5":"val5",
	}

	getter := go_cache.GetFunc(func(key string) interface{} {
		log.Println("[from DB] find key", key)
		if val, ok := db[key]; ok {
			return val
		}
		return nil
	})
	tourCache := go_cache.NewTourCache(getter, lru.New(0, nil))

	var wg sync.WaitGroup

	for k, v := range db {
		wg.Add(1)
		go func(k, v string) {
			defer wg.Done()
			assert.Equal(tourCache.Get(k), v)
			assert.Equal(tourCache.Get(k), v)
		}(k,v)
	}
	wg.Wait()

	assert.Equal(tourCache.Get("unknown"), nil)
	assert.Equal(tourCache.Get("unknown"), nil)

	assert.Equal(tourCache.Stat().NGet, 9)//10
	assert.Equal(tourCache.Stat().NHit, 3)//4

}

