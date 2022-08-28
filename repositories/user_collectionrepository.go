package repositories

import (
	"github.com/astaxie/beego"
	"github.com/weblfe/travel-app/common"
	"github.com/weblfe/travel-app/models"
	"github.com/weblfe/travel-app/services"
	"time"
)

type UserCollectionRepository interface {
	Add(id, userId string) common.ResponseJson
	Remove(id, userId string) common.ResponseJson
	Lists(id string, page, count int) common.ResponseJson
}

type userCollectionRepositoryImpl struct {
	ctx      common.BaseRequestContext
	service  services.UserCollectionService
	postRepo PostsRepository
}

func NewUserCollectionRepository(ctx common.BaseRequestContext) UserCollectionRepository {
	var repository = new(userCollectionRepositoryImpl)
	repository.ctx = ctx
	repository.service = services.UserCollectionServiceOf()
	repository.postRepo = NewPostsRepository(ctx)
	return repository
}

// Add 添加
func (this *userCollectionRepositoryImpl) Add(id string, userId string) common.ResponseJson {
	var err = this.service.Add(id, userId,this.ctx.GetString("type",models.CollectTargetTypePost))
	if err == nil {
		return common.NewSuccessResp(beego.M{"timestamp": time.Now().Unix()}, "收藏成功")
	}
	if e, ok := err.(common.Errors); ok {
		return common.NewErrorResp(e, "收藏失败")
	}
	return common.NewSuccessResp(beego.M{"timestamp": time.Now().Unix(), "id": id, "error": err.Error()}, "收藏失败")
}

// Remove 移除
func (this *userCollectionRepositoryImpl) Remove(id string, userId string) common.ResponseJson {
	var err = this.service.Remove(id, userId,this.ctx.GetString("type",models.CollectTargetTypePost))
	if err == nil {
		return common.NewSuccessResp(beego.M{"timestamp": time.Now().Unix()}, "移除成功")
	}
	return common.NewSuccessResp(beego.M{"id": id, "timestamp": time.Now().Unix()}, "移除失败")
}

// Lists 列表
func (this *userCollectionRepositoryImpl) Lists(userId string, page, limit int) common.ResponseJson {
	var	param = models.NewListParam(page, limit)
	param.SetArg("types",this.ctx.GetStrings("types"))
	param.Order(`updatedAt`,`desc`)
	var items, meta = this.service.Lists(userId,param)
	if items != nil {
		var (
			lists     []interface{}
			transform = this.postRepo.GetPostTransform()
		)
		for _, v := range items {
			lists = append(lists, v.M(transform))
		}
		return common.NewSuccessResp(beego.M{"items": lists, "meta": meta}, "罗列成功")
	}
	return common.NewSuccessResp(nil, "空")
}
