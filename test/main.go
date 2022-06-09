package main

import (
	GeeCahce "GeeCache"
	"GeeCache/Lru"
	"fmt"
	"log"
)

var db = map[string]string{
	"Tom":  "630",
	"Jack": "589",
	"Sam":  "567",
}

func main() {
	geeGroup := GeeCahce.NewGroup("LowDB", GeeCahce.GetterFunc(func(key string) ([]byte, error) {
		log.Printf("[LowDB]Search Key %s", key)
		if v, ok := db[key]; ok {
			return []byte(v), nil
		} else {
			return nil, fmt.Errorf("[LowDB]not found key:%s", key)
		}
	}), 8, Lru.OnEvicted(func(key string, value Lru.Value) {
		log.Printf("[GeeCache]:Cache is Full,Element:{key:%s,value:%s}remove from cache", key, value)
	}))
	for k, v := range db {
		if view, err := geeGroup.Get(k); err != nil || view.String() != v {
			log.Fatal("[GeeCahce],get error" + err.Error())
		}
		if _, err := geeGroup.Get(k); err != nil {
			log.Fatal("[GeeCahce],get error" + err.Error())
		}
	}
}
