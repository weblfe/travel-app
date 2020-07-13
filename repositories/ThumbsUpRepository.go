package repositories

import (
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
		ctx     common.BaseRequestContext
}

func NewThumbsUpRepository(ctx common.BaseRequestContext) ThumbsUpRepository {
		var repository = new(thumbsUpRepositoryImpl)
		repository.init()
		repository.ctx = ctx
		return repository
}

func (this *thumbsUpRepositoryImpl) init() {
		this.service = services.ThumbsUpServiceOf()
}

func (this *thumbsUpRepositoryImpl) Up() common.ResponseJson {

		return common.NewInDevResp(this.ctx.GetActionId())
}

func (this *thumbsUpRepositoryImpl) Down() common.ResponseJson {
		return common.NewInDevResp(this.ctx.GetActionId())
}

func (this *thumbsUpRepositoryImpl) Count() common.ResponseJson {
		return common.NewInDevResp(this.ctx.GetActionId())
}
