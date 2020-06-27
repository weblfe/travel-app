package controllers

import (
		"github.com/weblfe/travel-app/common"
		"github.com/weblfe/travel-app/repositories"
)

type UserController struct {
		BaseController
}

// 用户模块 controller
func UserControllerOf() *UserController {
		return new(UserController)
}

// 用户登录接口
// @router /login [post]
func (this *UserController) Login() {
		var res = repositories.NewLoginRepository(&this.BaseController.Controller).Login()
		if res.IsSuccess() {
				token := res.GetDataByKey("token", "")
				if token != "" && token != nil {
						this.Ctx.SetCookie(common.AppTokenCookie, token.(string))
				}
		}
		this.Send(res)
}

// 用户注册接口
// @router /register [post]
func (this *UserController) Register() {
		var res = repositories.NewUserRegisterRepository(&this.BaseController.Controller).Register()
		// 注册成功后 记录相关
		if res.IsSuccess() {
				repositories.GetEventProvider().Dispatch("registerSuccess", res.GetData(), "UserRegister")
		}
		this.Send(res)
}

// 获取用户基本信息接口
// @router /user/info [get]
func (this *UserController) GetUserInfo() {
		this.Send(repositories.NewUserInfoRepository(&this.BaseController.Controller).GetUserInfo())
}

// 重置用户密码接口
// @router /reset/password [put]
func (this *UserController) ResetPassword() {
		this.Send(repositories.NewUserInfoRepository(&this.BaseController.Controller).ResetPassword())
}

// 获取用户好友列表接口
// @router /user/friends [get]
func (this *UserController) GetUserFriends() {
		this.Send(repositories.NewUserInfoRepository(&this.BaseController.Controller).GetUserFriends())
}

// 关注用户接口
// @router /focus/on/:userId [post]
func (this *UserController) FocusOn() {
		this.Send(repositories.NewUserInfoRepository(&this.BaseController.Controller).FocusOn())
}

// 取消关注接口
// @router /focus/off/:userId [delete]
func (this *UserController) FocusOff() {
		this.Send(repositories.NewUserInfoRepository(&this.BaseController.Controller).FocusOff())
}

// 更新用户信息接口
// @router /user/info [put]
func (this *UserController) UpdateUserInfo() {
		this.Send(repositories.NewUserInfoRepository(&this.BaseController.Controller).UpdateUserInfo())
}
