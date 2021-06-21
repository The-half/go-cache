package lru

import (
	"container/list"
	"go-cache"
	//"sync"
)

type lru struct {
	maxBytes int

	onEvicted func(key string, value interface{})

	usedBytes int

	ll *list.List

	cache map[string]*list.Element
}

type entry struct {
	key string
	value interface{}
}

func (e *entry) Len() int {
	return go_cache.CalcLen(e.value)
}

func (l *lru) Set(key string, value interface{}) {
	if e, ok := l.cache[key]; ok {
		//将元素移到队尾
		l.ll.MoveToBack(e)

		en := e.Value.(*entry)
		l.usedBytes = l.usedBytes - go_cache.CalcLen(en.value) + go_cache.CalcLen(value)
		en.value = value
		return
	}

	en := &entry{key, value}
	//将元素移到队尾
	e := l.ll.PushBack(en)
	l.cache[key] =e

	l.usedBytes += en.Len()

	if l.maxBytes > 0 && l.usedBytes > l.maxBytes {
		l.DelOldest()
	}

}

func (l *lru) Get(key string) interface{} {
	if e, ok := l.cache[key]; ok {
		//将e移到最后
		l.ll.MoveToBack(e)
		return e.Value.(*entry).value
	}

	return nil
}


func (l *lru) Del(key string) {
	if e, ok := l.cache[key]; ok {
		l.removeElement(e)
	}
}

func (l *lru) DelOldest() {
	l.removeElement(l.ll.Front())
}

func (l *lru) Len() int {
	return l.ll.Len()
}

func (l *lru) removeElement(e *list.Element) {
	if e == nil {
		return
	}
	//删除节点ｅ
	l.ll.Remove(e)
	en := e.Value.(*entry)
	l.usedBytes -= en.Len()
	delete(l.cache, en.key)

	if l.onEvicted != nil {
		l.onEvicted(en.key, en.value)
	}
}

func New(maxBytes int, onEvicted func(key string, value interface{})) go_cache.Cache {
	return &lru{
		maxBytes: maxBytes,
		onEvicted: onEvicted,
		ll: list.New(),
		cache: make(map[string]*list.Element),
	}
}
