package mycache

import (
	"errors"
	"log"
	"sync"
)

// 数据获取接口
type Getter interface {
	Get(key string) ([]byte, error)
}

// 接口型函数
// 可以将任何符合func(key string) ([]byte, error)的函数
// 转换为Getter接口的实现
type GetterFunc func(key string) ([]byte, error)

func (f GetterFunc) Get(key string) ([]byte, error) {
	return f(key)
}

// Group 表示一个缓存命名空间及其相关的数据加载分布
type Group struct {
	name      string
	getter    Getter
	mainCache cache
}

var (
	mu     sync.RWMutex
	groups = make(map[string]*Group)
)

func NewGroup(name string, cacheBytes int64, getter Getter) *Group {
	mu.Lock()
	defer mu.Unlock()
	if g, ok := groups[name]; ok {
		return g
	}
	if getter == nil {
		panic("getter is nil")
	}
	g := &Group{
		name:      name,
		getter:    getter,
		mainCache: cache{cacheBytes: cacheBytes},
	}
	groups[name] = g
	return g
}

// GetGroup 返回之前使用 NewGroup 创建的指定名称的 group，
// 如果没有找到该 group，则返回 nil。
func GetGroup(name string) *Group {
	mu.RLock()
	defer mu.RUnlock()
	return groups[name]
}

func (g *Group) Get(key string) (ByteView, error) {
	if key == "" {
		return ByteView{}, errors.New("key error")
	}
	if b, ok := g.mainCache.get(key); ok {
		log.Printf("[Cache hit]: key %s", key)
		return b, nil
	}
	return g.load(key)
}

// private

func (g *Group) load(key string) (ByteView, error) {
	return g.getLocally(key)
}

func (g *Group) getLocally(key string) (ByteView, error) {
	bytes, err := g.getter.Get(key)
	if err != nil {
		return ByteView{}, err
	}
	value := ByteView{b: cloneBytes(bytes)}
	g.populateCache(key, value)
	return value, nil
}

func (g *Group) populateCache(key string, value ByteView) {
	g.mainCache.add(key, value)
}
