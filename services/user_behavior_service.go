package services

import (
		"fmt"
		"github.com/astaxie/beego"
		"github.com/globalsign/mgo/bson"
		"github.com/weblfe/travel-app/common"
		"github.com/weblfe/travel-app/libs"
		"github.com/weblfe/travel-app/models"
		"time"
)

type UserBehaviorService interface {
		GetUserFansNumber(userId string) int64
		IsFriend(userId, follower string) bool
		GetUserFollowNumber(userId string) int64
		IsFollowed(userId, follower string) bool
		UnFollow(userId string, targetUserId string, extras ...beego.M) error
		Follow(userId string, targetUserId string, extras ...beego.M) error
		GetFans(userId string, limit models.ListsParams) ([]bson.ObjectId, *models.Meta)
		GetFollows(userId string, limit models.ListsParams) ([]bson.ObjectId, *models.Meta)
		GetFollowPostsLists(query bson.M, limit models.ListsParams) ([]*models.TravelNotes, *models.Meta)
		ListsByUserId(userId string, limit models.ListsParams, extras ...beego.M) ([]*models.UserFollow, *models.Meta)
}

type userBehaviorServiceImpl struct {
		BaseService
}

func UserBehaviorServiceOf() UserBehaviorService {
		return newUserBehaviorServiceInstance()
}

func newUserBehaviorServiceInstance() *userBehaviorServiceImpl {
		var service = new(userBehaviorServiceImpl)
		service.init()
		return service
}

// 罗列用户关注的用户列表
func (this *userBehaviorServiceImpl) ListsByUserId(userId string, limit models.ListsParams, extras ...beego.M) ([]*models.UserFollow, *models.Meta) {
		if len(extras) == 0 {
				extras = append(extras, beego.M{"status": 1})
		}
		var (
				meta  = models.NewMeta()
				items = make([]*models.UserFollow, 2)
				query = beego.M{"userId": bson.ObjectIdHex(userId)}
		)
		items = items[:0]
		Query := this.getUserFocusModel().NewQuery(bson.M(models.Merger(query, extras[0])))
		err := Query.Limit(limit.Count()).Skip(limit.Skip()).All(&items)
		if err == nil {
				meta.Size = len(items)
				meta.Page = limit.Page()
				meta.Total, _ = Query.Count()
				meta.Count = limit.Count()
				meta.Boot()
				return items, meta
		}
		return nil, meta
}

// 用户关注
func (this *userBehaviorServiceImpl) Follow(userId string, targetUserId string, extras ...beego.M) error {
		if len(extras) == 0 {
				extras = append(extras, beego.M{"status": 1})
		}
		if userId == targetUserId {
				return common.NewErrors(common.InvalidParametersCode, "用户不能自己关注自己")
		}
		var (
				err   error
				data  = models.NewUserFocus()
				model = this.getUserFocusModel()
				query = beego.M{"userId": bson.ObjectIdHex(userId), "focusUserId": bson.ObjectIdHex(targetUserId)}
		)
		it := model.GetByUnique(query)
		if it != nil {
				it.Status = models.StatusOk
				it.UpdatedAt = time.Now().Local()
				err = model.Update(bson.M{"_id": it.Id}, it)
				if err == nil {
						go this.followAfter(data)
				}
				return err
		}
		var extrasInfo = extras[0]
		data.Status = models.StatusOk
		if id, ok := extrasInfo["targetId"]; ok {
				data.TargetId = this.id(id)
		}
		data.UserId = bson.ObjectIdHex(userId)
		data.FocusUserId = bson.ObjectIdHex(targetUserId)
		data.InitDefault()
		err = model.Add(data)
		if err == nil {
				go this.followAfter(data)
		}
		return err
}

// 用户关注
func (this *userBehaviorServiceImpl) UnFollow(userId string, targetUserId string, extras ...beego.M) error {
		if len(extras) == 0 {
				extras = append(extras, beego.M{"status": 1})
		}
		var (
				err   error
				model = this.getUserFocusModel()
				query = beego.M{"userId": bson.ObjectIdHex(userId), "focusUserId": bson.ObjectIdHex(targetUserId)}
		)
		it := model.GetByUnique(query)
		if it != nil {
				it.UpdatedAt = time.Now().Local()
				it.Status = models.StatusCancel
				err = model.Update(bson.M{"_id": it.Id}, it)
				if err == nil {
						go this.followAfter(it)
				}
				return err
		}
		return common.NewErrors(common.RecordNotFound, common.RecordNotFoundError)
}

// 获取用户关注用户
func (this *userBehaviorServiceImpl) GetUserFollowNumber(userId string) int64 {
		var (
				model = this.getUserFocusModel()
				query = bson.M{"userId": bson.ObjectIdHex(userId), "status": 1}
		)
		n, err := model.NewQuery(query).Count()
		if err == nil {
				return int64(n)
		}
		return 0
}

