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
		"regexp"
		"time"
)

type UserInfoRepository interface {
		GetUserInfo() common.ResponseJson
		ResetPassword() common.ResponseJson
		GetUserFriends() common.ResponseJson
		UpdateUserInfo() common.ResponseJson
		FocusOff() common.ResponseJson
		FocusOn() common.ResponseJson
}

type UserInfoRepositoryImpl struct {
		ctx         common.BaseRequestContext
		userService services.UserService
}

func NewUserInfoRepository(ctx common.BaseRequestContext) UserInfoRepository {
		var repository = new(UserInfoRepositoryImpl)
		repository.ctx = ctx
		repository.init()
		return repository
}

func (this *UserInfoRepositoryImpl) init() {
		this.userService = services.UserServiceOf()
}

func (this *UserInfoRepositoryImpl) FocusOn() common.ResponseJson {

		return nil
}

func (this *UserInfoRepositoryImpl) GetUserInfo() common.ResponseJson {
		var (
				id string
				v  = this.ctx.Session(middlewares.AuthUserId)
		)
		if v != nil {
				id = v.(string)
		}
		if id == "" {
				return common.NewUnLoginResp(common.NewErrors(common.UnLoginCode, "请先登陆"))
		}
		user := this.userService.GetById(id)
		if user == nil || models.IsForbid(user) {
				return common.NewErrorResp(common.NewErrors(common.PermissionCode, "账号禁用状态"))
		}
		data := user.M(transforms.FilterUserBase)
		return common.NewSuccessResp(beego.M{"user": data}, "获取成功")
}

func (this *UserInfoRepositoryImpl) ResetPassword() common.ResponseJson {
		var (
				ctx     = this.ctx.GetParent().Ctx
				request = new(transports.ResetPassword)
				userId  = this.ctx.GetParent().GetSession(middlewares.AuthUserId)
		)
		request.Load(ctx.Input)
		if request.Password == "" {
				return common.NewInvalidParametersResp("password参数缺失")
		}
		if err := this.PasswordCheck(request.Password); err != nil {
				return common.NewInvalidParametersResp(err)
		}
		// 登陆用户更新自己密码
		if request.Mobile == "" && request.Code == "" {
				if userId == "" || userId == nil {
						return common.NewInvalidParametersResp("异常请求")
				}
				request.UserId = userId.(string)
				if request.CurrentPassword == "" {
						return common.NewInvalidParametersResp("请输入当前登陆密码")
				}
				err := this.resetPasswordByCurrent(request.UserId, request.CurrentPassword, request.Password)
				if err == nil {
						return common.NewSuccessResp(beego.M{"timestamp": time.Now().Unix(), "tip": "请求重新登陆"}, "更新密码成功")
				}
				return common.NewFailedResp(common.ServiceFailed, "更新密码失败", err)
		}
		// 手机+验证码重置密码
		err := this.resetPasswordByMobileSms(request.Mobile, request.Code, request.Password)
		if err == nil {
				return common.NewSuccessResp(beego.M{"timestamp": time.Now().Unix(), "tip": "请求重新登陆"}, "更新密码成功")
		}
		return common.NewFailedResp(common.ServiceFailed, "更新密码失败", err)
}

// 密码规则检查
func (this *UserInfoRepositoryImpl) PasswordCheck(pass string) common.Errors {
		if len(pass) < 6 {
				return common.NewErrors(common.InvalidParametersCode, "密码长度必须不少于6位")
		}
		var reg = regexp.MustCompile(`[\w@+.$&]{6,20}`)
		if !reg.MatchString(pass) {
				return common.NewErrors(common.InvalidParametersCode, "特殊越界,密码特殊字符是能包含(@,_,+,-,.,$,&)")
		}
		return nil
}

func (this *UserInfoRepositoryImpl) resetPasswordByCurrent(userId string, current string, password string) error {
		var (
				user = this.userService.GetById(userId)
		)
		if user == nil {
				return common.NewErrors("用户不存在")
		}
		if !libs.PasswordVerify(user.PasswordHash, current) {
				return common.NewErrors("密码不正确")
		}
		user.PasswordHash = ""
		// 更新密码
		user.Load(beego.M{
				"passwordHash": password,
		})
		// 仅更新密码
		data := beego.M{
				"passwordHash":       user.PasswordHash,
				"updatedAt":          time.Now(),
				"resetPasswordTimes": user.ResetPasswordTimes + 1,
				"modifies":           []string{"passwordHash", "updatedAt", "resetPasswordTimes"},
		}
		err := this.userService.UpdateByUid(userId, data)
		if err != nil {
				return common.NewErrors(err)
		}
		// 释放所有登陆token
		defer services.AuthServiceOf().ReleaseByUserId(userId)
		return nil
}

func (this *UserInfoRepositoryImpl) resetPasswordByMobileSms(mobile string, code string, password string) error {
		var (
				user = this.userService.GetByMobile(mobile)
		)
		if user == nil {
				return common.NewErrors("用户不存在")
		}
		// 校验
		if !services.SmsCodeServiceOf().Verify(mobile, code, "reset_password") {
				return common.NewErrors("验证码不正确")
		}
		user.PasswordHash = ""
		// 更新密码
		user.Load(beego.M{
				"passwordHash": password,
		})
		// 仅更新密码
		data := beego.M{
				"passwordHash":       user.PasswordHash,
				"updatedAt":          time.Now(),
				"resetPasswordTimes": user.ResetPasswordTimes + 1,
				"modifies":           []string{"passwordHash", "updatedAt", "resetPasswordTimes"},
		}
		err := this.userService.UpdateByUid(user.Id.Hex(), data)
		if err != nil {
				return common.NewErrors(err)
		}
		// 释放所有登陆token
		defer services.AuthServiceOf().ReleaseByUserId(user.Id.Hex())
		return nil
}

func (this *UserInfoRepositoryImpl) GetUserFriends() common.ResponseJson {
		return common.NewInDevResp(this.ctx.GetActionId())
}

func (this *UserInfoRepositoryImpl) UpdateUserInfo() common.ResponseJson {
		var (
				updateInfo    beego.M
				request       = this.ctx.GetParent()
				updateRequest = new(transports.UpdateUserRequest)
				userId        = request.GetSession(middlewares.AuthUserId)
				err           = updateRequest.Load(request.Ctx.Input.RequestBody)
		)
		if userId == nil || userId == "" {
				return common.NewUnLoginResp("登陆失效,请重新登陆！")
		}
		if err != nil {
				return common.NewErrorResp(common.NewErrors(err.Error(), common.ServiceFailed), "参数解析异常！")
		}
		if updateRequest.Empty() {
				return common.NewErrorResp(common.NewErrors(common.InvalidParametersError, common.InvalidParametersCode), "缺失请求参数")
		}
		updateInfo = updateRequest.M()
		if len(updateInfo) == 0 {
				return common.NewInvalidParametersResp("参数缺失!")
		}
		// 指定更新字段
		err = this.userService.UpdateByUid(userId.(string), updateInfo)
		if err == nil {
				return common.NewSuccessResp(beego.M{"timestamp": time.Now().Unix()}, "更新成功")
		}
		return common.NewFailedResp(common.ServiceFailed, err, "更新失败")
}

func (this *UserInfoRepositoryImpl) FocusOff() common.ResponseJson {
		return common.NewInDevResp(this.ctx.GetActionId())
}
