package consistenthash

import (
	"hash/crc32"
	"sort"
	"strconv"
	"sync"
)

type Hash func(key []byte) uint32

//
// Map
// @Description: Map 是一致性哈希算法的主数据结构，包含 4 个成员变量：Hash 函数 hash；
//虚拟节点倍数 replicas；哈希环 keys；虚拟节点与真实节点的映射表 hashMap，键是虚拟节点的哈希值，值是真实节点的名称
//
type Map struct {
	hash     Hash
	replicas int
	keys     []int
	hashMap  map[int]string
	mutex    sync.Mutex
}

func New(replicas int, hash Hash) *Map {
	m := &Map{
		hash:     hash,
		replicas: replicas,
		keys:     nil,
		hashMap:  make(map[int]string),
	}
	if m.hash == nil {
		m.hash = crc32.ChecksumIEEE
	}
	return m
}

//
//Add
//@Description: 对每一个真实节点 key，对应创建 m.replicas 个虚拟节点， 虚拟节点的名称是：strconv.Itoa(i) + key，即通过添加编号的方式区分不同虚拟节点。
//@receiver m Maps
//@param keys 节点的key值
//
func (m *Map) Add(keys ...string) {
	m.mutex.Lock()
	for _, key := range keys {
		for i := 0; i < m.replicas; i++ {
			hash := int(m.hash([]byte(strconv.Itoa(i) + key)))
			m.keys = append(m.keys, hash)
			m.hashMap[hash] = key
		}
	}
	sort.Ints(m.keys)
	m.mutex.Unlock()
}

//
//Get
//@Description: 选择节点就非常简单了，第一步，计算 key 的哈希值。
//第二步，顺时针找到第一个匹配的虚拟节点的下标 idx，从 m.keys 中获取到对应的哈希值。如果 idx == len(m.keys)，说明应选择 m.keys[0]，因为 m.keys 是一个环状结构，所以用取余数的方式来处理这种情况。
//第三步，通过 hashMap 映射得到真实的节点
//@receiver m
//@param key
//@return string 真实节点的string
//
func (m *Map) Get(key string) string {
	m.mutex.Lock()
	if len(m.keys) == 0 {
		return ""
	}
	hash := int(m.hash([]byte(key)))
	idx := sort.Search(len(m.keys), func(i int) bool {
		return m.keys[i] >= hash
	})
	m.mutex.Unlock()
	return m.hashMap[m.keys[idx%len(m.keys)]]
}