// 用户粉丝数量
func (this *userBehaviorServiceImpl) GetUserFansNumber(userId string) int64 {
		var (
				model = this.getUserFocusModel()
				query = bson.M{"focusUserId": bson.ObjectIdHex(userId), "status": 1}
		)
		n, err := model.NewQuery(query).Count()
		if err == nil {
				return int64(n)
		}
		return 0
}

// 关注之后
func (this *userBehaviorServiceImpl) followAfter(focus *models.UserFollow) {
		if focus == nil {
				return
		}
		switch focus.Status {
		case models.StatusOk:
				this.success(focus)
		case models.StatusCancel:
				this.cancel(focus)
		}
}

// 取消好友关系
func (this *userBehaviorServiceImpl) cancel(focus *models.UserFollow) bool {
		var (
				extras = beego.M{
						"unFollow":   time.Now().Unix(),
						"status":     models.StatusCancel,
						"targetType": models.TargetTypeFriend,
				}
				err = this.getUserRelation().SaveInfo(focus.UserId.Hex(), focus.TargetId.Hex(), extras)
		)
		if err == nil {
				return true
		}
		return false
}

// 建立好友关系
func (this *userBehaviorServiceImpl) success(focus *models.UserFollow) bool {
		var (
				err    error
				model  = this.getUserRelation()
				extras = beego.M{
						"status":     models.StatusOk,
						"follow":     time.Now().Unix(),
						"targetType": models.TargetTypeFriend,
				}
				query = beego.M{"userId": focus.FocusUserId.Hex(), "targetUserId": focus.UserId.Hex()}
		)
		var data = this.getUserFocusModel().GetByUnique(query)
		// 对方未关注
		if data == nil {
				return false
		}
		// 用户未互相关注
		if data.Status == models.StatusCancel || data.FocusUserId != focus.UserId {
				return false
		}
		// 保存好友关系
		err = model.SaveInfo(focus.UserId.Hex(), focus.TargetId.Hex(), extras)
		if err == nil {
				return true
		}
		return false
}

// 获取用户好友列表
func (this *userBehaviorServiceImpl) ListsUserFriends(userId string, limit models.ListsParams, extras ...beego.M) ([]bson.ObjectId, *models.Meta) {
		if len(extras) == 0 {
				extras = append(extras, beego.M{})
		}
		var (
				err   error
				query = bson.M{
						"userId":     userId,
						"status":     models.StatusOk,
						"targetType": models.TargetTypeFriend,
				}
				results []bson.ObjectId
				meta    = models.NewMeta()
				model   = this.getUserRelation()
				items   = make([]*models.UserRelation, 2)
		)
		items = items[:0]
		Query := model.NewQuery(query)
		err = Query.Limit(limit.Count()).Skip(limit.Skip()).All(&items)
		if err == nil {
				for _, it := range items {
						results = append(results, bson.ObjectIdHex(it.TargetUserId))
				}
				meta.Count = limit.Count()
				meta.Size = len(results)
				meta.Page = limit.Page()
				meta.Total, _ = model.NewQuery(query).Count()
				meta.Boot()
				return results, meta
		}
		return nil, meta
}

// 添加好友
func (this *userBehaviorServiceImpl) AddFriend(userId string, targetUserId string, extras ...beego.M) error {
		var defaults = beego.M{"targetType": models.TargetTypeFriend, "status": models.StatusOk}
		if len(extras) == 0 {
				extras = append(extras, defaults)
		} else {
				extras[0] = models.Merger(extras[0], defaults)
		}
		return this.getUserRelation().SaveInfo(userId, targetUserId, extras...)
}

// 取消好友
func (this *userBehaviorServiceImpl) CancelFriend(userId string, targetUserId string, extras ...beego.M) error {
		var defaults = beego.M{"targetType": models.TargetTypeFriend, "status": models.StatusCancel}
		if len(extras) == 0 {
				extras = append(extras, defaults)
		} else {
				extras[0] = models.Merger(extras[0], defaults)
		}
		return this.getUserRelation().SaveInfo(userId, targetUserId, extras...)
}

func (this *userBehaviorServiceImpl) getUserFocusModel() *models.UserFocusModel {
		return models.UserFocusModelOf()
}

func (this *userBehaviorServiceImpl) getUserRelation() *models.UserRelationModel {
		return models.UserRelationModelOf()
}

