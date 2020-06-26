package repositories

import (
		"github.com/astaxie/beego"
		"github.com/weblfe/travel-app/common"
)

type UserRegisterRepository interface {
		Register() common.ResponseJson
}

type UserRegisterRepositoryImpl struct {
    ctx  *beego.Controller
}

func NewUserRegisterRepository(ctx *beego.Controller)UserRegisterRepository  {
		var repository = new(UserRegisterRepositoryImpl)
		repository.ctx = ctx
		return repository
}

func (this *UserRegisterRepositoryImpl)Register() common.ResponseJson  {
		

		return common.NewInvalidParametersResp("注册参数异常")
}
