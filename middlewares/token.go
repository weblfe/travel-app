package middlewares

import (
		"github.com/astaxie/beego/context"
		"github.com/weblfe/travel-app/services"
)

type tokenMiddleware struct {
		middlewareImpl
}

var (
		tokenInstance *tokenMiddleware
)

const (
		AuthUser             = "user"
		AuthUserId           = "userId"
		TokenMiddleware      = "token"
		AppAccessTokenHeader = "authorization"
		AppAccessMacHeader   = "visitor"
)

func newToken() {
		tokenInstance = new(tokenMiddleware)
		tokenInstance.Init()
}

func GetTokenMiddleware() *tokenMiddleware {
		if tokenInstance != nil {
				return tokenInstance
		}
		getLock(TokenMiddleware).Do(newToken)
		return tokenInstance
}

func (this *tokenMiddleware) Init() {
		this.Name = TokenMiddleware
		this.handler = this.filter
		Register(this.Middleware(), this)
}

func (this *tokenMiddleware) filter(ctx *context.Context) {
		this.clear(ctx)
		token := ctx.Request.Header.Get(AppAccessTokenHeader)
		if token != "" && this.initSessionByToken(token, ctx) {
				return
		}
		if cookie, err := ctx.Request.Cookie(AppAccessTokenHeader); err == nil {
				token = cookie.Value
				this.initSessionByToken(token, ctx)
		}
}

func (this *tokenMiddleware) clear(ctx *context.Context) {
		_ = ctx.Input.CruSession.Delete(AuthUser)
		_ = ctx.Input.CruSession.Delete(AuthUserId)
		_ = ctx.Input.CruSession.Delete(AppAccessTokenHeader)
		_ = ctx.Input.CruSession.Delete(AppAccessMacHeader)
		ctx.Input.SetParam("_userId", "")
}

func (this *tokenMiddleware) Filter(ctx *context.Context) bool {
		this.clear(ctx)
		token := ctx.Request.Header.Get(AppAccessTokenHeader)
		if token != "" && this.initSessionByToken(token, ctx) {
				return true
		}
		if cookie, err := ctx.Request.Cookie(AppAccessTokenHeader); err == nil {
				token = cookie.Value
				return this.initSessionByToken(token, ctx)
		}
		return false
}

func (this *tokenMiddleware) initSessionByToken(token string, ctx *context.Context) bool {
		user, err := services.AuthServiceOf().GetByAccessToken(token)
		if err != nil || user == nil {
				return false
		}
		if err := ctx.Input.CruSession.Set(AuthUser, user.M()); err == nil {
				uid := user.Id.Hex()
				_ = ctx.Input.CruSession.Set(AuthUserId, uid)
				_ = ctx.Input.CruSession.Set(AppAccessTokenHeader, token)
				ctx.Input.SetParam("_userId", uid)
				return true
		}
		return false
}
