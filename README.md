# go-cache
#### go-cache是一个线程安全的进程内缓存库，实现的缓存淘汰算法有:
- FIFO算法
- LFU算法
- LRU算法
获取库的方式:`go get -u github.com/The-half/go-cache`
#### Example
```go
//tourCache_test.go
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
```
