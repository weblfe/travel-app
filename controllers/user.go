package controllers

import "github.com/astaxie/beego"

type UserController struct {
		beego.Controller
}

// 用户模块 controller
func UserControllerOf() *UserController {
		return new(UserController)
}

// 用户登录接口
// @route /login [post]
func (this *UserController) Login() {

}

// 用户注册接口
// @route /register [post]
func (this *UserController) Register() {

}

// 获取用户基本信息接口
// @route /user/info [get]
func (this *UserController) GetUserInfo() {

}

// 重置用户密码接口
// @route /reset/password [put]
func (this *UserController) ResetPassword() {

}

// 获取用户好友列表接口
// @route /user/friends [get]
func (this *UserController) GetUserFriends() {

}

// 关注用户接口
// @route /focus/on/:userId [post]
func (this *UserController) FocusOn() {

}

// 取消关注接口
// @route /focus/off/:userId [delete]
func (this *UserController) FocusOff() {

}

// 更新用户信息接口
// @route /user/info [put]
func (this *UserController) UpdateUserInfo() {

}
