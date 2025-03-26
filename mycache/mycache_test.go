package mycache_test

import (
	"keshuigu/mycache"
	"reflect"
	"testing"
)

func TestGetter(t *testing.T) {
	var f mycache.Getter = mycache.GetterFunc(func(key string) ([]byte, error) {
		return []byte(key), nil
	})

	expect := []byte("key")
	if v, _ := f.Get("key"); !reflect.DeepEqual(v, expect) {
		t.Fatalf("callback error")
	}
}
