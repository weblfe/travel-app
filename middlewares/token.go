package middlewares

import (
		"github.com/astaxie/beego/context"
		"github.com/weblfe/travel-app/models"
		"github.com/weblfe/travel-app/services"
)

type tokenMiddleware struct {
		middlewareImpl
}

var (
		tokenInstance *tokenMiddleware
)

const (
		AppAccessTokenHeader = "authorization"
		AuthUser             = "user"
		AuthUserId           = "userId"
)

func newToken() {
		tokenInstance = new(tokenMiddleware)
		tokenInstance.Init()
}

func GetTokenMiddleware() *tokenMiddleware {
		if tokenInstance != nil {
				return tokenInstance
		}
		getLock("token").Do(newToken)
		return tokenInstance
}

func (this *tokenMiddleware) Init() {
		this.Name = "token"
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
		if err == nil || user == nil {
				return false
		}
		if err := ctx.Input.CruSession.Set(AuthUser, user); err == nil {
				_ = ctx.Input.CruSession.Set(AuthUserId, user.Id.Hex())
				this.dispatch(token, user)
				return true
		}
		return false
}

// 记录当前登录用户
func (this *tokenMiddleware) dispatch(token string, user *models.User) {
		// @todo
}
