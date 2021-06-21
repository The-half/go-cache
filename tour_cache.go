package go_cache

type Getter interface {
	Get(key string) interface{}
}

type GetFunc func(key string) interface{}

func (f GetFunc) Get(key string) interface{} {
	return f(key)
}

type TourCache struct {
	mainCache *safeCache
	getter Getter
}

func NewTourCache(getter Getter, cache Cache) *TourCache {
	return &TourCache{
		//并发安全的缓存实现
		mainCache: newSafeCache(cache),
		//回调，用于缓存未命中时从数据源获取的数据
		getter: getter,
	}
}

func (t *TourCache) Get(key string) interface{} {
	val := t.mainCache.get(key)
	if val != nil{
		return val
	}
	if t.getter != nil {
		val = t.getter.Get(key)
		if val == nil {
			return nil
		}
		t.mainCache.set(key, val)
		return val
	}
	return nil
}

func (t *TourCache) Stat() *Stat {
	return t.mainCache.stat()
}

func (t *TourCache) Set(key string, val interface{}) {
	if val == nil {
		return
	}
	t.mainCache.set(key, val)
}