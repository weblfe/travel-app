package repositories

import (
		"github.com/astaxie/beego"
		"github.com/weblfe/travel-app/common"
		"github.com/weblfe/travel-app/libs"
		"github.com/weblfe/travel-app/middlewares"
		"github.com/weblfe/travel-app/models"
		"github.com/weblfe/travel-app/services"
		"github.com/weblfe/travel-app/transforms"
		"github.com/weblfe/travel-app/transports"
		"time"
)

type LoginRepository interface {
		Login() common.ResponseJson
		Logout() common.ResponseJson
}

type LoginRepositoryImpl struct {
		userService services.UserService
		smsService  services.SmsCodeService
		authService services.AuthService
		ctx         common.BaseRequestContext
}

func NewLoginRepository(ctx common.BaseRequestContext) LoginRepository {
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

// 登录逻辑总汇
func (this *LoginRepositoryImpl) Login() common.ResponseJson {
		var (
				typ     string
				request = new(transports.Login)
				input   = this.ctx.GetInput()
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

// 判断登录方式
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

// 手机号 + 短信验证码 登录
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
		if models.IsForbid(user) {
				return common.NewErrorResp(common.NewErrors(common.PermissionCode, "账号禁用状态"))
		}
		token := this.authService.Token(user)
		if token == "" {
				return common.NewErrorResp(common.NewErrors(common.ServiceFailed, "登陆服务异常"))
		}
		data := user.M(transforms.FilterUser)
		delete(data, "access_tokens")
		return common.NewSuccessResp(beego.M{"user": data, "token": token}, "登陆成功")
}

// 账号 + 密码 登录
func (this *LoginRepositoryImpl) loginByAccountPassword(account string, password string) common.ResponseJson {
		user := this.userService.GetByUserName(account)
		if user == nil {
				return common.NewErrorResp(common.NewErrors(1021, "用户不存在"))
		}
		if models.IsForbid(user) {
				return common.NewErrorResp(common.NewErrors(common.PermissionCode, "账号禁用状态"))
		}
		if !libs.PasswordVerify(user.PasswordHash, password) {
				return common.NewErrorResp(common.NewErrors(common.VerifyNotMatch, "密码不正确"))
		}
		token := this.authService.Token(user)
		data := user.M(getBaseUserInfoTransform())
		return common.NewSuccessResp(beego.M{"user": data, "token": token}, "登陆成功")
}

// 手机号+用户密码 登录
func (this *LoginRepositoryImpl) loginByMobilePassword(mobile string, password string) common.ResponseJson {
		user := this.userService.GetByMobile(mobile)
		if user == nil {
				return common.NewErrorResp(common.NewErrors(1021, "用户不存在"))
		}
		if models.IsForbid(user) {
				return common.NewErrorResp(common.NewErrors(common.PermissionCode, "账号禁用状态"))
		}
		if !libs.PasswordVerify(user.PasswordHash, password) {
				return common.NewErrorResp(common.NewErrors(common.VerifyNotMatch, "密码不正确"))
		}
		token := this.authService.Token(user)
		data := user.M(getBaseUserInfoTransform())
		return common.NewSuccessResp(beego.M{"user": data, "token": token}, "登陆成功")
}

// 邮箱 + 邮件验证码登录
func (this *LoginRepositoryImpl) loginByEmailCode(email string, code string) common.ResponseJson {
		user := this.userService.GetByEmail(email)
		if user == nil {
				return common.NewErrorResp(common.NewErrors(1021, "用户不存在"))
		}
		if models.IsForbid(user) {
				return common.NewErrorResp(common.NewErrors(common.PermissionCode, "账号禁用状态"))
		}
		// @todo 邮件验证码
		if code == "" {
				return common.NewErrorResp(common.NewErrors(common.VerifyNotMatch, "验证码不正确"))
		}
		token := this.authService.Token(user)
		data := user.M(getBaseUserInfoTransform())
		return common.NewSuccessResp(beego.M{"user": data, "token": token}, "登陆成功")
}

// 邮箱+登录密码 登录
func (this *LoginRepositoryImpl) loginByEmail(email string, password string) common.ResponseJson {
		user := this.userService.GetByEmail(email)
		if user == nil {
				return common.NewErrorResp(common.NewErrors(1021, "用户不存在"))
		}
		if models.IsForbid(user) {
				return common.NewErrorResp(common.NewErrors(common.PermissionCode, "账号禁用状态"))
		}
		if !libs.PasswordVerify(user.PasswordHash, password) {
				return common.NewErrorResp(common.NewErrors(common.VerifyNotMatch, "验证码不正确"))
		}
		token := this.authService.Token(user)
		data := user.M(getBaseUserInfoTransform())
		return common.NewSuccessResp(beego.M{"user": data, "token": token}, "登陆成功")
}

// 登出
func (this *LoginRepositoryImpl) Logout() common.ResponseJson {
		var id = getUserId(this.ctx)
		if id == "" {
				return common.NewUnLoginResp("please login!")
		}
		var err = this.getTokenService().Logout(id, getToken(this.ctx))
		if err == nil {
				this.ctx.Cookie(middlewares.AppAccessTokenHeader, nil)
				return common.NewSuccessResp(beego.M{"id": id, "timestamp": time.Now().Unix()}, "logout success!")
		}
		return common.NewFailedResp(common.ServiceFailed, common.NewErrors(err, "logout success!"))
}

func (this *LoginRepositoryImpl) getUserService() services.UserService {
		return services.UserServiceOf()
}

func (this *LoginRepositoryImpl) getTokenService() services.AuthService {
		return services.AuthServiceOf()
}
