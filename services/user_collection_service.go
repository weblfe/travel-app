package services

import (
		"github.com/astaxie/beego"
		"github.com/globalsign/mgo/bson"
		"github.com/weblfe/travel-app/common"
		"github.com/weblfe/travel-app/models"
		"time"
)

type UserCollectionService interface {
		Add(id, userId string) error
		Remove(id, userId string) error
		Lists(userId string, limit models.ListsParams) ([]*models.TravelNotes, *models.Meta)
}

type userCollectionServiceImpl struct {
		model *models.CollectModel
		BaseService
}

func UserCollectionServiceOf() UserCollectionService {
		var service = new(userCollectionServiceImpl)
		service.model = models.CollectModelOf()
		return service
}

// 添加收藏
func (this *userCollectionServiceImpl) Add(id, userId string) error {
		var collect = models.NewCollect()
		collect.TargetType = models.CollectTargetTypePost
		collect.UserId = userId
		collect.TargetId = id
		collect.Status = models.StatusOk
		collect.InitDefault()
		var query = bson.M{
				"targetType": collect.TargetType,
				"userId":     collect.UserId,
				"targetId":   collect.TargetId,
		}
		if this.model.Exists(query) {
				return nil
		}
		if !PostServiceOf().Exists(bson.M{"_id": bson.ObjectIdHex(id), "status": models.StatusOk}) {
				return common.NewErrors(common.RecordNotFound, "作品已下架")
		}
		if !UserServiceOf().Exists(beego.M{"_id": bson.ObjectIdHex(userId), "status": models.StatusOk}) {
				return common.NewErrors(common.RecordNotFound, "用户异常")
		}
		return this.model.Add(collect)
}

// 移除收藏
func (this *userCollectionServiceImpl) Remove(id, userId string) error {
		var collect = models.NewCollect()
		collect.TargetType = models.CollectTargetTypePost
		collect.UserId = userId
		collect.TargetId = id
		collect.Status = models.StatusActivity
		var query = bson.M{
				"targetType": collect.TargetType,
				"userId":     collect.UserId,
				"targetId":   collect.TargetId,
				"status":     collect.Status,
		}
		var err = this.model.FindOne(query, collect)
		if err != nil {
				return common.NewErrors(common.NotFound, "记录不存在")
		}
		collect.UpdatedAt = time.Now().Local()
		collect.Status = models.StatusCancel
		return this.model.UpdateById(collect.Id.Hex(), collect)
}

// 罗列收藏
func (this *userCollectionServiceImpl) Lists(userId string, limit models.ListsParams) ([]*models.TravelNotes, *models.Meta) {

		return nil, nil
}
