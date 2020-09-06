package middlewares

import (
		"github.com/astaxie/beego/context"
		"github.com/astaxie/beego/logs"
		"github.com/weblfe/travel-app/common"
		"github.com/weblfe/travel-app/libs"
		"github.com/weblfe/travel-app/plugins"
		"strings"
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
		return this.limiter(ctx)
}

func (this *HeaderWare) limiter(ctx *context.Context) bool {
		var (
				ctxLimiter = this.createLimiterCtx(ctx)
				handler    = plugins.GetLimiter().New(ctxLimiter)
				result     = handler(ctxLimiter)
		)
		if result.Ok {
				return true
		}
		res := common.NewUnLoginResp(common.NewErrors(common.LimitCode, common.LimitError))
		err := ctx.Output.JSON(res, this.hasIndex(), true)
		if err != nil {
				logs.Error(err)
		}
		return false
}

func (this *HeaderWare) createLimiterCtx(ctx *context.Context) plugins.ContextLimit {
		var ctxLimiter = plugins.NewContextLimit()
		ctxLimiter = ctxLimiter.SetPolicy(plugins.LimiterTokenPolicy).
				SetValue(plugins.TokenCtxValueKey, this.getToken(ctx)).
				SetValue(plugins.MacCtxValueKey, this.getMac(ctx))
		return ctxLimiter
}

func (this *HeaderWare) getToken(ctx *context.Context) string {
		var token = ctx.Request.Header.Get(AppAccessTokenHeader)
		if token == "" {
				return ""
		}
		var _, value = AuthorizationParse(token)
		return value
}

func (this *HeaderWare) getMac(ctx *context.Context) string {
		var mac = ctx.Request.Header.Get(AppAccessMacHeader)
		return mac
}

func AuthorizationParse(token string) (string, string) {
		if token == "" || !strings.Contains(token, " ") {
				return "", token
		}
		var arr = strings.SplitN(token, " ", 2)
		if len(arr) == 2 {
				return strings.TrimSpace(strings.ToLower(arr[0])), strings.TrimSpace(arr[1])
		}
		if strings.Contains(token, "Bearer ") {
				return "bearer", strings.Replace(token, "Bearer ", "", 1)
		}
		if strings.Contains(token, "bearer ") {
				return "bearer", strings.Replace(token, "bearer ", "", 1)
		}
		if strings.Contains(token, "Basic ") {
				return "basic", strings.Replace(token, "Basic ", "", 1)
		}
		if strings.Contains(token, "basic ") {
				return "basic", strings.Replace(token, "basic ", "", 1)
		}
		return "", token
}
