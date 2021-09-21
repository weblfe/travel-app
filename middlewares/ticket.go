package middlewares

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/context"
	"github.com/astaxie/beego/logs"
	"github.com/weblfe/travel-app/common"
	"github.com/weblfe/travel-app/services"
	"os"
	"strconv"
	"strings"
)

type AttachTicketMiddleware struct {
	middlewareImpl
	allowNoAuth int
}

var (
	attachTicketMiddlewareInstance *AttachTicketMiddleware
)

const (
	AttachTicket               = "ticket"
	AttachTicketMiddlewareName = "attach.ticket"
	TicketOk                   = "_ticketPass"
	LoginOk                    = "_loginPass"
)

// GetAttachTicketMiddleware 获取附件防盗链中间键
func GetAttachTicketMiddleware() *AttachTicketMiddleware {
	if attachTicketMiddlewareInstance == nil {
		getLock(AttachTicket).Do(newAttachTicketWare)
	}
	return attachTicketMiddlewareInstance
}

func newAttachTicketWare() {
	attachTicketMiddlewareInstance = new(AttachTicketMiddleware)
	attachTicketMiddlewareInstance.allowNoAuth = -1
	attachTicketMiddlewareInstance.Init()
}

func (this *AttachTicketMiddleware) Init() {
	this.Name = AttachTicketMiddlewareName
	this.SetHandler(this.verify, this.auth, this.forbid)
	Register(this.Middleware(), this)
}

// 放行令牌验证
func (this *AttachTicketMiddleware) verify(ctx *context.Context) bool {
	var (
		ticket string
		_      = ctx.Input.Bind(&ticket, "ticket")
	)
	ctx.Input.SetParam(TicketOk, "")
	// @todo 验证令牌
	if ticket != "" {
		ctx.Input.SetParam(TicketOk, "1")
		return true
	}
	var arr = strings.SplitN(ctx.Request.URL.Path, "/", -1)
	if len(arr) > 0 {
		if services.UrlTicketServiceOf().Expired(arr[len(arr)-1]) {
			ctx.Input.SetParam(TicketOk, "1")
		}
	}
	return true
}

// 登陆用户验证
func (this *AttachTicketMiddleware) auth(ctx *context.Context) bool {
	v := ctx.Input.Session(AuthUserId)
	ctx.Input.SetParam(LoginOk, "")
	if v == nil {
		GetTokenMiddleware().Filter(ctx)
		v = ctx.Input.Session(AuthUserId)
		if v == nil {
			return true
		}
	}
	userId := v.(string)
	if userId != "" {
		ctx.Input.SetParam(LoginOk, "1")
	}
	return true
}

// 是否禁止
func (this *AttachTicketMiddleware) forbid(ctx *context.Context) bool {
	var (
		v        = ctx.Input.Param(LoginOk)
		pass     = ctx.Input.Param(TicketOk)
		hasIndex = beego.BConfig.RunMode != beego.PROD
	)
	if v == "1" || pass == "1" || this.allows() {
		return true
	}
	res := common.NewUnLoginResp(common.NewErrors(common.PermissionCode, "权限不足无法访问!"))
	err := ctx.Output.JSON(res, hasIndex, true)
	if err != nil {
		logs.Error(err)
	}
	return false
}

func (this *AttachTicketMiddleware) allows() bool {
	if this.allowNoAuth != -1 {
		return this.allowNoAuth > 0
	}
	var read = os.Getenv("ALLOW_ATTACH_PUBLIC_READ")
	if read == "" {
		this.allowNoAuth = 0
		return false
	}
	if b, err := strconv.ParseBool(read); err == nil {
		if b {
			this.allowNoAuth = 1
		} else {
			this.allowNoAuth = 0
		}
		return b
	}
	if n, err := strconv.Atoi(read); err == nil {
		if n > 0 {
			this.allowNoAuth = 1
		} else {
			this.allowNoAuth = 0
		}
	}
	return this.allowNoAuth > 0
}
