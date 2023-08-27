package controllers

import (
		"github.com/weblfe/travel-app/common"
		"github.com/weblfe/travel-app/middlewares"
		"github.com/weblfe/travel-app/repositories"
		"log"
)

type UserController struct {
	BaseController
}

// UserControllerOf 用户模块 controller
func UserControllerOf() *UserController {
	return new(UserController)
}

// Login 用户登录接口
// @router /login [post]
func (this *UserController) Login() {
	var res = repositories.NewLoginRepository(this).Login()
	if res.IsSuccess() {
		token := res.GetDataByKey("token", "")
		if token != "" && token != nil {
			this.Cookie(common.AppTokenCookie, token.(string))
		}
	}
	this.Send(res)
}

// Register 用户注册接口
// @router /register [post]
func (this *UserController) Register() {
	var res = repositories.NewUserRegisterRepository(this).Register()
	// 注册成功后 记录相关
	if res.IsSuccess() {
		repositories.GetEventProvider().Dispatch("registerSuccess", res.GetData(), "UserRegister")
	}
	this.Send(res)
}

// Logout 用户注销登陆接口
// @router /logout  [delete]
func (this *UserController) Logout() {
	this.Send(repositories.NewLoginRepository(this).Logout())
}

// GetUserInfo 获取用户基本信息接口
// @router /user/info [get]
func (this *UserController) GetUserInfo() {
	this.Send(repositories.NewUserInfoRepository(this).GetUserInfo())
}

// GetUserInfoById 获取用户信息接口
// @router /user/info/public [get]
func (this *UserController) GetUserInfoById() {
	var id = this.GetString("id")
	this.Send(repositories.NewUserInfoRepository(this).GetUserPublicInfo(id))
}

// UpdateUserInfo 更新用户信息接口
// @router /user/info [put]
func (this *UserController) UpdateUserInfo() {
	this.Send(repositories.NewUserInfoRepository(this).UpdateUserInfo())
}

// ResetPassword 重置用户密码接口
// @router /reset/password [put]
func (this *UserController) ResetPassword() {
	this.Send(repositories.NewUserInfoRepository(this).ResetPassword())
}

// GetUserFriends 获取用户好友列表接口
// @router /user/friends [get]
func (this *UserController) GetUserFriends() {
	this.Send(repositories.NewUserInfoRepository(this).GetUserFriends())
}

// GetFriends 获取关注列表接口
// @router /friends/:userId [get]
func (this *UserController) GetFriends() {
	var userId = this.GetString(":userId", "0")
	this.Send(repositories.NewUserInfoRepository(this).GetUserFriends(userId))
}

// GetFriendsQuery 获取关注列表接口
// @router /friends [get]
func (this *UserController) GetFriendsQuery() {
	this.Send(repositories.NewUserInfoRepository(this).GetUserFriends(this.GetString("userId", "0")))
}

// GetUserFollowsQuery 获取其他用户关注接口
// @router /follows/public  [get]
func (this *UserController) GetUserFollowsQuery() {
	var userId = this.GetString("userId")
	this.Send(repositories.NewBehaviorRepository(this).GetUserFollows(userId))
}

// GetUserFollows 获取其他用户关注接口
// @router /follows/:userId [get]
func (this *UserController) GetUserFollows() {
	var userId, _ = this.GetParam(":userId", "0")
	this.Send(repositories.NewBehaviorRepository(this).GetUserFollows(userId.(string)))
}

// GetFollows 获取关注列表接口
// @router /follows [get]
func (this *UserController) GetFollows() {
	session := this.GetRequestContext().GetSession()
	userId := session.Get(middlewares.AuthUserId)
	if userId == nil || userId== "" {
		log.Printf("[Error] miss userId for follows")
	}
	this.Send(repositories.NewBehaviorRepository(this).GetUserFollows(userId.(string)))
}

// FocusOn 关注用户接口
// @router /follow/:userId [post]
func (this *UserController) FocusOn() {
	this.Send(repositories.NewBehaviorRepository(this).FocusOn())
}

// FocusOnQuery 关注用户接口
// @router /follow [post]
func (this *UserController) FocusOnQuery() {
	this.Send(repositories.NewBehaviorRepository(this).FocusOn())
}

// FocusOff 取消关注接口
// @router /follow/:userId [delete]
func (this *UserController) FocusOff() {
	this.Send(repositories.NewBehaviorRepository(this).FocusOff())
}

// FocusOffQuery 取消关注接口
// @router /follow [delete]
func (this *UserController) FocusOffQuery() {
	this.Send(repositories.NewBehaviorRepository(this).FocusOff())
}

// GetFans 获取粉丝接口
// @router /fans [get]
func (this *UserController) GetFans() {
	this.Send(repositories.NewBehaviorRepository(this).GetUserFans())
}

// GetUserFans 获取粉丝接口
// @router /fans/:userId [get]
func (this *UserController) GetUserFans() {
	var userId, _ = this.GetParam(":userId", "0")
	this.Send(repositories.NewBehaviorRepository(this).GetUserFans(userId.(string)))
}

// Search 用户搜索
// @router /user/search  [get]
func (this *UserController) Search() {
	this.Send(repositories.NewUserInfoRepository(this).Search(this.GetString("query")))
}

// GetProfile 用户个人页信息
//
//	@router /user/profile [get]
func (this *UserController) GetProfile() {
	this.Send(repositories.NewUserInfoRepository(this).GetProfile(this.GetString("userId")))
}

// AddCollect 收藏文章
// @router /user/collect/post  [post]
func (this *UserController) AddCollect() {
	var postId = this.GetString("postId")
	var userId = this.GetString("userId")
	if userId == "" {
		userId = this.GetUserId()
	}
	this.Send(repositories.NewUserCollectionRepository(this).Add(postId, userId))
}

// ListsCollect 罗列 收藏
// @router /user/collect/post [get]
func (this *UserController) ListsCollect() {
	var userId = this.GetString("userId")
	var page = this.GetInt("page")
	var limit = this.GetInt("count", 20)
	if userId == "" {
		userId = this.GetUserId()
	}
	this.Send(repositories.NewUserCollectionRepository(this).Lists(userId, page, limit))
}

// RemoveCollects 移除收藏
// @router /user/collect/post  [delete]
func (this *UserController) RemoveCollects() {
	var postId = this.GetString("postId")
	var userId = this.GetString("userId")
	if userId == "" {
		userId = this.GetUserId()
	}
	this.Send(repositories.NewUserCollectionRepository(this).Remove(postId, userId))
}
