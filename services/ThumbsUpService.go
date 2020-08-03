package services

import (
		"github.com/globalsign/mgo/bson"
		"github.com/weblfe/travel-app/models"
)

type ThumbsUpService interface {
		Exists(query bson.M) bool
		Count(string, string, ...string) int
		Up(typ string, typeId string, userId string) int
		Down(typ string, typeId string, userId string) int
		GetUserLikeLists(query bson.M, limit models.ListsParams) ([]*models.TravelNotes, *models.Meta)
}

type thumbsUpServiceImpl struct {
		BaseService
		model *models.ThumbsUpModel
}

const (
		ThumbsUpActUp   = 1
		ThumbsUpActDown = 0
)

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

// 点赞
func (this *thumbsUpServiceImpl) Up(typ string, typeId string, userId string) int {
		var (
				data = bson.M{
						"type": typ, "typeId": typeId, "userId": userId,
				}
				up = new(models.ThumbsUp)
		)
		err := this.model.FindOne(data, up)
		// 新的点赞
		if err != nil {
				up.Type = typ
				up.TypeId = typeId
				up.UserId = userId
				err = up.Defaults().Save()
				if err == nil {
						defer this.After(typ, typeId, userId, ThumbsUpActUp)
				}
				return this.Count(typ, typeId)
		}
		if up.Status == 1 && up.Count == 1 {
				return this.Count(typ, typeId)
		}
		// 更新取消的点赞
		up.Count = 1
		up.Status = 1
		up.UserId = userId
		up.Type = typ
		up.TypeId = typeId
		err = up.Save()
		if err == nil {
				defer this.After(typ, typeId, userId, ThumbsUpActUp)
		}
		return this.Count(typ, typeId)
}

// 取消点赞
func (this *thumbsUpServiceImpl) Down(typ string, typeId string, userId string) int {
		var (
				data = bson.M{
						"type": typ, "typeId": typeId, "userId": userId,
				}
				up = new(models.ThumbsUp)
		)
		err := this.model.FindOne(data, up)
		// 点赞记录不存在
		if err != nil {
				return this.Count(typ, typeId)
		}
		if up.Status == 0 {
				return this.Count(typ, typeId)
		}
		up.Status = 0
		up.Count = 0
		up.TypeId = typeId
		up.Type = typ
		up.UserId = userId
		err = up.Save()
		// 取消点赞成功
		if err == nil {
				defer this.After(typ, typeId, userId, ThumbsUpActDown)
		}
		return this.Count(typ, typeId)
}

// 获取对于数据点赞数量
func (this *thumbsUpServiceImpl) Count(typ string, typId string, userId ...string) int {
		var (
				data = bson.M{
						"type": typ, "typeId": typId,
				}
		)
		// 对用用户点赞对用数据点赞次数
		if len(userId) != 0 {
				data["userId"] = userId[0]
		}
		// 未指定ID 仅限定类型
		if typId == "" {
				delete(data, "typeId")
		}
		return this.model.Sum(data, "count")
}

// 获取点赞列表
func (this *thumbsUpServiceImpl) Lists(query bson.M, limit models.ListsParams) ([]*models.ThumbsUp, *models.Meta) {
		var (
				err   error
				meta  = models.NewMeta()
				items = make([]*models.ThumbsUp, 2)
		)
		items = items[:0]
		Query := this.model.NewQuery(query)
		err = Query.Limit(limit.Count()).Skip(limit.Skip()).All(&items)
		if err == nil {
				meta.Count = limit.Count()
				meta.Size = len(items)
				meta.Page = limit.Page()
				meta.Total, _ = this.model.NewQuery(query).Count()
				meta.Boot()
				return items, meta
		}
		return nil, meta
}

// 获取点赞列表
func (this *thumbsUpServiceImpl) GetUserLikeLists(query bson.M, limit models.ListsParams) ([]*models.TravelNotes, *models.Meta) {
		var (
				err     error
				meta    = models.NewMeta()
				ids     []bson.ObjectId
				results []*models.TravelNotes
				items   = make([]*models.ThumbsUp, 2)
		)
		items = items[:0]
		query["type"] = models.ThumbsTypePost
		query["status"] = models.StatusOk
		Query := this.model.NewQuery(query)
		err = Query.Limit(limit.Count()).Skip(limit.Skip()).Sort("-createdAt").All(&items)
		// 获取用户喜欢列表
		if err == nil {
				for _, it := range items {
						ids = append(ids, bson.ObjectIdHex(it.TypeId))
				}
				results, meta = PostServiceOf().ListsQuery(bson.M{"_id": bson.M{"$in": ids}}, limit)
				// 排序
				if results != nil {
						for i, id := range ids {
								for j, it := range results {
										if id == it.Id && i != j {
												results[i], results[j] = results[j], results[i]
										}
								}
						}
				}
				return results, meta
		}
		return nil, meta
}

// 点赞之后 [log,updateNum]
func (this *thumbsUpServiceImpl) After(ty, id, userId string, act int) bool {
		if ty == "" || id == "" || userId == "" {
				return false
		}
		var (
				err         error
				isPostUp    = false
				isCommentUp = false
		)
		// 作品点赞
		if ty == models.ThumbsTypePost {
				isPostUp = true
		}
		// 评论点赞
		if ty == models.ThumbsTypeComment {
				isCommentUp = false
		}
		// 是否 游记相关点赞
		if isPostUp {
				switch act {
				case ThumbsUpActUp:
						err = PostServiceOf().IncrThumbsUp(id, 1)
				case ThumbsUpActDown:
						err = PostServiceOf().IncrThumbsUp(id, -1)
				}
		}
		// 点赞评论
		if isCommentUp {
				switch act {
				case ThumbsUpActUp:
						err = CommentServiceOf().IncrThumbsUp(id, 1)
				case ThumbsUpActDown:
						err = CommentServiceOf().IncrThumbsUp(id, -1)
				}
		}

		return err == nil
}

// 是否存在
func (this *thumbsUpServiceImpl) Exists(query bson.M) bool {
		var n, err = this.model.NewQuery(query).Count()
		if err != nil {
				return false
		}
		return n >= 1
}
