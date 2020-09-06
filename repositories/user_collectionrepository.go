package repositories

import (
		"github.com/astaxie/beego"
		"github.com/weblfe/travel-app/common"
		"github.com/weblfe/travel-app/services"
		"time"
)

type UserCollectionRepository interface {
		Add(id ,userId string) common.ResponseJson
		Remove(id ,userId string) common.ResponseJson
}

type userCollectionRepositoryImpl struct {
		ctx     common.BaseRequestContext
		service services.UserCollectionService
}

func NewUserCollectionRepository(ctx common.BaseRequestContext) UserCollectionRepository {
		var repository = new(userCollectionRepositoryImpl)
		repository.ctx = ctx
		repository.service = services.UserCollectionServiceOf()
		return repository
}

// 添加
func (this *userCollectionRepositoryImpl) Add(id string,userId string) common.ResponseJson {
		var err = this.service.Add(id,userId)
		if err == nil {
				return common.NewSuccessResp(beego.M{"timestamp": time.Now().Unix()}, "收藏成功")
		}
		return common.NewSuccessResp(beego.M{"timestamp": time.Now().Unix(), "id": id, "error": err.Error()}, "收藏失败")
}

// 移除
func (this *userCollectionRepositoryImpl) Remove(id string,userId string) common.ResponseJson {
		var err = this.service.Remove(id,userId)
		if err == nil {
				return common.NewSuccessResp(beego.M{"timestamp": time.Now().Unix()}, "移除成功")
		}
		return common.NewSuccessResp(beego.M{"id": id, "timestamp": time.Now().Unix()}, "移除失败")
}

// 移除
func (this *userCollectionRepositoryImpl) Lists(userId string) common.ResponseJson {
		var items, meta = this.service.Lists(userId,nil)
		if items != nil {
				return common.NewSuccessResp(beego.M{"items": items, "meta": meta}, "罗列成功")
		}
		return common.NewSuccessResp(nil, "空")
}
