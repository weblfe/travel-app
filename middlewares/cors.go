package middlewares

import (
		"github.com/astaxie/beego"
		"github.com/astaxie/beego/context"
		"github.com/astaxie/beego/plugins/cors"
)

type CorsMiddleware struct {
		middlewareImpl
		handler beego.FilterFunc
}

var (
		corsMiddlewareInstance *CorsMiddleware
)

const (
		Cors               = "cors"
		CorsMiddlewareName = "api.cors"
)

// GetCorsMiddleware 垮域中间件
func GetCorsMiddleware() *CorsMiddleware {
		if corsMiddlewareInstance == nil {
				getLock(Cors).Do(newCorsMiddleware)
		}
		return corsMiddlewareInstance
}

func newCorsMiddleware() {
		corsMiddlewareInstance = new(CorsMiddleware)
		corsMiddlewareInstance.Init()
}

func (this *CorsMiddleware) Init() {
		this.Name = CorsMiddlewareName
		this.SetHandler(this.cors)
		Register(this.Middleware(), this)
}

// 跨域options
func (this *CorsMiddleware) cors(ctx *context.Context) bool {
		var handler = this.getHandler(ctx)
		if handler == nil {
				return false
		}
		handler(ctx)
		return true
}

// 获取处理器
func (this *CorsMiddleware) getHandler(ctx *context.Context) beego.FilterFunc {
		if ctx == nil {
				return nil
		}
		if this.handler != nil {
				return this.handler
		}
		this.handler = cors.Allow(&cors.Options{
				//允许访问所有源
				AllowAllOrigins: true,
				//可选参数"GET", "POST", "PUT", "DELETE", "OPTIONS" (*为所有)
				//其中Options跨域复杂请求预检
				AllowMethods: []string{"*"},
				//指的是允许的Header的种类
				AllowHeaders: []string{"*"},
				//公开的HTTP标头列表
				ExposeHeaders: []string{"Content-Length", "Content-Type"},
				//如果设置，则允许共享身份验证凭据，例如cookie
				AllowCredentials: true,
		})
		return this.handler
}
