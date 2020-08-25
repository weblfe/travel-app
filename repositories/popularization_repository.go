package repositories

import (
		"github.com/astaxie/beego"
		"github.com/weblfe/travel-app/common"
		"github.com/weblfe/travel-app/services"
)

type PopularizationRepository interface {
		GetChannel() common.ResponseJson
		GetChannelInfo() common.ResponseJson
		GetChannelQrcode() common.ResponseJson
}

// 推广渠道逻辑业务
type popularizationRepositoryImpl struct {
		ctx     common.BaseRequestContext
		service services.PopularizationService
}

func NewPopularizationRepository(ctx common.BaseRequestContext) PopularizationRepository {
		var repository = new(popularizationRepositoryImpl)
		repository.Init()
		repository.ctx = ctx
		return repository
}

func (this *popularizationRepositoryImpl) Init() {
		this.service = services.PopularizationServiceOf()
}

func (this *popularizationRepositoryImpl) GetChannelInfo() common.ResponseJson {
		var channel = this.ctx.GetString("ch")
		var info = this.service.GetChannelInfo(channel)
		if info == nil {
				return common.NewFailedResp(common.NotFound, "渠道不存在")
		}
		return common.NewSuccessResp(beego.M{"info": info}, "获取信息成功")
}

func (this *popularizationRepositoryImpl) GetChannelQrcode() common.ResponseJson {
		return common.NewErrorResp(nil, 1)
}

func (this *popularizationRepositoryImpl) GetChannel() common.ResponseJson {
		return common.NewErrorResp(nil, 1)
}
