package consistenthash_test

import (
	"cache/mycache/consistenthash"
	"strconv"
	"testing"
)

func TestHashing(t *testing.T) {
	hash := consistenthash.New(3, func(key []byte) uint32 {
		i, _ := strconv.Atoi(string(key))
		return uint32(i)
	})
	// Given the above hash function, this will give replicas with "hashes":
	// 2, 4, 6, 12, 14, 16, 22, 24, 26
	hash.Add("6", "4", "2")

	// value是对应的节点
	testCases := map[string]string{
		"2":  "2",
		"11": "2",
		"23": "4",
		"27": "2",
	}

	for k, v := range testCases {
		if hash.Get(k) != v {
			t.Fatalf("Asking for %s, should have yielded %s", k, v)
		}
	}

	hash.Add("8")
	testCases["27"] = "8"
	for k, v := range testCases {
		if hash.Get(k) != v {
			t.Fatalf("Asking for %s, should have yielded %s", k, v)
		}
	}
}
