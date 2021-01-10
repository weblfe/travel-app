package repositories

import (
		"github.com/astaxie/beego"
		"github.com/astaxie/beego/logs"
		"github.com/weblfe/travel-app/common"
	"github.com/weblfe/travel-app/libs"
	"github.com/weblfe/travel-app/models"
		"github.com/weblfe/travel-app/services"
		"github.com/weblfe/travel-app/transforms"
		"github.com/weblfe/travel-app/transports"
		"time"
)

type UserRegisterRepository interface {
		Register() common.ResponseJson
		RegisterByQuick() (*models.User, error)
}

type UserRegisterRepositoryImpl struct {
		ctx         common.BaseRequestContext
		userService services.UserService
		smsService  services.SmsCodeService
}

const (
		RegisterByAccount = "account"
		RegisterByMobile  = "mobile"
		RegisterByEmail   = "email"
		RegisterByThird   = "thirdParty"
)

func NewUserRegisterRepository(ctx common.BaseRequestContext) UserRegisterRepository {
		var repository = new(UserRegisterRepositoryImpl)
		repository.ctx = ctx
		repository.init()
		return repository
}

// 初始化
func (this *UserRegisterRepositoryImpl) init() {
		this.smsService = services.SmsCodeServiceOf()
		this.userService = services.UserServiceOf()
}

// username + password
// mobile + sms-code
// email + email-code
// third-party => { auth_code, user_info }
// 注册逻辑
func (this *UserRegisterRepositoryImpl) Register() common.ResponseJson {
		var request = new(transports.RegisterRequest)
		request.Load(this.ctx.GetInput())
		if request.Type == "" {
				request.Type = this.choose(request)
		}
		if request.Type == "" {
				return common.NewInvalidParametersResp(common.NewErrors(common.InvalidParametersCode, "请求异常"))
		}
		switch request.Type {
		case RegisterByAccount:
				return this.registerAccount(request.Account, request.Password, request)
		case RegisterByMobile:
				return this.registerByMobile(request.Mobile, request.Code, request)
		case RegisterByEmail:
				return this.registerByEmail(request.Email, request.Code, request)
		case RegisterByThird:
				return this.registerThirdParty(request, request.Third)
		}
		return common.NewResponse(common.RegisterFail, common.RegisterFailTip, common.NewErrors(common.RegisterFail, "未知注册方式 :"+request.Type))
}

// 自动选择注册方式
func (this *UserRegisterRepositoryImpl) choose(request *transports.RegisterRequest) string {
		if request.Mobile != "" && request.Code != "" {
				return RegisterByMobile
		}
		if request.Email != "" && request.Code != "" {
				return RegisterByEmail
		}
		if request.Account != "" && request.Password != "" {
				return RegisterByAccount
		}
		if request.Mobile != "" && request.Password != "" {
				request.Account = request.Mobile
				return RegisterByAccount
		}
		if request.Third != "" {
				return RegisterByThird
		}
		return ""
}

// 用户账号+密码登录
func (this *UserRegisterRepositoryImpl) registerAccount(account string, password string, request *transports.RegisterRequest) common.ResponseJson {
		if account == "" {
				return common.NewInvalidParametersResp(common.NewErrors(common.InvalidParametersCode, "用户账号不能为空"))
		}
		if password == "" {
				return common.NewInvalidParametersResp(common.NewErrors(common.InvalidParametersCode, "用户密码不能为空"))
		}
		var (
				user = new(models.User)
				data = beego.M{"username": account, "passwordHash": password, "register_way": "account"}
		)
		user.Load(data).Defaults()
		user.ResetPasswordTimes++
		if request.Mobile != "" {
				user.Mobile = request.Mobile
		}
		if this.userService.GetByUserName(user.UserName) != nil || this.userService.GetByMobile(user.Mobile) != nil {
				return common.NewInvalidParametersResp(common.NewErrors(common.InvalidParametersCode, "用户账号已注册"))
		}
		err := this.getUserService().Create(user)
		if err == nil {
				return common.NewSuccessResp(beego.M{"user": user}, "注册成功")
		}
		logs.Error(err)
		return common.NewResponse(common.RegisterFail, common.RegisterFailTip, common.NewErrors(err, common.RegisterFail))
}

// 手机号注册
func (this *UserRegisterRepositoryImpl) registerByEmail(email string, code string, request *transports.RegisterRequest) common.ResponseJson {
		return common.NewResponse(common.RegisterFail, common.RegisterFailTip)
}

// 手机号注册
func (this *UserRegisterRepositoryImpl) registerByMobile(mobile string, code string, request *transports.RegisterRequest) common.ResponseJson {
		if !this.smsService.Verify(mobile, code, "register") {
				return common.NewResponse(common.RegisterFail, common.RegisterFailTip, common.NewErrors(common.VerifyNotMatch, "验证码错误"))
		}
		var (
				user = models.NewUser()
				data = beego.M{
						"mobile":      mobile,
						"nickname":    "travelGo_" + libs.RandNumbers(8),
						"registerWay": "mobile",
						"username":    mobile,
						"created_at":  time.Now(),
				}
		)
		if request.Password != "" {
				data["passwordHash"] = request.Password
				data["resetPasswordTimes"] = 1
		}
		if request.Way != "" {
				data["registerWay"] = request.Way
		}
		// 创建用户
		user.Load(data)
		if u := this.userService.GetByMobile(mobile); u != nil {
				return common.NewResponse(common.RegisterFail, common.RegisterFailTip, common.NewErrors(common.RegisterFail, "手机号已注册"))
		}
		if err := this.userService.Create(user.Defaults()); err == nil {
				return common.NewSuccessResp(beego.M{"user": user.M(transforms.FilterUser)}, "注册成功")
		}
		return common.NewResponse(common.RegisterFail, common.RegisterFailTip)
}

// 第三方注册
func (this *UserRegisterRepositoryImpl) registerThirdParty(request *transports.RegisterRequest, third string) common.ResponseJson {
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

func (this *UserRegisterRepositoryImpl) RegisterByQuick() (*models.User, error) {
		var (
				data = new(transports.QuickRegister)
				user = models.NewUser()
				info = data.Load(this.ctx.GetInput()).M()
		)
		if len(info) == 0 {
				return nil, common.NewErrors(common.ServiceFailed, "参数不足")
		}
		err := this.userService.Create(user.Load(info).Defaults())
		if err == nil {
				return user, nil
		}
		return nil, err
}

