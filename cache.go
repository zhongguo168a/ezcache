package ezcache

import (
	"context"
	"sync"
	"time"
)

type ICache interface {
	Set(key string, val interface{})
	Delete(key string)
	Expire(key string, duration time.Duration)
	Get(key string) (val interface{}, ok bool)
}

var (
	cacheGlobal *Cache
)

func NewCacheSyncNoExpire() *Cache {
	c := &Cache{
		data:         &sync.Map{},
		expiredSpace: time.Second * 5,
		expiredSet:   map[string]*expiredItem{},
	}
	return c
}

func NewCache() *Cache {
	c := &Cache{
		data:         &sync.Map{},
		expiredSpace: time.Second * 5,
		expiredSet:   map[string]*expiredItem{},
	}
	c.goCleanExpireList()
	return c
}

type Cache struct {
	// 缓存数据库中获取的数据
	data *sync.Map
	//
	mu sync.RWMutex
	//
	expiredSet map[string]*expiredItem
	//
	expiredRoot *expiredItem
	// 过期检查间隔
	expiredSpace time.Duration
}

func (c *Cache) Set(key string, val interface{}) {
	c.data.Store(key, val)
}

func (c *Cache) Delete(key string) {
	c.data.Delete(key)
}

func (c *Cache) Get(key string) (val interface{}, ok bool) {
	return c.data.Load(key)
}

// 从cache中获取对象, 如果存在, 通过has方法进行处理; 如果不存在, 通过create方法创建, 接着调用has方法
// 该方法没有返回值
func (c *Cache) GetOrNew(key string, has func(ival interface{}), create func() interface{}) {
	ival, ok := c.Get(key)
	if ok == false {
		ival = create()
		c.Set(key, ival)
	}
	has(ival)
	return
}

func NewContext(ctx context.Context) *Context {
	return &Context{
		Context: ctx,
	}
}

type Context struct {
	context.Context

	//
	cache *Cache
}

func (ctx *Context) GetCacheGlobal() *Cache {
	if cacheGlobal == nil {
		cacheGlobal = NewCache()
	}
	return cacheGlobal
}

func (ctx *Context) GetCache() *Cache {
	if ctx.cache == nil {
		ctx.cache = NewCache()
	}
	return ctx.cache
}
