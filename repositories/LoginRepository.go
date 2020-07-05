package repositories

import (
		"encoding/json"
		"github.com/astaxie/beego"
		"github.com/astaxie/beego/context"
		"github.com/weblfe/travel-app/common"
		"github.com/weblfe/travel-app/libs"
		"github.com/weblfe/travel-app/services"
)

type LoginRepository interface {
		Login() common.ResponseJson
}

type LoginRepositoryImpl struct {
		ctx         *beego.Controller
		userService services.UserService
		smsService  services.SmsCodeService
		authService services.AuthService
}

func NewLoginRepository(ctx *beego.Controller) LoginRepository {
		var repository = new(LoginRepositoryImpl)
		repository.ctx = ctx
		repository.init()
		return repository
}

func (this *LoginRepositoryImpl) init() {
		this.userService = services.UserServiceOf()
		this.smsService = services.SmsCodeServiceOf()
		this.authService = services.AuthServiceOf()
}

func (this *LoginRepositoryImpl) Login() common.ResponseJson {
		var (
				typ     string
				request = new(Login)
				input   = this.ctx.Ctx.Input
		)
		request.Load(input)
		typ = this.choose(request.Mobile, request.Code, request.Username, request.Password, request.Email)
		if typ == "" {
				return common.NewInvalidParametersResp(common.NewErrors(1020, "参数不足无法登陆"))
		}
		switch typ {
		case "mobile":
				return this.loginByMobile(request.Mobile, request.Code)
		case "mobile_password":
				return this.loginByMobilePassword(request.Mobile, request.Password)
		case "account":
				return this.loginByAccountPassword(request.Username, request.Password)
		case "email":
				return this.loginByEmailCode(request.Email, request.Code)
		case "email_password":
				return this.loginByEmail(request.Email, request.Password)
		}
		return common.NewInvalidParametersResp(common.NewSuccessResp(1023, "未知登陆请求,登陆失败"))
}

func (this *LoginRepositoryImpl) choose(mobile, code, username, password, email string) string {
		if mobile != "" && code != "" {
				return "mobile"
		}
		if mobile != "" && password != "" {
				return "mobile_password"
		}
		if username != "" && password != "" {
				return "account"
		}
		if email != "" && password != "" {
				return "email_password"
		}
		if email != "" && code != "" {
				return "email"
		}
		return ""
}

func (this *LoginRepositoryImpl) loginByMobile(mobile string, code string) common.ResponseJson {
		if !libs.IsCnMobile(mobile) {
				return common.NewErrorResp(common.NewErrors(common.VerifyNotMatch, "手机号格式不正确"))
		}
		var (
				err         error
				smsCodeType = "login"
		)
		// 短信验证
		if !this.smsService.Verify(mobile, code, smsCodeType) {
				smsCodeType = "quick"
				if !this.smsService.Verify(mobile, code, smsCodeType) {
						return common.NewErrorResp(common.NewErrors(common.VerifyNotMatch, "验证码不匹配"))
				}
		}
		user := this.userService.GetByMobile(mobile)
		if user == nil && smsCodeType == "login" {
				return common.NewErrorResp(common.NewErrors(1021, "用户不存在"))
		}
		if smsCodeType == "quick" && user == nil {
				user, err = NewUserRegisterRepository(this.ctx).RegisterByQuick()
				if user == nil {
						return common.NewErrorResp(common.NewErrors(1021, "用户不存在"))
				}
				if err != nil {
						return common.NewErrorResp(err.(common.Errors), "快捷登陆失败")
				}
		}
		if isForbid(user) {
				return common.NewErrorResp(common.NewErrors(common.PermissionCode, "账号禁用状态"))
		}
		token := this.authService.Token(user)
		if token == "" {
				return common.NewErrorResp(common.NewErrors(common.ServiceFailed, "登陆服务异常"))
		}
		data := user.M(filterUser)
		delete(data, "access_tokens")
		return common.NewSuccessResp(beego.M{"user": data, "token": token}, "登陆成功")
}

