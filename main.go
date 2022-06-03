package main

import (
	"GeeCache/Lru"
	"fmt"
)

type vstring string

func (s vstring) Len() int {
	return len(s)
}

func main() {
	lru := Lru.New(8, func(key string, value Lru.Value) {
		fmt.Printf("key:%s,value:%s in back,has been remove", key, value)
	})
	lru.Add("cheng", vstring("18"))
	lru.Add("che", vstring("18"))
	lru.Add("c", vstring("18"))
}
