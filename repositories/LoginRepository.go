package repositories

import (
		"github.com/astaxie/beego"
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
				mobile   string
				code     string
				username string
				password string
				email    string
				typ      string
				input    = this.ctx.Ctx.Input
		)
		_ = input.Bind(&mobile, "mobile")
		_ = input.Bind(&code, "code")
		_ = input.Bind(&password, "password")
		_ = input.Bind(&username, "username")
		_ = input.Bind(&email, "email")
		typ = this.choose(mobile, code, username, password, email)
		if typ == "" {
				return common.NewInvalidParametersResp(common.NewErrors(1020, "参数不足无法登陆"))
		}
		switch typ {
		case "mobile":
				return this.loginByMobile(mobile, code)
		case "account":
				return this.loginByAccountPassword(username,password)
		case "email":
				return this.loginByEmailCode(email,code)
		case "email_password":
				return this.loginByEmail(email,password)
		}
		return common.NewInvalidParametersResp(common.NewSuccessResp(1023,"未知登陆请求,登陆失败"))
}

func (this *LoginRepositoryImpl) choose(mobile, code, username, password, email string) string {
		if mobile != "" && code != "" {
				return "mobile"
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
		if !this.smsService.Verify(mobile, code, "login") {
				return common.NewErrorResp(common.NewErrors(common.VerifyNotMatch, "验证码不匹配"))
		}
		user := this.userService.GetByMobile(mobile)
		if user == nil {
				return common.NewErrorResp(common.NewErrors(1021, "用户不存在"))
		}
		if isForbid(user) {
				return common.NewErrorResp(common.NewErrors(common.PermissionCode,"账号禁用状态"))
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
				return common.NewErrorResp(common.NewErrors(common.PermissionCode,"账号禁用状态"))
		}
		if !libs.PasswordVerify(user.Password, password) {
				return common.NewErrorResp(common.NewErrors(common.VerifyNotMatch, "密码不正确"))
		}
		token := this.authService.Token(user)
		data := user.M(filterUser)
		delete(data, "access_tokens")
		return common.NewSuccessResp(beego.M{"user": data, "token": token}, "登陆成功")
}

func (this *LoginRepositoryImpl) loginByEmailCode(email string, code string) common.ResponseJson {
		user := this.userService.GetByEmail(email)
		if user == nil {
				return common.NewErrorResp(common.NewErrors(1021, "用户不存在"))
		}
		if isForbid(user) {
				return common.NewErrorResp(common.NewErrors(common.PermissionCode,"账号禁用状态"))
		}
		// @todo 邮件验证码
		if code == "" {
				return common.NewErrorResp(common.NewErrors(common.VerifyNotMatch, "验证码不正确"))
		}
		token := this.authService.Token(user)
		data := user.M(filterUser)
		delete(data, "access_tokens")
		return common.NewSuccessResp(beego.M{"user": data, "token": token}, "登陆成功")
}

func (this *LoginRepositoryImpl) loginByEmail(email string, password string) common.ResponseJson {
		user := this.userService.GetByEmail(email)
		if user == nil {
				return common.NewErrorResp(common.NewErrors(1021, "用户不存在"))
		}
		if isForbid(user) {
				return common.NewErrorResp(common.NewErrors(common.PermissionCode,"账号禁用状态"))
		}
		if !libs.PasswordVerify(user.Password,password) {
				return common.NewErrorResp(common.NewErrors(common.VerifyNotMatch, "验证码不正确"))
		}
		token := this.authService.Token(user)
		data := user.M(filterUser)
		delete(data, "access_tokens")
		return common.NewSuccessResp(beego.M{"user": data, "token": token}, "登陆成功")
}