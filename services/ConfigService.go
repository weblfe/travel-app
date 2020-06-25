package services

import "github.com/weblfe/travel-app/models"

type ConfigService interface {

}

type ConfigServiceImpl struct {
		BaseService
		configModel *models.ConfigModel
}

func ConfigServiceOf() ConfigService  {
		var service = new(ConfigServiceImpl)
		service.Init()
		return service
}

func (this *ConfigServiceImpl)Init()  {
		this.configModel = models.ConfigModelOf()
		this.Constructor = func(args ...interface{}) interface{} {
				return ConfigServiceOf()
		}
}

