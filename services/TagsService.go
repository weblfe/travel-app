package services

import "github.com/weblfe/travel-app/models"

type TagsService interface {
		GetTags(group string) ([]string, *models.Meta)
}

type tagsServiceImpl struct {
		model *models.TagModel
		BaseService
}

const (
		PostTagGroup = "post"
)

func TagsServiceOf() TagsService {
		var service = new(tagsServiceImpl)
		service.Init()
		return service
}

func (this *tagsServiceImpl) Init() {
		this.init()
		this.model = models.TagsModelOf()
		this.Constructor = func(args ...interface{}) interface{} {
				return TagsServiceOf()
		}
}

func (this *tagsServiceImpl) GetTags(group string) ([]string, *models.Meta) {
		var (
				items []string
				meta  = models.NewMeta()
		)
		items = this.model.GetTagsByGroup(group)
		meta.Page = 1
		meta.Count = len(items)
		meta.Total = meta.Count
		meta.Size = meta.Count
		return items, meta
}
