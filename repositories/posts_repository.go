package repositories

import (
		"github.com/astaxie/beego"
		"github.com/globalsign/mgo/bson"
		"github.com/weblfe/travel-app/common"
		"github.com/weblfe/travel-app/libs"
		"github.com/weblfe/travel-app/middlewares"
		"github.com/weblfe/travel-app/models"
		"github.com/weblfe/travel-app/services"
		"strconv"
		"strings"
		"time"
)

type PostsRepository interface {
		Create() common.ResponseJson
		Update() common.ResponseJson
		Audit() common.ResponseJson
		GetLikes(ids ...string) common.ResponseJson
		GetRanking() common.ResponseJson
		GetFollows() common.ResponseJson
		Lists(...string) common.ResponseJson
		GetById(...string) common.ResponseJson
		RemoveId(...string) common.ResponseJson
		AutoVideosCover() common.ResponseJson
		GetAll() common.ResponseJson
		ListsByPostType(typ string) common.ResponseJson
}

type postRepositoryImpl struct {
		service services.PostService
		ctx     common.BaseRequestContext
}

const (
		PostImagesInfoKey = "imagesInfo"
		PostVideoInfoKey  = "videosInfo"
)

func NewPostsRepository(ctx common.BaseRequestContext) PostsRepository {
		var repository = new(postRepositoryImpl)
		repository.init()
		repository.ctx = ctx
		repository.GetDto()
		return repository
}

func (this *postRepositoryImpl) GetDto() *DtoRepository {
		return GetDtoRepository()
}

func (this *postRepositoryImpl) init() {
		this.service = services.PostServiceOf()
}

func (this *postRepositoryImpl) Create() common.ResponseJson {
		var (
				err    error
				data   = new(models.TravelNotes)
				userId = this.GetUserId()
		)
		if err = this.ctx.JsonDecode(data); err != nil {
				err = this.ctx.GetParent().ParseForm(data)
				if err != nil {
						return common.NewErrorResp(common.NewErrors(common.EmptyParamCode, err), "参数不足")
				}
		}
		if data.IsEmpty() {
				return common.NewErrorResp(common.NewErrors(common.EmptyParamCode, "post create failed"), "参数不足")
		}
		data.UserId = userId
		// 自动过滤敏感词
		data.Content = models.GetDfaInstance().ChangeSensitiveWords(data.Content)
		// 仅内容时 自动通过
		if data.Type == models.ContentType {
				data.Status = models.StatusAuditPass
		}
		err = this.service.Create(data.Defaults())
		if err != nil {
				return common.NewErrorResp(common.NewErrors(common.ServiceFailed, err), "发布失败")
		}
		return common.NewSuccessResp(beego.M{"timestamp": time.Now().Unix(), "id": data.Id.Hex()}, "发布成功")
}

func (this *postRepositoryImpl) GetUserId() string {
		var (
				ctx = this.ctx.GetParent()
				id  = ctx.GetSession(middlewares.AuthUserId)
		)
		if id == nil || id == "" {
				return ""
		}
		if v, ok := id.(string); ok {
				return v
		}
		return ""
}

func (this *postRepositoryImpl) Update() common.ResponseJson {
		var (
				id = this.ctx.GetParent().GetString(":id")
		)
		if id == "" {
				return common.NewInvalidParametersResp("游记id缺失")
		}
		// @todo 更新
		return common.NewInDevResp(this.ctx.GetActionId())
}

