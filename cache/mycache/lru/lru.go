package lru

import (
	"container/list"
)

type Cache struct {
	maxBytes int64
	nbytes   int64
	ll       *list.List
	cache    map[string]*list.Element
	// 移除最近最少使用时的回调函数
	OnEvicted func(key string, value Value)
}

type entry struct {
	key   string
	value Value
}

type Value interface {
	Len() int
}

func New(maxBytes int64, onEvicted func(string, Value)) *Cache {
	return &Cache{
		maxBytes:  maxBytes,
		ll:        list.New(),
		cache:     make(map[string]*list.Element),
		OnEvicted: onEvicted,
	}
}

func (c *Cache) Get(key string) (value Value, ok bool) {
	if e, ok := c.cache[key]; ok {
		c.ll.MoveToFront(e)
		kv := e.Value.(*entry)
		return kv.value, true
	}
	// 返回nil, false
	return
}

func (c *Cache) RemoveOldest() {
	e := c.ll.Back()
	if e != nil {
		kv := e.Value.(*entry)
		bytes := int64(len(kv.key)) + int64(kv.value.Len())
		c.nbytes -= bytes
		delete(c.cache, kv.key)
		c.ll.Remove(e)
		if c.OnEvicted != nil {
			c.OnEvicted(kv.key, kv.value)
		}
	}
}

func (c *Cache) Add(key string, value Value) {
	// 考虑修改
	if e, ok := c.cache[key]; ok {
		c.ll.MoveToFront(e)
		kv := e.Value.(*entry)
		bytesDiff := int64(value.Len()) - int64(kv.value.Len())
		// 确保空间足够
		for c.nbytes+bytesDiff > c.maxBytes {
			c.RemoveOldest()
		}
		kv.value = value
		c.nbytes += bytesDiff
	} else { // 新增情况
		bytes := int64(len(key)) + int64(value.Len())
		for c.nbytes+bytes > c.maxBytes {
			c.RemoveOldest()
		}
		kv := c.ll.PushFront(&entry{
			key:   key,
			value: value,
		})
		c.cache[key] = kv
		c.nbytes += int64(bytes)
	}
}

func (c *Cache) Len() int {
	return c.ll.Len()
}
