package GeeCahce

import (
	"GeeCache/Lru"
	"sync"
)

type cache struct {
	m          sync.Mutex
	lru        *Lru.Cache
	cacheBytes int64
	onEvicted  Lru.OnEvicted
}

func NewCache(maxBytes int64, onEvicted Lru.OnEvicted) *cache {
	return &cache{
		onEvicted:  onEvicted,
		cacheBytes: maxBytes,
	}
}

func (c *cache) add(key string, value ByteView) bool {
	c.m.Lock()
	defer c.m.Unlock()
	if c.lru == nil {
		c.lru = Lru.New(c.cacheBytes, c.onEvicted)
	}
	c.lru.Add(key, value)
	return true
}

func (c *cache) get(key string) (ByteView, bool) {
	c.m.Lock()
	defer c.m.Unlock()
	if c.lru == nil {
		return ByteView{}, false
	}
	if v, ok := c.lru.Get(key); !ok {
		return ByteView{}, false
	} else {
		return v.(ByteView), ok
	}

}