func (this *postRepositoryImpl) Lists(typ ...string) common.ResponseJson {
		var (
				meta     *models.Meta
				ctx      = this.ctx.GetParent()
				items    []*models.TravelNotes
				ty       = ctx.GetString("type")
				page, _  = ctx.GetInt("page", 1)
				count, _ = ctx.GetInt("count", 20)
				limit    = models.NewListParam(page, count)
		)
		if ty == "" && len(typ) != 0 {
				ty = typ[0]
		}

		var extras = beego.M{"privacy": models.PublicPrivacy, "status": models.StatusAuditPass}
		switch ty {
		case "my":
				id := ctx.GetSession(middlewares.AuthUserId)
				if id == nil {
						return common.NewFailedResp(common.RecordNotFound, common.RecordNotFoundError)
				}
				delete(extras, "privacy")
				extras["status"] = beego.M{"$ne": models.StatusAuditNotPass}
				// 查看自己的作品
				items, meta = this.service.Lists(id.(string), limit, extras)
		case "address":
				// 通过地址罗列作品
				items, meta = this.service.ListByAddress(ctx.GetString("address"), limit, extras)
		case "tags":
				// 通过标签罗列作品
				tags := ctx.GetString("tags")
				if tags != "" {
						items, meta = this.service.ListByTags(strings.SplitN(tags, ",", -1), limit, extras)
				} else {
						tags := ctx.GetStrings("tags")
						items, meta = this.service.ListByTags(tags, limit, extras)
				}
		case "user":
				if len(typ) <= 1 {
						break
				}
				// 查询他人的作品列表
				userId := typ[1]
				items, meta = this.service.Lists(userId, limit, extras)
		case "search":
				items, meta = this.service.Search(this.parseSearchQuery(this.ctx.GetString("query")), limit)
		}
		if items != nil && len(items) > 0 && meta != nil {
				var arr []beego.M
				for _, item := range items {
						arr = append(arr, item.M(this.getPostTransform()))
				}
				return common.NewSuccessResp(beego.M{"items": arr, "meta": meta}, "罗列成功")
		}
		return common.NewFailedResp(common.RecordNotFound, common.RecordNotFoundError)
}

// 查询解析
func (this *postRepositoryImpl) parseSearchQuery(query string) beego.M {
		if query == "" {
				return beego.M{}
		}
		var (
				queryMapper = beego.M{}
				supportKeys = []string{"address", "content"}
				// "startAt", "endAt", "nickname",
		)
		if !strings.Contains(query, ":") {
				var or = make([]bson.M, 1)
				or = or[:0]
				for _, key := range supportKeys {
						or = append(or, bson.M{key: bson.RegEx{Pattern: query, Options: "i"}})
				}
				queryMapper["$or"] = or
				return queryMapper
		}
		if strings.Contains(query, "&") {
				items := strings.SplitN(query, "&", -1)
				for _, value := range items {
						values := strings.SplitN(value, ":", 2)
						if len(values) < 2 {
								continue
						}
						if !libs.InArray(values[0], supportKeys) {
								continue
						}
						queryMapper[values[0]] = values[1]
				}
				return queryMapper
		}
		if strings.Contains(query, ":") {
				values := strings.SplitN(query, ":", 2)
				if len(values) < 2 {
						return beego.M{}
				}
				if !libs.InArray(values[0], supportKeys) {
						return beego.M{}
				}
				return beego.M{values[0]: values[1]}
		}
		return beego.M{}
}

func (this *postRepositoryImpl) GetById(id ...string) common.ResponseJson {
		if len(id) == 0 {
				id = append(id, this.ctx.GetParent().GetString(":id"))
		}
		var data = this.service.GetById(id[0])
		if data == nil || data.DeletedAt != 0 {
				return common.NewFailedResp(common.RecordNotFound, common.RecordNotFoundError)
		}
		return common.NewSuccessResp(data.M(this.getPostTransform()), "获取成功")
}

func (this *postRepositoryImpl) RemoveId(id ...string) common.ResponseJson {
		if len(id) == 0 {
				id = append(id, this.ctx.GetParent().GetString(":id"))
		}
		var (
				userId = this.GetUserId()
				data   = this.service.GetById(id[0])
		)
		if data == nil {
				return common.NewFailedResp(common.RecordNotFound, common.RecordNotFoundError)
		}
		if data.UserId != userId {
				return common.NewFailedResp(common.PermissionCode, common.PermissionError)
		}
		data.DeletedAt = time.Now().Unix()
		err := data.Save()
		if err == nil {
				return common.NewSuccessResp(beego.M{"timestamp": time.Now().Unix()}, "删除成功")
		}
		return common.NewFailedResp(common.RecordNotFound, "删除失败")
}

