package middlewares

import (
		"fmt"
		"github.com/astaxie/beego"
		"github.com/astaxie/beego/context"
		"github.com/astaxie/beego/logs"
		"github.com/weblfe/travel-app/common"
		"github.com/weblfe/travel-app/services"
)

type roleMiddleware struct {
		middlewareImpl
}

var (
		roleInstance *roleMiddleware
)

const (
		RoleMiddleware = "role"
)

func newRoleWare() {
		roleInstance = new(roleMiddleware)
		roleInstance.Init()
}

func GetRoleMiddleware() *roleMiddleware {
		if roleInstance != nil {
				return roleInstance
		}
		getLock(RoleMiddleware).Do(newRoleWare)
		return roleInstance
}

func (this *roleMiddleware) Init() {
		this.Name = RoleMiddleware
		this.handler = this.filter
		Register(this.Middleware(), this)
}

func (this *roleMiddleware) filter(ctx *context.Context) {
		var (
				value      = ctx.Input.CruSession.Get(AuthUserId)
				hasIndex   = beego.BConfig.RunMode != beego.PROD
				unLogin    = common.NewUnLoginResp(common.NewErrors(common.UnLoginCode, "请先登录!"))
				permission = common.NewErrorResp(common.NewErrors(common.PermissionCode, common.PermissionError), "权限不足")
		)
		fmt.Println(ctx.Request.RequestURI)
		fmt.Println(ctx.Request.URL)
		// 未登陆
		if value == "" {
				err := ctx.Output.JSON(unLogin, hasIndex, true)
				if err != nil {
						logs.Error(err)
				}
				return
		}
		var user = services.UserServiceOf().GetById(value.(string))
		// 用户角色
		if user == nil || !user.IsRootRole() {
				err := ctx.Output.JSON(permission, hasIndex, true)
				if err != nil {
						logs.Error(err)
				}
				return
		}
}

func (this *roleMiddleware) Filter(ctx *context.Context) bool {
		var (
				value = ctx.Input.CruSession.Get(AuthUserId)
		)
		// 未登陆
		if value == "" {
				return false
		}
		var user = services.UserServiceOf().GetById(value.(string))
		// 用户角色
		if user == nil || !user.IsRootRole() {
				return false
		}
		return true
}
