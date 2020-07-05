package repositories

import (
		"encoding/json"
		"github.com/astaxie/beego"
		"github.com/astaxie/beego/context"
		"github.com/weblfe/travel-app/common"
		"github.com/weblfe/travel-app/models"
		"github.com/weblfe/travel-app/services"
		"time"
)

type UserRegisterRepository interface {
		Register() common.ResponseJson
		RegisterByQuick() (*models.User, error)
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
		var (
				mobile   string
				password string
				username string
				email    string
				code     string
				ctx      = this.ctx
				typ      = ctx.GetString("type")
		)
		if typ == "" {
				if typ = this.choose(ctx); typ == "" {
						return common.NewInvalidParametersResp(common.NewErrors(common.InvalidParametersCode, "请求异常"))
				}
		}
		switch typ {
		case RegisterByAccount:
				_ = ctx.Ctx.Input.Bind(&password, "password")
				_ = ctx.Ctx.Input.Bind(&username, "username")
				return this.registerAccount(username, password, ctx)
		case RegisterByMobile:
				_ = ctx.Ctx.Input.Bind(&mobile, "mobile")
				_ = ctx.Ctx.Input.Bind(&code, "code")
				return this.registerByMobile(mobile, code, ctx)
		case RegisterByEmail:
				_ = ctx.Ctx.Input.Bind(&email, "email")
				_ = ctx.Ctx.Input.Bind(&code, "code")
				return this.registerByEmail(ctx.GetString("email"), ctx.GetString("code"), ctx)
		case RegisterByThird:
				return this.registerThirdParty(ctx)
		}
		return common.NewResponse(common.RegisterFail, common.RegisterFailTip, common.NewErrors(common.RegisterFail, "未知注册方式 :"+typ))
}

// 自动选择注册方式
func (this *UserRegisterRepositoryImpl) choose(ctx *beego.Controller) string {
		var (
				code     string
				mobile   string
				account  string
				password string
				email    string
				third    = ctx.GetString(RegisterByThird)
		)
		_ = ctx.Ctx.Input.Bind(&code, "code")
		_ = ctx.Ctx.Input.Bind(&mobile, "mobile")
		_ = ctx.Ctx.Input.Bind(&account, "username")
		_ = ctx.Ctx.Input.Bind(&password, "password")
		_ = ctx.Ctx.Input.Bind(&email, "email")
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
				data = beego.M{"username": account, "passwordHash": password, "register_way": "account"}
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
		if !this.smsService.Verify(mobile, code, "register") {
				return common.NewResponse(common.RegisterFail, common.RegisterFailTip, common.NewErrors(common.VerifyNotMatch, "验证码错误"))
		}
		var (
				user     = models.NewUser()
				way      string
				password string
				_        = ctx.Ctx.Input.Bind(&way, "_register")
				_        = ctx.Ctx.Input.Bind(&password, "password")
				data     = beego.M{
						"mobile":      mobile,
						"nickname":    "nick_" + mobile,
						"registerWay": "mobile",
						"username":    mobile,
						"created_at":  time.Now(),
				}
		)
		if password != "" {
				data["passwordHash"] = password
				data["resetPasswordTimes"] = 1
		}
		if way != "" {
				data["registerWay"] = way
		}
		// 创建用户
		user.Load(data)
		if u := this.userService.GetByMobile(mobile); u != nil {
				return common.NewResponse(common.RegisterFail, common.RegisterFailTip, common.NewErrors(common.RegisterFail, "手机号已注册"))
		}
		if err := this.userService.Create(user.Defaults()); err == nil {
				return common.NewSuccessResp(beego.M{"user": user.M(filterUser)}, "注册成功")
		}
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

func (this *UserRegisterRepositoryImpl) RegisterByQuick() (*models.User, error) {
		var (
				request = this.ctx
				data    = new(QuickRegister)
				user    = models.NewUser()
				info    = data.Load(request.Ctx.Input).M()
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

// 快捷注册
type QuickRegister struct {
		Mobile      string `json:"mobile"`
		Password    string `json:"password"`
		RegisterWay string `json:"_register"`
		Gender      int    `json:"gender"`
		UserName    string `json:"username"`
		NickName    string `json:"nickname"`
		Email       string `json:"email"`
		AvatarId    string `json:"avatarId"`
}

func (this *QuickRegister)Load(ctx *context.BeegoInput) *QuickRegister  {
		var (
			_       = json.Unmarshal(ctx.RequestBody, this)
			mapper  = map[string]interface{}{
					"mobile":&this.Mobile,
					"password":&this.Password,
					"_register":&this.RegisterWay,
					"gender":&this.Gender,
					"username":&this.UserName,
					"nickname":&this.NickName,
					"email":&this.Email,
					"avatarId":&this.AvatarId,
			}
		)
		if this.Mobile == ""  {
				for key,v:=range mapper {
						_ = ctx.Bind(v,key)
				}
		}
		return this
}

func (this *QuickRegister) M(filters ...func(m beego.M) beego.M) beego.M {
		if this.RegisterWay == "" {
				this.RegisterWay = "quick"
		}
		var data = beego.M{
				"mobile":       this.Mobile,
				"passwordHash": this.Password,
				"registerWay":  this.RegisterWay,
				"gender":       this.Gender,
				"username":     this.UserName,
				"nickname":     this.NickName,
				"email":        this.Email,
				"avatarId":     this.AvatarId,
		}
		filters = append(filters, filterEmpty)
		for _, filter := range filters {
				data = filter(data)
		}
		return data
}
