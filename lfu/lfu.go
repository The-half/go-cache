package lfu

import (
	"container/heap"
	go_cache "go-cache"
	"sort"
)

type lfu struct {
	maxBytes int

	onEvicted func(key string, value interface{})

	usedBytes int

	queue *queue
	cache map[string]*entry
}

type entry struct {
	key string
	value interface{}
	//表示entry在ｑｕｅｕｅ中的权重，访问次数越多，权重越大
	weight int
	//表示在对中的缩影
	index int
}

func (e *entry) Len() int {
	return go_cache.CalcLen(e.value) + 4 + 4
}

//使用切片来构造最小堆
type queue []*entry

type Interface interface {
	sort.Interface
	Push(x interface{}) //add x as element Len()
	Pop() interface{} //remove and return element Len() - 1
}

func (q queue) Len() int{
	return len(q)
}

func (q queue) Less(i, j int) bool {
	return q[i].weight < q[j].weight
}

func (q queue) Swap( i, j int) {
	q[i], q[j] = q[j], q[i]
	q[i].index = i
	q[j].index = j
}

func (q *queue) Push(x interface{}) {
	n := len(*q)
	en := x.(*entry)
	en.index = n
	*q = append(*q, en)
}

func (q *queue) Pop() interface{} {
	old := *q
	n := len(old)
	en := old[n-1]
	old[n-1] = nil //防止对象逃逸
	en.index = -1
	*q = old[0:n-1]
	return en
}

func (q *queue) update(e *entry, value interface{}, weight int) {
	e.value = value
	e.weight = weight
	heap.Fix(q, e.index)
}

func New(maxBytes int, onEvicted func(key string, value interface{})) go_cache.Cache {
	q := make(queue, 0, 1024)
	return &lfu{
		maxBytes:  maxBytes,
		onEvicted: onEvicted,
		queue:     &q,
		cache:     make(map[string]*entry),
	}
}

func (l *lfu) Set(key string, value interface{}) {
	if e, ok := l.cache[key]; ok {
		l.usedBytes = l.usedBytes - go_cache.CalcLen(e.value) + go_cache.CalcLen(value)
		l.queue.update(e, value, e.weight + 1)
		return
	}
	en := &entry{key:key,value:value}
	//将en加入到最小堆中
	heap.Push(l.queue, en)

	l.cache[key] = en

	l.usedBytes += en.Len()
	if l.maxBytes > 0 && l.usedBytes > l.maxBytes {
		l.DelOldest()
	}
}

func (l *lfu) Get(key string) interface{} {
	if e, ok := l.cache[key]; ok {
		l.queue.update(e, e.value, e.weight)
		return e.value
	}
	return nil
}


func (l *lfu) Del(key string) {
	if e, ok := l.cache[key]; ok {
		heap.Remove(l.queue, e.index)
		l.removeElement(e)
	}
}

func (l *lfu) DelOldest() {
	if l.queue.Len() == 0 {
		return
	}
	l.removeElement(heap.Pop(l.queue))
}

func (l *lfu) Len() int {
	return l.queue.Len()
}

func (l *lfu) removeElement(x interface{}) {
	if x == nil {
		return
	}
	en := x.(*entry)
	delete(l.cache, en.key)
	l.usedBytes -= en.Len()

	if l.onEvicted != nil {
		l.onEvicted(en.key, en.value)
	}
}
