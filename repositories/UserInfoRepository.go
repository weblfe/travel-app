package repositories

import (
		"github.com/astaxie/beego"
		"github.com/weblfe/travel-app/common"
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
}

func NewUserInfoRepository(ctx *beego.Controller) UserInfoRepository {
		var repository = new(UserInfoRepositoryImpl)
		repository.ctx = ctx
		return repository
}

func (this *UserInfoRepositoryImpl)FocusOn()common.ResponseJson  {
		
		return nil
}

func (this *UserInfoRepositoryImpl) GetUserInfo() common.ResponseJson {
		return common.NewInvalidParametersResp()
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


