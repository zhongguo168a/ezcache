package ezcache

type IContext interface {
	GetCache() *Cache
	GetCacheGlobal() *Cache
}

type IContextRequest interface {
	GetCache() *Cache
	GetCacheGlobal() *Cache
	GetCacheRequest() *Cache
}
