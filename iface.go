package ezcache

import "context"

type IContext interface {
	context.Context

	GetCache() *Cache
	GetCacheGlobal() *Cache
}

type IContextRequest interface {
	IContext

	GetCacheRequest() *Cache
}