func (this *postRepositoryImpl) Audit() common.ResponseJson {
		var (
				userId  = getUserId(this.ctx)
				ids     = this.ctx.GetStrings("ids")
				comment = this.ctx.GetString("comment")
				typ     = this.ctx.GetString("type")
		)
		if len(ids) == 0 {
				var data = struct {
						Ids     []string `json:"ids"`
						Type    string   `json:"type"`
						Comment string   `json:"comment"`
				}{}
				_ = this.ctx.JsonDecode(&data)
				if len(data.Ids) > 0 && this.service.Audit(data.Type, data.Ids...) {
						this.service.AddAuditLog(userId, typ, comment, ids)
						return common.NewSuccessResp(bson.M{"timestamp": time.Now().Unix()}, "审核完成")
				}
				return common.NewFailedResp(common.ServiceFailed, "审核失败")
		}
		if this.service.Audit(typ, ids...) {
				this.service.AddAuditLog(userId, typ, comment, ids)
				return common.NewSuccessResp(bson.M{"timestamp": time.Now().Unix()}, "审核完成")
		}
		return common.NewFailedResp(common.ServiceFailed, "审核失败")
}

// 获取文章内容转换器
func (this *postRepositoryImpl) getPostTransform() func(m beego.M) beego.M {
		var dto = this.GetDto()
		return func(m beego.M) beego.M {
				m = getMediaInfoTransform()(m)
				m = TransBigNumberToText(m, "commentNum", "thumbsUpNum")
				var (
						id, _      = m["id"]
						userId, ok = m["userId"]
				)
				m["userInfo"] = nil
				// 获取用户是否已点赞
				if id != nil && id != "" {
						var (
								currentUserId = this.GetUserId()
						)
						value := dto.IsThumbsUp(id.(string), currentUserId)
						m["isUp"] = value
						m["isFollowed"] = services.UserBehaviorServiceOf().IsFollowed(currentUserId, userId.(string))
						// m["isThumbsUp"] = dto.IsThumbsUp(id.(string),currentUserId)
				}
				if !ok {
						return m
				}
				if id, ok := userId.(string); ok {
						var (
								key   = dto.Key(id)
								value = dto.Get(key)
						)
						if value == nil {
								value = dto.GetUserById(id)
						}
						dto.Cache(key, value)
						m["userInfo"] = value
				}
				return m
		}
}

// 是否作者
func (this *postRepositoryImpl) IsAuthor(postId, userId string) bool {
		var post = services.PostServiceOf().GetById(postId)
		if post == nil {
				return false
		}
		if post.UserId == userId {
				return true
		}
		return false
}

// 获取喜欢列表
func (this *postRepositoryImpl) GetLikes(ids ...string) common.ResponseJson {
		if len(ids) == 0 {
				ids = append(ids, getUserId(this.ctx))
		}
		var (
				meta     *models.Meta
				ctx      = this.ctx.GetParent()
				items    []*models.TravelNotes
				page, _  = ctx.GetInt("page", 1)
				count, _ = ctx.GetInt("count", 20)
				limit    = models.NewListParam(page, count)
		)
		var query = bson.M{"userId": ids[0]}
		items, meta = services.ThumbsUpServiceOf().GetUserLikeLists(query, limit)
		if items != nil && len(items) > 0 && meta != nil {
				var arr []beego.M
				for _, item := range items {
						arr = append(arr, item.M(this.getPostTransform()))
				}
				return common.NewSuccessResp(beego.M{"items": arr, "meta": meta}, "罗列成功")
		}
		return common.NewFailedResp(common.RecordNotFound, common.RecordNotFoundError)
}

// 获取排行榜列表
func (this *postRepositoryImpl) GetRanking() common.ResponseJson {
		var (
				meta     *models.Meta
				ctx      = this.ctx.GetParent()
				items    []*models.TravelNotes
				ty, _    = ctx.GetInt("type", 0)
				page, _  = ctx.GetInt("page", 1)
				count, _ = ctx.GetInt("count", 20)
				limit    = models.NewListParam(page, count)
		)
		var extras = bson.M{"privacy": models.PublicPrivacy, "status": models.StatusAuditPass}
		if ty == models.ImageType || ty == models.VideoType || ty == models.ContentType {
				extras["type"] = ty
		}
		items, meta = this.service.GetRankingLists(extras, limit)
		// 分页
		if items != nil && len(items) > 0 && meta != nil {
				var arr []beego.M
				for _, item := range items {
						arr = append(arr, item.M(this.getPostTransform()))
				}
				return common.NewSuccessResp(beego.M{"items": arr, "meta": meta}, "罗列成功")
		}
		return common.NewFailedResp(common.RecordNotFound, common.RecordNotFoundError)
}

