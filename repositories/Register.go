package repositories

import (
		"github.com/astaxie/beego"
		"github.com/weblfe/travel-app/common"
		"github.com/weblfe/travel-app/models"
		"github.com/weblfe/travel-app/services"
)

type UserRegisterRepository interface {
		Register() common.ResponseJson
}

type UserRegisterRepositoryImpl struct {
		ctx         *beego.Controller
		userService services.UserService
		smsService  services.SmsCodeService
}

const (
		RegisterByAccount = "account"
		RegisterByMobile  = "mobile"
		RegisterByEmail   = "email"
		RegisterByThird   = "thirdParty"
)

func NewUserRegisterRepository(ctx *beego.Controller) UserRegisterRepository {
		var repository = new(UserRegisterRepositoryImpl)
		repository.ctx = ctx
		return repository
}

// username + password
// mobile + sms-code
// email + email-code
// third-party => { auth_code, user_info }
// 注册逻辑
func (this *UserRegisterRepositoryImpl) Register() common.ResponseJson {
		var (
				ctx = this.ctx
				typ = ctx.GetString("type")
		)
		if typ == "" {
				if typ = this.choose(ctx); typ == "" {
						return common.NewInvalidParametersResp(common.NewErrors(common.InvalidParametersCode, "请求异常"))
				}
		}
		switch typ {
		case RegisterByAccount:
				return this.registerAccount(ctx.GetString("username"), ctx.GetString("password"), ctx)
		case RegisterByMobile:
				return this.registerByMobile(ctx.GetString("mobile"), ctx.GetString("code"), ctx)
		case RegisterByEmail:
				return this.registerByEmail(ctx.GetString("email"), ctx.GetString("code"), ctx)
		case RegisterByThird:
				return this.registerThirdParty(ctx)
		}
		return common.NewResponse(common.RegisterFail, common.RegisterFailTip, common.NewErrors(common.RegisterFail, "未知注册方式 :"+typ))
}

// 自动选择注册方式
func (this *UserRegisterRepositoryImpl) choose(ctx *beego.Controller) string {
		var (
				code     = ctx.GetString("code")
				mobile   = ctx.GetString("mobile")
				account  = ctx.GetString("username")
				password = ctx.GetString("password")
				email    = ctx.GetString("email")
				third    = ctx.GetString(RegisterByThird)
		)
		if mobile != "" && code != "" {
				return RegisterByMobile
		}
		if email != "" && code != "" {
				return RegisterByEmail
		}
		if account != "" && password != "" {
				return RegisterByAccount
		}
		if third != "" {
				return RegisterByThird
		}
		return ""
}

// 用户账号+密码登录
func (this *UserRegisterRepositoryImpl) registerAccount(account string, password string, ctx *beego.Controller) common.ResponseJson {
		if account == "" {
				return common.NewInvalidParametersResp(common.NewErrors(common.InvalidParametersCode, "用户账号不能为空"))
		}
		if password == "" {
				return common.NewInvalidParametersResp(common.NewErrors(common.InvalidParametersCode, "用户密码不能为空"))
		}
		var (
				user = new(models.User)
				data = beego.M{"username": account, "password": password}
		)
		err := this.getUserService().Create(user.Load(data).Defaults())
		if err == nil {
				return common.NewSuccessResp(beego.M{"user": user}, "注册成功")
		}
		return common.NewResponse(common.RegisterFail, common.RegisterFailTip)
}

// 手机号注册
func (this *UserRegisterRepositoryImpl) registerByEmail(email string, code string, ctx *beego.Controller) common.ResponseJson {
		return common.NewResponse(common.RegisterFail, common.RegisterFailTip)
}

// 手机号注册
func (this *UserRegisterRepositoryImpl) registerByMobile(mobile string, code string, ctx *beego.Controller) common.ResponseJson {
		return common.NewResponse(common.RegisterFail, common.RegisterFailTip)
}

// 第三方注册
func (this *UserRegisterRepositoryImpl) registerThirdParty(ctx *beego.Controller) common.ResponseJson {
		return common.NewResponse(common.RegisterFail, common.RegisterFailTip)
}

func (this *UserRegisterRepositoryImpl) getUserService() services.UserService {
		if this.userService == nil {
				this.userService = services.UserServiceOf()
		}
		return this.userService
}

func (this *UserRegisterRepositoryImpl) getSmsService() services.SmsCodeService {
		if this.smsService == nil {
				this.smsService = services.SmsCodeServiceOf()
		}
		return this.smsService
}

func (this *UserRegisterRepositoryImpl) getEmailService() services.SmsCodeService {
		if this.smsService == nil {
				this.smsService = services.SmsCodeServiceOf()
		}
		return this.smsService
}