// 用户被关注列表(粉丝列表)
func (this *userBehaviorServiceImpl) GetFans(userId string, limit models.ListsParams) ([]bson.ObjectId, *models.Meta) {
		var (
				err   error
				query = bson.M{
						"focusUserId": this.id(userId),
						"status":      models.StatusOk,
				}
				results []bson.ObjectId
				meta    = models.NewMeta()
				model   = this.getUserFocusModel()
				items   = make([]*models.UserFollow, 2)
		)
		items = items[:0]
		Query := model.NewQuery(query)
		err = Query.Limit(limit.Count()).Skip(limit.Skip()).All(&items)
		if err == nil {
				for _, it := range items {
						results = append(results, it.UserId)
				}
				meta.Count = limit.Count()
				meta.Size = len(results)
				meta.Page = limit.Page()
				meta.Total, _ = model.NewQuery(query).Count()
				meta.Boot()
				return results, meta
		}
		return nil, meta
}

// 用户关注列表
func (this *userBehaviorServiceImpl) GetFollows(userId string, limit models.ListsParams) ([]bson.ObjectId, *models.Meta) {
		var (
				err   error
				query = bson.M{
						"userId": this.id(userId),
						"status": models.StatusOk,
				}
				results []bson.ObjectId
				meta    = models.NewMeta()
				model   = this.getUserFocusModel()
				items   = make([]*models.UserFollow, 2)
		)
		items = items[:0]
		Query := model.NewQuery(query)
		err = Query.Limit(limit.Count()).Skip(limit.Skip()).All(&items)
		if err == nil {
				for _, it := range items {
						results = append(results, it.FocusUserId)
				}
				meta.Count = limit.Count()
				meta.Size = len(results)
				meta.Page = limit.Page()
				meta.Total, _ = model.NewQuery(query).Count()
				meta.Boot()
				return results, meta
		}
		return nil, meta
}

// 获取用户关注用户最新作品列表
func (this *userBehaviorServiceImpl) GetFollowPostsLists(query bson.M, limit models.ListsParams) ([]*models.TravelNotes, *models.Meta) {
		var (
				err   error
				ids   []string
				meta  = models.NewMeta()
				model = this.getUserFocusModel()
				ty    = this.getStr(beego.M(query), "$type")
				items = make([]*models.UserFollow, 2)
		)
		items = items[:0]
		delete(query, "$type")
		if this.isSet(beego.M(query), "userId") {
				query["userId"] = this.getObjectId(beego.M(query), "userId")
		}
		query["status"] = models.StatusOk
		Query := model.NewQuery(query)
		err = Query.Sort("-createdAt").All(&items)
		// 获取用户喜欢列表
		if err == nil {
				for _, it := range items {
						ids = append(ids, it.FocusUserId.Hex())
				}
				query["type"] = ty
				return PostServiceOf().ListsQuery(this.parserQuery(query, ids), limit, "-createdAt")
		}
		return nil, meta
}

// 是否关注
func (this *userBehaviorServiceImpl) IsFollowed(userId, follower string) bool {
		if userId == follower || userId == "" || follower == "" {
				return false
		}
		return models.UserFocusModelOf().Exists(bson.M{"userId": bson.ObjectIdHex(userId), "focusUserId": bson.ObjectIdHex(follower), "status": models.StatusOk})
}

// 是朋友
func (this *userBehaviorServiceImpl) IsFriend(userId, userId2 string) bool {
		if userId == userId2 || userId == "" || userId2 == "" {
				return false
		}
		var n, _ = models.UserFocusModelOf().NewQuery(bson.M{
				"$and": []bson.M{
						{"userId": bson.ObjectIdHex(userId), "focusUserId": bson.ObjectIdHex(userId2), "status": models.StatusOk},
						{"userId": bson.ObjectIdHex(userId2), "focusUserId": bson.ObjectIdHex(userId), "status": models.StatusOk},
				},
		}).Count()
		if n == 2 {
				return true
		}
		return false
}

func (this *userBehaviorServiceImpl) parserQuery(query bson.M, ids []string) bson.M {
		// 7 天内更新的作品
		var conditions = bson.M{
				"userId":  bson.M{"$in": ids},
				"privacy": models.PublicPrivacy,
				"status":  models.StatusAuditPass,
		}
		var t, ok = this.getTime(beego.M(query), "timestamp")
		if !ok {
				t = time.Time{}
		}
		if t.IsZero() {
				t = time.Now()
		}
		var ty = this.getStr(beego.M(query), "type")
		// t.Local()
		if ty == "recently" {
				day, _ := time.ParseDuration("3day")
				conditions["createdAt"] = bson.M{
						"$lte": t.Local(), "$gt": t.Add(-day).Local(),
				}
		}
		if ty == "weekly" {
				day, _ := time.ParseDuration("7day")
				conditions["createdAt"] = bson.M{
						"$lte": t.Local(), "$gt": t.Add(-day).Local(),
				}
		}
		if ty == "month" {
				day, _ := time.ParseDuration("60day")
				conditions["createdAt"] = bson.M{
						"$lte": t.Local(), "$gt": t.Add(-day).Local(),
				}
		}
		fmt.Printf("%v\n", string(libs.JsonEncode(conditions)))
		return conditions
}
