package GeeCahce

import (
	"GeeCache/Lru"
	"fmt"
	"log"
	"sync"
)

// Getter
// @Description: interface of getting the data
//
type Getter interface {
	Get(key string) ([]byte, error)
}

type GetterFunc func(key string) ([]byte, error)

func (f GetterFunc) Get(key string) ([]byte, error) {
	return f(key)
}

type Group struct {
	name      string
	getter    Getter
	mainCache *cache
}

var (
	m      sync.Mutex
	groups = make(map[string]*Group)
)

//
//NewGroup
//@Description: return a new group
//@param name:the name of group
//@param getter: getter function
//@param maxBytes: max bytes of cache
//@param evicted: onEvicted function
//@return *Group
//
func NewGroup(name string, getter Getter, maxBytes int64, evicted Lru.OnEvicted) *Group {
	if getter == nil {
		panic("getter is nil")
	}
	m.Lock()
	defer m.Unlock()
	g := &Group{
		name:      name,
		getter:    getter,
		mainCache: NewCache(maxBytes, evicted),
	}
	groups[name] = g
	return g
}

func GetGroup(name string) *Group {
	m.Lock()
	defer m.Unlock()
	if g, ok := groups[name]; ok {
		return g
	} else {
		return nil
	}
}

//
//Get
//@Description: get the value in group
//@receiver g
//@param key the description string of value
//@return ByteView
//@return error
//
func (g Group) Get(key string) (ByteView, error) {
	if key == "" {
		return ByteView{}, fmt.Errorf("[GeeCache]key is null\n")
	}
	if v, ok := g.mainCache.get(key); ok {
		log.Printf("[GeeCache]:Cache hit the %s,value is %s\n", key, v)
		return v, nil
	}
	return g.Load(key)
}

func (g Group) Load(key string) (ByteView, error) {
	return g.getLocally(key)
}

//
//getLocally
//@Description:when cache do not have the value,get it by Getter
//@receiver g
//@param key the description string of value
//@return ByteView
//@return error
//
func (g Group) getLocally(key string) (ByteView, error) {
	v, err := g.getter.Get(key)
	if err != nil {
		return ByteView{}, err
	}
	value := ByteView{b: v}
	g.populateCache(key, value)
	log.Printf("[GeeCache]:not hit key %s,get it from Getter function\n", key)
	return value, nil
}

//
//populateCache
//@Description: add the value to the cache when the value is not in it
//@receiver g
//@param key the description string of value
//@param value ByteView
//
func (g Group) populateCache(key string, value ByteView) {
	g.mainCache.add(key, value)
}
