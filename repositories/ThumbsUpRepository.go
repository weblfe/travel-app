package repositories

import (
		"github.com/astaxie/beego"
		"github.com/weblfe/travel-app/common"
		"github.com/weblfe/travel-app/services"
)

type ThumbsUpRepository interface {
		Up() common.ResponseJson
		Down() common.ResponseJson
		Count() common.ResponseJson
}

type thumbsUpRepositoryImpl struct {
		service services.ThumbsUpService
		ctx     *beego.Controller
}

func NewThumbsUpRepository(ctx *beego.Controller) ThumbsUpRepository {
		var repository = new(thumbsUpRepositoryImpl)
		repository.init()
		repository.ctx = ctx
		return repository
}

func (this *thumbsUpRepositoryImpl) init() {
		this.service = services.ThumbsUpServiceOf()
}

func (this *thumbsUpRepositoryImpl) Up() common.ResponseJson {
		var (
				ctx = this.ctx.Ctx
		)
		return common.NewInDevResp(ctx.Request.RequestURI)
}

func (this *thumbsUpRepositoryImpl) Down() common.ResponseJson {
		var (
				ctx = this.ctx.Ctx
		)
		return common.NewInDevResp(ctx.Request.RequestURI)
}

func (this *thumbsUpRepositoryImpl) Count() common.ResponseJson {
		var (
				ctx = this.ctx.Ctx
		)
		return common.NewInDevResp(ctx.Request.RequestURI)
}
