package repositories

import (
		"fmt"
		"github.com/astaxie/beego"
		"github.com/weblfe/travel-app/common"
		"github.com/weblfe/travel-app/middlewares"
		"github.com/weblfe/travel-app/services"
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
		ctx *beego.Controller
		userService services.UserService
}

func NewUserInfoRepository(ctx *beego.Controller) UserInfoRepository {
		var repository = new(UserInfoRepositoryImpl)
		repository.ctx = ctx
		repository.init()
		return repository
}

func (this *UserInfoRepositoryImpl)init()  {
		this.userService = services.UserServiceOf()
}

func (this *UserInfoRepositoryImpl)FocusOn()common.ResponseJson  {
		
		return nil
}

func (this *UserInfoRepositoryImpl) GetUserInfo() common.ResponseJson {
		var (
				id string
				v = this.ctx.GetSession(middlewares.AuthUserId)
		)
		fmt.Println("123231")
		if v != nil {
				id = v.(string)
		}
		if id == "" {
				return common.NewUnLoginResp(common.NewErrors(common.UnLoginCode,"请先登陆"))
		}
		user:=this.userService.GetById(id)
		if user == nil || isForbid(user) {
				return common.NewErrorResp(common.NewErrors(common.PermissionCode,"账号禁用状态"))
		}
		data:=user.M(filterUser)
		delete(data,"access_tokens")
		return common.NewSuccessResp(beego.M{"user":data},"获取成功")
}

func (this *UserInfoRepositoryImpl) ResetPassword() common.ResponseJson {
		panic("implement me")
}

func (this *UserInfoRepositoryImpl) GetUserFriends() common.ResponseJson {
		panic("implement me")
}

func (this *UserInfoRepositoryImpl) UpdateUserInfo() common.ResponseJson {
		panic("implement me")
}

func (this *UserInfoRepositoryImpl) FocusOff() common.ResponseJson {
		panic("implement me")
}


