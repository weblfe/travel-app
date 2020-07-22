package repositories

import "github.com/weblfe/travel-app/common"

type AppRepository interface {
		GetConfig(typ string) common.ResponseJson
}

type appRepository struct {
		ctx common.ResponseJson
}


func (this *appRepository)GetConfig(typ string) common.ResponseJson  {

		return common.NewErrorResp(common.NewErrors(common.NotFound,"config empty"),"配置为空")
}