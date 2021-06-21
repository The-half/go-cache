package fifo

import (
	"container/list"
	go_cache "go-cache"
)

type fifo struct {
	maxBytes int

	onEvicted func(key string, value interface{})

	usedBytes int

	ll *list.List

	cache map[string]*list.Element
}

func (f *fifo) Set(key string, value interface{}) {
	if e, ok := f.cache[key]; ok {
		//移到最后
		f.ll.MoveToBack(e)
		en := e.Value.(*entry)
		f.usedBytes = f.usedBytes - go_cache.CalcLen(en.value) + go_cache.CalcLen(value)
		en.value = value
		return
	}

	en := &entry{key, value}
	//将新元素放到队尾
	e := f.ll.PushBack(en)
	f.cache[key] = e

	f.usedBytes += en.Len()
	if f.maxBytes > 0 && f.usedBytes > f.maxBytes {
		f.DelOldest()
	}
}

func (f *fifo) Get(key string) interface{} {
	if e, ok := f.cache[key]; ok {
		return e.Value.(*entry).value
	}
	return nil
}

func (f *fifo) Del(key string) {
	if e, ok := f.cache[key]; ok {
		f.removeElement(e)
	}
}

func (f *fifo) DelOldest() {
	//删除头结点
	f.removeElement(f.ll.Front())
}

func (f *fifo) Len() int {
	return f.ll.Len()
}

func (f *fifo) removeElement(e *list.Element) {
	if e == nil {
		return
	}
	//删除节点
	f.ll.Remove(e)
	en := e.Value.(*entry)
	f.usedBytes -= en.Len()
	delete(f.cache, en.key)

	if f.onEvicted != nil {
		f.onEvicted(en.key, en.value)
	}
}

type entry struct {
	key string
	value interface{}
}

func (e *entry) Len() int {
	return go_cache.CalcLen(e.value)
}

func New(maxBytes int, onEvicted func(key string, value interface{})) go_cache.Cache {
	return &fifo{
		maxBytes: maxBytes,
		onEvicted: onEvicted,
		ll: list.New(),
		cache: make(map[string]*list.Element),
	}
}