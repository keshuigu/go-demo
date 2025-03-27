package consistenthash

import (
	"hash/crc32"
	"slices"
	"strconv"
)

type Hash func(data []byte) uint32

type Map struct {
	hash     Hash
	replicas int // 虚拟节点倍数
	keys     []int
	hashMap  map[int]string // 虚拟节点映射真实节点
}

func New(replicas int, hash Hash) *Map {
	m := &Map{
		replicas: replicas,
		hash:     hash,
		hashMap:  make(map[int]string),
	}
	if m.hash == nil {
		m.hash = crc32.ChecksumIEEE
	}
	return m
}

// 创建节点，每个key对应一个节点
func (m *Map) Add(keys ...string) {
	for _, key := range keys {
		for i := 0; i < m.replicas; i++ {
			hash := int(m.hash([]byte(strconv.Itoa(i) + key)))
			m.keys = append(m.keys, hash)
			m.hashMap[hash] = key
		}
	}
	// 排序以确定节点顺寻
	slices.Sort(m.keys)
}

// 返回key对应的存储节点
func (m *Map) Get(key string) string {
	if key == "" {
		return ""
	}
	hash := int(m.hash([]byte(key)))
	idx, _ := slices.BinarySearch(m.keys, hash)
	return m.hashMap[m.keys[idx%len(m.keys)]]
}
