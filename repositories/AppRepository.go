package repositories

import (
		"github.com/weblfe/travel-app/common"
		"github.com/weblfe/travel-app/services"
)

type AppRepository interface {
		GetConfig(typ string) common.ResponseJson
}

type appRepository struct {
		ctx     common.BaseRequestContext
		service services.AppService
}

func NewAppRepository(ctx common.BaseRequestContext) AppRepository {
		var repository = new(appRepository)
		repository.ctx = ctx
		repository.Init()
		return repository
}

func (this *appRepository) Init() {
		this.service = services.AppServiceOf()
}

func (this *appRepository) GetConfig(driver string) common.ResponseJson {
		var items = this.service.GetAppInfos(driver)
		if len(items) == 0 {
				return common.NewErrorResp(common.NewErrors(common.NotFound, "config empty"), "配置为空")
		}
		return common.NewSuccessResp(items, "获取配置成功")
}
