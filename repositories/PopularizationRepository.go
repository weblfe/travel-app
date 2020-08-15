package repositories

import (
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
		service services.PostService
		ctx     common.BaseRequestContext
}

func NewPopularizationRepository(ctx common.BaseRequestContext) PopularizationRepository {
		var repository = new(popularizationRepositoryImpl)
		repository.Init()
		repository.ctx = ctx
		return repository
}

func (this *popularizationRepositoryImpl) Init() {

}

func (this *popularizationRepositoryImpl)GetChannelInfo() common.ResponseJson {
		return common.NewErrorResp(nil,1)
}

func (this *popularizationRepositoryImpl)GetChannelQrcode() common.ResponseJson {
		return common.NewErrorResp(nil,1)
}

func (this *popularizationRepositoryImpl)GetChannel() common.ResponseJson {
		return common.NewErrorResp(nil,1)
}
