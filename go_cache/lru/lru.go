package lru

import "container/list"

type Value interface {
	Len() int
}

type Entry struct {
	key   string
	value Value
}

type Cache struct {
	maxBytes  int64      //内存最大值
	nbytes    int64      //内存当前使用的值
	ll        *list.List //双向链表
	cache     map[string]*list.Element
	OnEvicted func(key string, value Value)
}

func (c *Cache) Len() int {
	return c.ll.Len()
}

// initialize
func New(maxBytes int64, onEvicted func(string, Value)) *Cache {
	return &Cache{
		maxBytes:  maxBytes,
		ll:        list.New(),
		cache:     make(map[string]*list.Element),
		OnEvicted: onEvicted,
	}
}

func (c *Cache) Get(key string) (value Value, ok bool) {
	if ele, ok := c.cache[key]; ok {
		c.ll.MoveToFront(ele)
		kv := ele.Value.(*Entry)
		return kv.value, true
	}
	return
}

func (c *Cache) RemoveOldest() {
	ele := c.ll.Back()
	if ele != nil {
		c.ll.Remove(ele)
		kv := ele.Value.(*Entry)
		delete(c.cache, kv.key)
		c.nbytes -= int64(len(kv.key)) + int64(kv.value.Len())
		if c.OnEvicted != nil {
			c.OnEvicted(kv.key, kv.value)
		}
	}
}

func (c *Cache) Add(key string, value Value) {
	if ele, ok := c.cache[key]; ok {
		c.ll.MoveToFront(ele)
		kv := ele.Value.(*Entry)
		c.nbytes += int64((value.Len())) - int64(kv.value.Len())
		kv.value = value
	} else {
		//先计算，如果超出最大内存，则先删除键和对应的节点,再添加
		c.nbytes += int64(len(key)) + int64(value.Len())
		for c.maxBytes != 0 && c.maxBytes < c.nbytes {
			c.RemoveOldest()
		}
		ele := c.ll.PushFront(&Entry{key, value})
		c.cache[key] = ele
	}

}
