package cachex

import "container/list"

type LRUCache struct {
	capacity int
	cache    map[int]*list.Element
	list     *list.List
}

type Pair struct {
	key   int
	value int
}

func NewLRUCache(capacity int) *LRUCache {
	return &LRUCache{
		capacity: capacity,
		cache:    make(map[int]*list.Element),
		list:     list.New(),
	}
}

func (c *LRUCache) Get(key int) int {
	if elem, found := c.cache[key]; found {
		c.list.MoveToFront(elem)
		return elem.Value.(*Pair).value
	}
	return -1
}

func (c *LRUCache) Put(key int, value int) {
	if elem, found := c.cache[key]; found {
		c.list.MoveToFront(elem)
		elem.Value.(*Pair).value = value
	} else {
		if c.list.Len() == c.capacity {
			back := c.list.Back()
			c.list.Remove(back)
			delete(c.cache, back.Value.(*Pair).key)
		}
		pair := &Pair{key, value}
		elem := c.list.PushFront(pair)
		c.cache[key] = elem
	}
}