func (this *LoginRepositoryImpl) loginByAccountPassword(account string, password string) common.ResponseJson {
		user := this.userService.GetByUserName(account)
		if user == nil {
				return common.NewErrorResp(common.NewErrors(1021, "用户不存在"))
		}
		if isForbid(user) {
				return common.NewErrorResp(common.NewErrors(common.PermissionCode, "账号禁用状态"))
		}
		if !libs.PasswordVerify(user.PasswordHash, password) {
				return common.NewErrorResp(common.NewErrors(common.VerifyNotMatch, "密码不正确"))
		}
		token := this.authService.Token(user)
		data := user.M(filterUserBase)
		return common.NewSuccessResp(beego.M{"user": data, "token": token}, "登陆成功")
}

func (this *LoginRepositoryImpl) loginByMobilePassword(mobile string, password string) common.ResponseJson {
		user := this.userService.GetByMobile(mobile)
		if user == nil {
				return common.NewErrorResp(common.NewErrors(1021, "用户不存在"))
		}
		if isForbid(user) {
				return common.NewErrorResp(common.NewErrors(common.PermissionCode, "账号禁用状态"))
		}
		if !libs.PasswordVerify(user.PasswordHash, password) {
				return common.NewErrorResp(common.NewErrors(common.VerifyNotMatch, "密码不正确"))
		}
		token := this.authService.Token(user)
		data := user.M(filterUserBase)
		return common.NewSuccessResp(beego.M{"user": data, "token": token}, "登陆成功")
}

func (this *LoginRepositoryImpl) loginByEmailCode(email string, code string) common.ResponseJson {
		user := this.userService.GetByEmail(email)
		if user == nil {
				return common.NewErrorResp(common.NewErrors(1021, "用户不存在"))
		}
		if isForbid(user) {
				return common.NewErrorResp(common.NewErrors(common.PermissionCode, "账号禁用状态"))
		}
		// @todo 邮件验证码
		if code == "" {
				return common.NewErrorResp(common.NewErrors(common.VerifyNotMatch, "验证码不正确"))
		}
		token := this.authService.Token(user)
		data := user.M(filterUserBase)
		return common.NewSuccessResp(beego.M{"user": data, "token": token}, "登陆成功")
}

func (this *LoginRepositoryImpl) loginByEmail(email string, password string) common.ResponseJson {
		user := this.userService.GetByEmail(email)
		if user == nil {
				return common.NewErrorResp(common.NewErrors(1021, "用户不存在"))
		}
		if isForbid(user) {
				return common.NewErrorResp(common.NewErrors(common.PermissionCode, "账号禁用状态"))
		}
		if !libs.PasswordVerify(user.PasswordHash, password) {
				return common.NewErrorResp(common.NewErrors(common.VerifyNotMatch, "验证码不正确"))
		}
		token := this.authService.Token(user)
		data := user.M(filterUserBase)
		return common.NewSuccessResp(beego.M{"user": data, "token": token}, "登陆成功")
}

// 登陆参数
type Login struct {
		Mobile   string `json:"mobile"`
		Code     string `json:"code"`
		Username string `json:"username"`
		Password string `json:"password"`
		Email    string `json:"email"`
}

func (this *Login) Load(ctx *context.BeegoInput) *Login {
		var (
				_      = json.Unmarshal(ctx.RequestBody, this)
				mapper = map[string]interface{}{
						"mobile":   &this.Mobile,
						"code":     &this.Code,
						"username": &this.Username,
						"email":    &this.Email,
						"password": &this.Password,
				}
		)
		if this.Email == "" && this.Username == "" && this.Mobile == "" {
				for key, addr := range mapper {
						_ = ctx.Bind(addr, key)
				}
		}
		return this
}

func (this *Login) M(filters ...func(m beego.M) beego.M) beego.M {
		var (
				data = beego.M{
						"mobile":   this.Mobile,
						"code":     this.Code,
						"username": this.Username,
						"password": this.Password,
						"email":    this.Email,
				}
		)
		filters = append(filters, filterEmpty)
		for _, filter := range filters {
				data = filter(data)
		}
		return data
}
