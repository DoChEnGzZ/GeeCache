package Lru

import "container/list"

// Cache My cache
type Cache struct {
	MaxBytes  int64                         //申请的最大内存空间
	NowBytes  int64                         //当前使用的内存空间
	List      *list.List                    //LRU队列
	cache     map[string]*list.Element      //Cache主体
	OnEvicted func(key string, value Value) //回调函数,当去除最进未使用的内存块时调用该函数
}

type Entry struct {
	key   string
	value Value
}

type Value interface {
	Len() int
}

func New(maxBytes int64, onEvicted func(key string, value Value)) *Cache {
	return &Cache{
		MaxBytes:  maxBytes,
		List:      list.New(),
		cache:     make(map[string]*list.Element),
		OnEvicted: onEvicted,
	}
}

func (c *Cache) Add(key string, value Value) {
	if ele, ok := c.cache[key]; ok {
		c.List.MoveToFront(ele)
		kv := ele.Value.(Entry)
		c.NowBytes += int64(value.Len() - kv.value.Len())
		kv.value = value
	} else {
		ele := c.List.PushFront(Entry{
			key:   key,
			value: value,
		})
		c.cache[key] = ele
		c.NowBytes += int64(value.Len() + len(key))
	}
	for c.MaxBytes != 0 && c.NowBytes > c.MaxBytes {
		c.RemoveOldEle()
	}
}

// Get 获得内存块
func (c *Cache) Get(key string) (value Value, ok bool) {
	if ele, ok := c.cache[key]; ok {
		c.List.MoveToFront(ele)
		kv := ele.Value.(Entry)
		return kv.value, true
	} else {
		return nil, false
	}
}

// RemoveOldEle 移除队列尾的最不常用内存块
func (c *Cache) RemoveOldEle() {
	ele := c.List.Back()
	if ele != nil {
		c.List.Remove(ele)
		kv := ele.Value.(Entry)
		delete(c.cache, kv.key)
		c.NowBytes -= int64(len(kv.key)) + int64(kv.value.Len())
		if c.OnEvicted != nil {
			c.OnEvicted(kv.key, kv.value)
		}
	}
}

func (c *Cache) Len() int {
	return c.List.Len()
}
