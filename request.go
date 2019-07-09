package ezcache

import "context"

func NewRequestContextDefault() (r *RequestContext) {
	return NewRequestContext(NewContext(context.Background()))
}

func NewRequestContext(ctx IContext) (r *RequestContext) {
	r = &RequestContext{
		request: NewCache(),
	}
	r.IContext = ctx
	return
}

type RequestContext struct {
	IContext

	request *Cache
}

func (ctx *RequestContext) GetCacheRequest() *Cache {
	return ctx.request
}
