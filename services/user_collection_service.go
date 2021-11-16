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
	Lists(userId string, limit models.ListsParams) ([]*models.TravelNotes, models.ListsParams)
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

// Add 添加收藏
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

// Remove 移除收藏
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

// Lists 罗列收藏
func (this *userCollectionServiceImpl) Lists(userId string, limit models.ListsParams) ([]*models.TravelNotes, models.ListsParams) {
	var (
		query = bson.M{
			"userId":     userId,
			"targetType": models.CollectTargetTypePost,
		}
		lists      = new([]*models.Collect)
		total, err = this.model.Lists(query, lists, limit)
	)
	if err != nil {
		return nil, nil
	}
	if lists == nil {
		return nil, models.NewMeta()
	}
	var (
		ids   []string
		meta  = models.NewMeta()
		notes []*models.TravelNotes
	)
	for _, v := range *lists {
		if v.TargetType != models.CollectTargetTypePost {
			continue
		}
		ids = append(ids, v.TargetId)
	}
	meta.Set("total", total)
	defer meta.Boot()
	notes = this.model.GetTravelNotesByIds(ids)
	if notes == nil {
		meta.Set("size", 0)
		return nil, meta
	}
	meta.Set("size", len(notes))
	return notes, meta
}
