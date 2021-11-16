package repositories

import (
		"github.com/astaxie/beego"
		"github.com/globalsign/mgo/bson"
		"github.com/weblfe/travel-app/common"
		"github.com/weblfe/travel-app/models"
		"github.com/weblfe/travel-app/services"
		"time"
)

// BehaviorRepository 用户相关行为
type BehaviorRepository interface {
		FocusOn(ids ...string) common.ResponseJson
		FocusOff(ids ...string) common.ResponseJson
		GetUserFollows(ids ...string) common.ResponseJson
		GetUserFans(ids ...string) common.ResponseJson
}

type userBehaviorRepositoryImpl struct {
		ctx             common.BaseRequestContext
		userService     services.UserService
		behaviorService services.UserBehaviorService
		dto             *DtoRepository
}

func NewBehaviorRepository(ctx common.BaseRequestContext) BehaviorRepository {
		var repository = newUserBehaviorRepository()
		repository.ctx = ctx
		return repository
}

func newUserBehaviorRepository() *userBehaviorRepositoryImpl {
		var repository = new(userBehaviorRepositoryImpl)
		repository.init()
		return repository
}

func (this *userBehaviorRepositoryImpl) init() {
		this.dto = GetDtoRepository()
		this.userService = services.UserServiceOf()
		this.behaviorService = services.UserBehaviorServiceOf()
}

// FocusOff 取消关注
func (this *userBehaviorRepositoryImpl) FocusOff(ids ...string) common.ResponseJson {
		if len(ids) == 0 {
				ids = append(ids, "")
		}
		var (
				err             error
				userId          = getUserId(this.ctx)
				targetUserId, _ = this.ctx.GetParam(":userId", this.ctx.GetString("userId", ids[0]))
				extras, _       = this.ctx.GetParam("extras", beego.M{})
				query           = beego.M{"_id": bson.ObjectIdHex(targetUserId.(string)), "deletedAt": 0}
		)
		if targetUserId == "" {
				return common.NewUnLoginResp("targetUser empty!")
		}
		if targetUserId == userId {
				return common.NewUnLoginResp("targetUser error!")
		}
		if userId == "" {
				return common.NewUnLoginResp("please login!")
		}
		if !this.userService.Exists(query) {
				return common.NewUnLoginResp("targetUser not exists!")
		}
		err = this.behaviorService.UnFollow(userId, targetUserId.(string), extras.(beego.M))
		if err == nil {
				return common.NewSuccessResp(bson.M{"timestamp": time.Now().Unix()}, "取关成功")
		}
		return common.NewFailedResp(common.ServiceFailed, "取关失败")
}

// FocusOn 关注
func (this *userBehaviorRepositoryImpl) FocusOn(ids ...string) common.ResponseJson {
		if len(ids) == 0 {
				ids = append(ids, "")
		}
		var (
				err             error
				userId          = getUserId(this.ctx)
				extras, _       = this.ctx.GetParam("extras", beego.M{})
				targetUserId, _ = this.ctx.GetParam(":userId", this.ctx.GetString("userId", ids[0]))
				query           = beego.M{"_id": bson.ObjectIdHex(targetUserId.(string)), "deletedAt": 0}
		)
		if targetUserId == "" {
				return common.NewFailedResp(common.ServiceFailed, "follow targetUser required!")
		}
		if targetUserId == userId {
				return common.NewFailedResp(common.ServiceFailed, "follow yourself error!")
		}
		if userId == "" {
				return common.NewUnLoginResp("please login!")
		}
		if !this.userService.Exists(query) {
				return common.NewFailedResp(common.ServiceFailed, "targetUser not exists!")
		}
		var targetId = targetUserId.(string)
		if userId == targetId {
				return common.NewFailedResp(common.ServiceFailed, "follower error!")
		}
		err = this.behaviorService.Follow(userId, targetId, extras.(beego.M))
		if err == nil {
				return common.NewSuccessResp(bson.M{"timestamp": time.Now().Unix()}, "关注成功")
		}
		return common.NewFailedResp(common.ServiceFailed, "关注失败")
}

// GetUserFans 用户粉丝列表
func (this *userBehaviorRepositoryImpl) GetUserFans(ids ...string) common.ResponseJson {
		if len(ids) == 0 {
				ids = append(ids, this.ctx.GetString("id", ""))
		}
		var (
				userId        = ids[0]
				users         = make([]*BaseUser, 2)
				currentUserId = getUserId(this.ctx)
				page, count   = this.ctx.GetInt("page", 1), this.ctx.GetInt("count", 20)
				limit         = models.NewListParam(page, count)
		)
		if userId == "" && currentUserId != "" {
				userId = currentUserId
		}
		if userId == "0" {
				return common.NewUnLoginResp("error params")
		}
		var items, meta = this.behaviorService.GetFans(userId, limit)
		if items == nil || meta == nil {
				return common.NewFailedResp(common.RecordNotFound, "空")
		}
		users = users[:0]
		var dto = this.getDto()
		for _, user := range items {
				it := dto.GetUserById(user.Hex())
				users = append(users, it)
		}
		if len(users) == 0 {
				return common.NewFailedResp(common.RecordNotFound, "空")
		}
		return common.NewSuccessResp(bson.M{"items": users, "meta": meta}, "获取成功")
}

// GetUserFollows 用户关注列表
func (this *userBehaviorRepositoryImpl) GetUserFollows(ids ...string) common.ResponseJson {
		if len(ids) == 0 {
				ids = append(ids, "")
		}
		var (
				userId        = ids[0]
				users         = make([]beego.M, 2)
				currentUserId = getUserId(this.ctx)
				page, count   = this.ctx.GetInt("page", 1), this.ctx.GetInt("count", 20)
				limit         = models.NewListParam(page, count)
		)

		if userId == "" && currentUserId != "" {
				userId = currentUserId
		}
		if userId == "0" || userId == "" {
				return common.NewUnLoginResp("error params")
		}
		var items, meta = this.behaviorService.GetFollows(userId, limit)
		if items == nil || meta == nil {
				return common.NewFailedResp(common.RecordNotFound, "空")
		}
		users = users[:0]
		var dto = this.getDto()
		for _, user := range items {
				it := dto.GetUserById(user.Hex())
				users = append(users, it.M(dto.appendFollowStatus(userId)))
		}
		if len(users) == 0 {
				return common.NewFailedResp(common.RecordNotFound, "空")
		}
		return common.NewSuccessResp(bson.M{"items": users, "meta": meta}, "获取成功")
}

func (this *userBehaviorRepositoryImpl) getDto() *DtoRepository {
		return GetDtoRepository()
}
