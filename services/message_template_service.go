package services

import (
		"github.com/globalsign/mgo/bson"
		"github.com/weblfe/travel-app/models"
)

type TemplateService interface {
		Add(template *models.MessageTemplate) error
		GetByNameType(name string, typ string) *models.MessageTemplate
		UpdateByNameType(name string, typ string, data map[string]interface{}) bool
		Lists(query interface{}, limit models.ListsParams, selects ...interface{}) ([]*TemplateServiceImpl, int, bool)
}

type TemplateServiceImpl struct {
		BaseService
		templateModel *models.MessageTemplateModel
}

func TemplateServiceOf() TemplateService {
		var service = new(TemplateServiceImpl)
		service.Init()
		return service
}

func (this *TemplateServiceImpl) Init() {
		this.init()
		this.Constructor = func(args ...interface{}) interface{} {
				return TemplateServiceOf()
		}
		this.templateModel = models.MessageTemplateModelOf()
}

// Add 添加
func (this *TemplateServiceImpl) Add(temp *models.MessageTemplate) error {
		return this.templateModel.Add(temp)
}

// GetByNameType 获取名
func (this *TemplateServiceImpl) GetByNameType(name string, typ string) *models.MessageTemplate {
		var template = models.NewMessageTemplate()
		if err := this.templateModel.FindOne(bson.M{"name": name, "type": typ}, template); err == nil {
				return template
		}
		return nil
}

// UpdateByNameType 更新
func (this *TemplateServiceImpl) UpdateByNameType(name string, typ string, data map[string]interface{}) bool {
		template := this.GetByNameType(name, typ)
		if template == nil {
				return false
		}
		if err := this.templateModel.Update(bson.M{"id": template.Id}, template.Load(data)); err == nil {
				return true
		}
		return false
}

// 列表
func (this *TemplateServiceImpl) Lists(query interface{}, limit models.ListsParams, selects ...interface{}) ([]*TemplateServiceImpl, int, bool) {
		var result []*TemplateServiceImpl
		total, err := this.templateModel.Lists(query, result, limit, selects...)
		if err != nil {
				return nil, 0, false
		}
		limit.SetTotal(total)
		return result, total, limit.More()
}
