package repositories

import (
		"github.com/astaxie/beego"
		"github.com/weblfe/travel-app/common"
		"github.com/weblfe/travel-app/services"
)

type TagsRepository interface {
		GetPostTags() common.ResponseJson
}

type tagsRepositoryImpl struct {
		ctx     common.BaseRequestContext
		service services.TagsService
}

func NewTagRepository(ctx common.BaseRequestContext) TagsRepository {
		var repository = new(tagsRepositoryImpl)
		repository.ctx = ctx
		repository.Init()
		return repository
}

func (this *tagsRepositoryImpl) Init() {
		this.service = services.TagsServiceOf()
}

// 获取作品标签
func (this *tagsRepositoryImpl) GetPostTags() common.ResponseJson {
		var tags, meta = this.service.GetTags(services.PostTagGroup)
		if tags == nil || meta == nil {
				return common.NewErrorResp(common.NewErrors(common.NotFound, "empty"), "标签不存在")
		}
		return common.NewSuccessResp(beego.M{"items": tags, "meta": meta}, "获取成功")
}
