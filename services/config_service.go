package services

import (
		"github.com/globalsign/mgo/bson"
		"github.com/weblfe/travel-app/models"
)

type ConfigService interface {
		Adds(items []map[string]interface{}) error
		Update(data map[string]string) error
}

type ConfigServiceImpl struct {
		BaseService
		configModel *models.ConfigModel
}

func ConfigServiceOf() ConfigService {
		return newConfigService()
}

func newConfigService() *ConfigServiceImpl {
		var service = new(ConfigServiceImpl)
		service.Init()
		return service
}

func (this *ConfigServiceImpl) Init() {
		this.configModel = models.ConfigModelOf()
		this.Constructor = func(args ...interface{}) interface{} {
				return ConfigServiceOf()
		}
}

func (this *ConfigServiceImpl) Adds(items []map[string]interface{}) error {
		if len(items) == 0 {
				return models.ErrEmptyData
		}
		var result []interface{}
		for _, it := range items {
				config := this.configModel.GetByUnique(it)
				if config != nil {
						_ = this.configModel.Update(bson.M{"_id": config.Id}, it)
				} else {
						var config = models.NewConfig()
						config.SetAttributes(it, false)
						config.InitDefault()
						result = append(result, config)
				}
		}
		if len(result) == 0 {
				return nil
		}
		if err := this.configModel.Inserts(result); err != nil {
				return err
		}
		return nil
}

// 更新配置
func (this *ConfigServiceImpl) Update(data map[string]string) error {
		return nil
}
