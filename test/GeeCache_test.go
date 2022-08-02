package test

import (
	"GeeCache"
	"GeeCache/Lru"
	"fmt"
	"log"
	"testing"
)

var db = map[string]string{
	"Tom":  "630",
	"Jack": "589",
	"Sam":  "567",
}

func TestGeeCache(t *testing.T) {
	GeeCahce.NewGroup("LowDB", GeeCahce.GetterFunc(func(key string) ([]byte, error) {
		log.Printf("[LowDB]Search Key %s", key)
		if v, ok := db[key]; ok {
			return []byte(v), nil
		} else {
			return nil, fmt.Errorf("[LowDB]not found key:%s", key)
		}
	}), 8, Lru.OnEvicted(func(key string, value Lru.Value) {
		log.Printf("[GeeCache]:Cache is Full,Element:{key:%s,value:%s}remove from cache", key, value)
	}))
	addr := "localhost:8080"
	peer := GeeCahce.NewHTTPPool(addr)
	peer.Log("peer is running at %s", addr)
	err := peer.Run()
	if err != nil {
		log.Fatal(err.Error())
	}
}
