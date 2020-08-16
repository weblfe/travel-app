package middlewares

import (
		"github.com/astaxie/beego/context"
		"github.com/weblfe/travel-app/libs"
		"time"
)

// 请求 头部参数处理中间件
type HeaderWare struct {
		middlewareImpl
}

const (
		Header               = "header"
		HeaderMiddleWareName = "api.header"
		RequestId            = "ApiRequestId"
)

var (
		headerMiddlewareInstance *HeaderWare
)

func GetHeaderWares() *HeaderWare {
		if headerMiddlewareInstance == nil {
				getLock(Header).Do(newHeaderWare)
		}
		return headerMiddlewareInstance
}

func newHeaderWare() {
		var ware = new(HeaderWare)
		ware.Init()
}

func (this *HeaderWare) Init() {
		this.Name = HeaderMiddleWareName
		this.SetHandler(this.done)
		Register(this.Middleware(), this)
}

func (this *HeaderWare) done(ctx *context.Context) bool {
		var handler = this.getHandler(ctx)
		return handler(ctx)
}

func (this *HeaderWare) getHandler(ctx *context.Context) func(*context.Context) bool {
		var id = this.getId(ctx)
		if id != "" {
				return this.next
		}
		return this.before
}

func (this *HeaderWare) getId(ctx *context.Context) string {
		var v = ctx.Input.GetData(RequestId)
		if v == nil || v == "" {
				return ""
		}
		return v.(string)
}

func (this *HeaderWare) setId(ctx *context.Context) {
		ctx.Input.SetData(RequestId, libs.HashCode(time.Now().Unix()))
}

func (this *HeaderWare) next(ctx *context.Context) bool {
		return true
}

func (this *HeaderWare) before(ctx *context.Context) bool {
		this.setId(ctx)
		return true
}
