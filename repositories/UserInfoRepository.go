package repositories

import (
		"encoding/json"
		"github.com/astaxie/beego"
		"github.com/astaxie/beego/context"
		"github.com/weblfe/travel-app/common"
		"github.com/weblfe/travel-app/libs"
		"github.com/weblfe/travel-app/middlewares"
		"github.com/weblfe/travel-app/services"
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
		ctx         *beego.Controller
		userService services.UserService
}

func NewUserInfoRepository(ctx *beego.Controller) UserInfoRepository {
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
				v  = this.ctx.GetSession(middlewares.AuthUserId)
		)
		if v != nil {
				id = v.(string)
		}
		if id == "" {
				return common.NewUnLoginResp(common.NewErrors(common.UnLoginCode, "请先登陆"))
		}
		user := this.userService.GetById(id)
		if user == nil || isForbid(user) {
				return common.NewErrorResp(common.NewErrors(common.PermissionCode, "账号禁用状态"))
		}
		data := user.M(filterUserBase)
		return common.NewSuccessResp(beego.M{"user": data}, "获取成功")
}

func (this *UserInfoRepositoryImpl) ResetPassword() common.ResponseJson {
		var (
				ctx     = this.ctx.Ctx
				request = new(ResetPassword)
				userId  =  this.ctx.GetSession(middlewares.AuthUserId)
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
		if services.SmsCodeServiceOf().Verify(mobile, code, "reset_password") {
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
		return common.NewInDevResp(this.ctx.Ctx.Request.URL.String())
}

func (this *UserInfoRepositoryImpl) UpdateUserInfo() common.ResponseJson {
		var (
				request       = this.ctx
				updateInfo    beego.M
				updateRequest = new(UpdateUserRequest)
				userId        = request.GetSession(middlewares.AuthUserId)
				err           = json.Unmarshal(request.Ctx.Input.RequestBody, updateRequest)
		)
		if userId == nil || userId == "" {
				return common.NewUnLoginResp("登陆失效,请重新登陆！")
		}
		if err != nil {
				return common.NewErrorResp(common.NewErrors(err.Error(), common.ServiceFailed), "参数解析异常！")
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
		return common.NewInDevResp(this.ctx.Ctx.Request.URL.String())
}

// 请求参数
type UpdateUserRequest struct {
		AvatarId string   `json:"avatarId,omitempty"`
		NickName string   `json:"nickname,omitempty"`
		Email    string   `json:"email,omitempty"`
		Gender   int      `json:"gender,omitempty"`
		Intro string `json:"intro,omitempty"`
		BackgroundCoverId string `json:"backgroundCoverId,omitempty"`
		Modifies []string `json:"modifies,omitempty"`
}

func (this *UpdateUserRequest) M(filters ...func(m beego.M) beego.M) beego.M {
		var data = beego.M{
				"avatarId": this.AvatarId,
				"nickname": this.NickName,
				"email":    this.Email,
				"gender":   this.Gender,
				"modifies": this.Modifies,
		}
		if this.Gender == 0 {
				delete(data, "gender")
		}
		if len(this.Modifies) == 0 {
				delete(data, "modifies")
		}
		filters = append(filters, filterEmpty)
		for _, filter := range filters {
				data = filter(data)
		}
		return data
}

// 重置密码请求体
type ResetPassword struct {
		Password        string `json:"password"`                  // 新密码
		CurrentPassword string `json:"currentPassword,omitempty"` // 当前登陆使用的密码
		UserId          string `json:"userId,omitempty"`          // 当前用户ID
		Code            string `json:"code,omitempty"`            // 手机重置密码使用的验证码
		Mobile          string `json:"mobile,omitempty"`          // 手机号
}

func (this *ResetPassword) Load(ctx *context.BeegoInput) *ResetPassword {
		var (
			_ = json.Unmarshal(ctx.RequestBody, this)
			mapper =map[string]interface{}{
					"code":&this.Code,
					"mobile":&this.Mobile,
					"password":&this.Password,
					"currentPassword":&this.CurrentPassword,
			}
		)
		if this.Password == "" {
				for key,addr:=range mapper{
						_ = ctx.Bind(addr,key)
				}
		}
		return this
}

func (this *ResetPassword) M(filters ...func(m beego.M) beego.M) beego.M {
		var data = beego.M{
				"userId":          this.UserId,
				"password":        this.Password,
				"currentPassword": this.CurrentPassword,
				"code":            this.Code,
				"mobile":          this.Mobile,
		}
		filters = append(filters, filterEmpty)
		for _, filter := range filters {
				data = filter(data)
		}
		return data
}
