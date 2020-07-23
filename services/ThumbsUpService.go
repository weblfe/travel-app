package services

import (
		"github.com/globalsign/mgo/bson"
		"github.com/weblfe/travel-app/models"
)

type ThumbsUpService interface {
		Count(string, string, ...string) int
		Up(typ string, typeId string, userId string) int
		Down(typ string, typeId string, userId string) int
}

type thumbsUpServiceImpl struct {
		BaseService
		model *models.ThumbsUpModel
}

func ThumbsUpServiceOf() ThumbsUpService {
		var service = new(thumbsUpServiceImpl)
		service.Init()
		return service
}

func (this *thumbsUpServiceImpl) Init() {
		this.init()
		this.model = models.ThumbsUpModelOf()
		this.Constructor = func(args ...interface{}) interface{} {
				return ThumbsUpServiceOf()
		}
}

func (this *thumbsUpServiceImpl) Up(typ string, typeId string, userId string) int {
		var (
				data = bson.M{
						"type": typ, "typeId": typeId, "userId": userId,
				}
				up = new(models.ThumbsUp)
		)
		err := this.model.FindOne(data, up)
		if err != nil {
				err = up.Defaults().Save()
				if err == nil {
						return this.Count(typ, typeId)
				}
		}
		up.Count = 1
		up.Status = 1
		up.UserId = userId
		up.Type =typ
		up.TypeId = typeId
		_ = up.Save()
		return this.Count(typ, typeId)
}

func (this *thumbsUpServiceImpl) Down(typ string, typeId string, userId string) int {
		var (
				data = bson.M{
						"type": typ, "typeId": typeId, "userId": userId,
				}
				up = new(models.ThumbsUp)
		)
		err := this.model.FindOne(data, up)
		if err != nil {
				return this.Count(typ, typeId)
		}
		up.Count = 0
		up.Status = 0
		up.TypeId= typeId
		up.Type = typ
		up.UserId = userId
		_ = up.Save()
		return this.Count(typ, typeId)
}

func (this *thumbsUpServiceImpl) Count(typ string, typId string, userId ...string) int {
		var (
				data = bson.M{
						"type": typ, "typeId": typId,
				}
		)
		if len(userId) != 0 {
				data["userId"] = userId[0]
		}
		return this.model.Sum(data, "count")
}
