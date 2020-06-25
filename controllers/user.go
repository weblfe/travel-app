package controllers

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

}

// 用户注册接口
// @router /register [post]
func (this *UserController) Register() {

}

// 获取用户基本信息接口
// @router /user/info [get]
func (this *UserController) GetUserInfo() {

}

// 重置用户密码接口
// @router /reset/password [put]
func (this *UserController) ResetPassword() {

}

// 获取用户好友列表接口
// @router /user/friends [get]
func (this *UserController) GetUserFriends() {

}

// 关注用户接口
// @router /focus/on/:userId [post]
func (this *UserController) FocusOn() {

}

// 取消关注接口
// @router /focus/off/:userId [delete]
func (this *UserController) FocusOff() {

}

// 更新用户信息接口
// @router /user/info [put]
func (this *UserController) UpdateUserInfo() {

}
