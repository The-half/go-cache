package fast

import go_cache "go-cache"

type fastCache struct {
	//包含所有分片的切片,通常为1024
	shards []*cacheShard
	//方便通过位运算计算余数
	shardMask uint64
	//使用sum64
	hash fnc64a
}

func (c *fastCache) Set(key string, value interface{}) {
	c.getShard(key).set(key, value)
}

func (c *fastCache) Get(key string) interface{} {
	return c.getShard(key).get(key)
}

func (c *fastCache) Del(key string) {
	c.getShard(key).del(key)
}

func (c *fastCache) Len() int {
	length := 0
	for _, shard := range c.shards {
		length += shard.len()
	}
	return length
}

func (c *fastCache) DelOldest() {
	panic("no implements")
}
func NewFastCache(maxEntries ,shardsNum int, onEvicted func(key string, value interface{})) go_cache.Cache {
	fastCache := &fastCache{
		hash: newDefaultHasher(),
		shards: make([]*cacheShard, shardsNum),
		shardMask: uint64(shardsNum - 1),
	}
	for i := 0; i < shardsNum; i++ {
		fastCache.shards[i] = newCacheShard(maxEntries, onEvicted)
	}
	return fastCache
}

//获取分片
func (c *fastCache) getShard(key string) *cacheShard {
	hashedKey := c.hash.Sum64(key)
	return c.shards[hashedKey&c.shardMask]
}





