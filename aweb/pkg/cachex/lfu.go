package cachex

import "container/heap"

type LFUCache struct {
	capacity int
	cache    map[int]*Item
	freqHeap *FrequencyHeap
}

type Item struct {
	key       int
	value     int
	frequency int
	index     int
}

type FrequencyHeap []*Item

func (h FrequencyHeap) Len() int           { return len(h) }
func (h FrequencyHeap) Less(i, j int) bool { return h[i].frequency < h[j].frequency }
func (h FrequencyHeap) Swap(i, j int) {
	h[i], h[j] = h[j], h[i]
	h[i].index = i
	h[j].index = j
}

func (h *FrequencyHeap) Push(x interface{}) {
	item := x.(*Item)
	item.index = len(*h)
	*h = append(*h, item)
}

func (h *FrequencyHeap) Pop() interface{} {
	old := *h
	n := len(old)
	item := old[n-1]
	*h = old[0 : n-1]
	return item
}

func NewLFUCache(capacity int) *LFUCache {
	return &LFUCache{
		capacity: capacity,
		cache:    make(map[int]*Item),
		freqHeap: &FrequencyHeap{},
	}
}

func (c *LFUCache) Get(key int) int {
	if item, found := c.cache[key]; found {
		item.frequency++
		heap.Fix(c.freqHeap, item.index)
		return item.value
	}
	return -1
}

func (c *LFUCache) Put(key int, value int) {
	if c.capacity == 0 {
		return
	}
	if item, found := c.cache[key]; found {
		item.value = value
		item.frequency++
		heap.Fix(c.freqHeap, item.index)
	} else {
		if len(c.cache) == c.capacity {
			least := heap.Pop(c.freqHeap).(*Item)
			delete(c.cache, least.key)
		}
		item := &Item{key: key, value: value, frequency: 1}
		heap.Push(c.freqHeap, item)
		c.cache[key] = item
	}
}