// 获取关注列表
func (this *postRepositoryImpl) GetFollows() common.ResponseJson {
		var (
				meta     *models.Meta
				ctx      = this.ctx.GetParent()
				items    []*models.TravelNotes
				userId   = getUserId(this.ctx)
				page, _  = ctx.GetInt("page", 1)
				count, _ = ctx.GetInt("count", 20)
				limit    = models.NewListParam(page, count)
		)
		if userId == "" {
				return common.NewFailedResp(common.RecordNotFound, common.RecordNotFoundError)
		}
		var query = bson.M{"userId": userId}
		query["$type"] = ctx.GetString("type", "all")
		items, meta = services.UserBehaviorServiceOf().GetFollowPostsLists(query, limit)
		if items != nil && len(items) > 0 && meta != nil {
				var arr []beego.M
				for _, item := range items {
						arr = append(arr, item.M(this.getPostTransform()))
				}
				return common.NewSuccessResp(beego.M{"items": arr, "meta": meta}, "罗列成功")
		}
		return common.NewFailedResp(common.RecordNotFound, common.RecordNotFoundError)
}

// 通过作品类型获取列表
func (this *postRepositoryImpl) ListsByPostType(typ string) common.ResponseJson {
		var (
				ty       int
				meta     *models.Meta
				ctx      = this.ctx.GetParent()
				items    []*models.TravelNotes
				page, _  = ctx.GetInt("page", 1)
				count, _ = ctx.GetInt("count", 20)
				limit    = models.NewListParam(page, count)
		)
		switch typ {
		case models.ImageTypeCode:
				ty = models.ImageType
		case models.VideoTypeCode:
				ty = models.VideoType
		case models.ContentTypeCode:
				ty = models.ContentType
		}
		var extras = bson.M{"privacy": models.PublicPrivacy, "status": models.StatusAuditPass}
		if ty != 0 {
				extras["type"] = ty
		}
		items, meta = this.service.GetRecommendLists(extras, limit)
		if items != nil && len(items) > 0 && meta != nil {
				var arr []beego.M
				for _, item := range items {
						arr = append(arr, item.M(this.getPostTransform()))
				}
				return common.NewSuccessResp(beego.M{"items": arr, "meta": meta}, "罗列成功")
		}
		return common.NewFailedResp(common.RecordNotFound, common.RecordNotFoundError)
}

func (this *postRepositoryImpl) AutoVideosCover() common.ResponseJson {
		var ids = this.ctx.GetStrings("ids")
		defer this.service.AutoVideoCoverImageTask(ids)
		return common.NewSuccessResp(bson.M{"count": len(ids), "timestamp": time.Now().Unix()}, "自动截图任务已经下放成功")
}

// 获取所有
func (this *postRepositoryImpl) GetAll() common.ResponseJson {
		var (
				user     = getUser(this.ctx)
				ctx      = this.ctx.GetParent()
				page, _  = ctx.GetInt("page", 1)
				count, _ = ctx.GetInt("count", 20)
				limit    = models.NewListParam(page, count)
				typ      = this.ctx.GetInt("type", 0)
				status   = this.ctx.GetString("status")
		)
		if user == nil {
				return common.NewFailedResp(common.RecordNotFound, common.RecordNotFoundError)
		}
		if !user.IsRootRole() {
				return common.NewFailedResp(common.PermissionCode, common.PermissionError)
		}
		//
		var query = beego.M{
				"status": beego.M{"$nin": []int{models.StatusAuditNotPass, models.StatusAuditPass}},
		}
		// 自定义状态查询
		if status != "" {
				var statusArr []int
				arr := strings.SplitN(status, ",", -1)
				for _, it := range arr {
						n, err := strconv.Atoi(it)
						if err == nil {
								continue
						}
						statusArr = append(statusArr, n)
				}
				if len(statusArr) > 0 {
						query["status"] = beego.M{
								"$in": statusArr,
						}
				}
		}
		// 作品类型
		if typ != 0 {
				query["type"] = typ
		}
		// 查看自己的作品
		items, meta := this.service.All(query, limit, "-createdAt", "-updatedAt")
		if items != nil && len(items) > 0 && meta != nil {
				var arr []beego.M
				for _, item := range items {
						arr = append(arr, item.M(this.getPostTransform()))
				}
				return common.NewSuccessResp(beego.M{"items": arr, "meta": meta}, "罗列成功")
		}
		return common.NewFailedResp(common.RecordNotFound, common.RecordNotFoundError)
}
