package ezcache

func NewRequestContextDefault() (r *RequestContext) {
	return NewRequestContext(NewContext())
}

func NewRequestContext(ctx IContext) (r *RequestContext) {
	r = &RequestContext{}
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
