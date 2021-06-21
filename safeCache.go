package go_cache

import (
	"log"
	"sync"
)

//默认允许占用的最大内存
const DefaultMaxBytes = 1 << 29

//并发安全缓存
type safeCache struct {
	m sync.RWMutex
	cache Cache

	//记录缓存获取次数和命中次数
	nhit, nget int

	//sence sync.Map
}

func newSafeCache(cache Cache) *safeCache {
	return &safeCache{
		cache: cache,
	}
}

func (sc *safeCache) set(key string, value interface{}) {
	sc.m.Lock()
	defer sc.m.Unlock()
	sc.cache.Set(key,value)
}

func (sc *safeCache) get(key string) interface{} {
	sc.m.RLock()
	defer sc.m.RUnlock()
	sc.nget++
	if sc.cache == nil {
		return nil
	}

	v := sc.cache.Get(key)
	if v != nil {
		log.Println("[TourCache] hint")
		sc.nhit ++
	}
	return v
}

type Stat struct{
	NGet, NHit int
}

func (sc *safeCache) stat() *Stat {
	sc.m.RLock()
	defer sc.m.RUnlock()
	return &Stat{
		NGet: sc.nget,
		NHit: sc.nhit,
	}
}
