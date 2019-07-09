package ezcache

import "context"

func NewRequestContextDefault() (r *RequestContext) {
	return NewRequestContext(context.Background())
}

func NewRequestContext(ctx context.Context) (r *RequestContext) {
	r = &RequestContext{
		request: NewCache(),
	}
	r.Context = *NewContext(ctx)
	return
}

type RequestContext struct {
	Context

	request *Cache
}

func (ctx *RequestContext) GetCacheRequest() *Cache {
	return ctx.request
}
