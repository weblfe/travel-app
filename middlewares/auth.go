package middlewares

import (
		"github.com/astaxie/beego"
		"github.com/astaxie/beego/context"
		"github.com/astaxie/beego/logs"
		"github.com/weblfe/travel-app/common"
)

type AuthMiddleware struct {
		tokenMiddleware
}

var (
		authMiddlewareInstance *AuthMiddleware
)

const (
		Auth               = "auth"
		AuthMiddlewareName = "api.auth"
)

func GetAuthMiddleware() *AuthMiddleware {
		if authMiddlewareInstance == nil {
				getLock(Auth).Do(newAuthMiddleware)
		}
		return authMiddlewareInstance
}

func newAuthMiddleware() {
		authMiddlewareInstance = new(AuthMiddleware)
		authMiddlewareInstance.Init()
}

func (this *AuthMiddleware) Init() {
		this.Name = AuthMiddlewareName
		this.SetHandler(this.auth, this.forbid)
		Register(this.Middleware(), this)
}

func (this *AuthMiddleware) auth(ctx *context.Context) bool {
		this.Filter(ctx)
		return true
}

func (this *AuthMiddleware) forbid(ctx *context.Context) bool {
		hasIndex := beego.BConfig.RunMode != beego.PROD
		v := ctx.Input.Session(AuthUserId)
		if v == nil {
				v = ""
		}
		userId := v.(string)
		if userId != "" {
				return true
		}
		res:=common.NewUnLoginResp(common.NewErrors(common.UnLoginCode,"请先登录!"))
		err := ctx.Output.JSON(res, hasIndex, true)
		if err != nil {
				logs.Error(err)
		}
		return false
}
