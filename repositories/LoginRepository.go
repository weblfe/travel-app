package repositories

import (
		"github.com/astaxie/beego"
		"github.com/weblfe/travel-app/common"
)

type LoginRepository interface {
		Login() common.ResponseJson
}

type LoginRepositoryImpl struct {
		ctx *beego.Controller
}

func NewLoginRepository(ctx *beego.Controller) LoginRepository {
		var repository = new(LoginRepositoryImpl)
		repository.ctx = ctx
		return repository
}

func (this *LoginRepositoryImpl)Login() common.ResponseJson  {
		
		return common.NewInvalidParametersResp()
}