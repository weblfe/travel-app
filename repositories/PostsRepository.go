package repositories

import (
		"github.com/astaxie/beego"
		"github.com/weblfe/travel-app/common"
		"github.com/weblfe/travel-app/middlewares"
		"github.com/weblfe/travel-app/models"
		"github.com/weblfe/travel-app/services"
		"time"
)

type PostsRepository interface {
		Create() common.ResponseJson
		Update() common.ResponseJson
		Lists(...string) common.ResponseJson
		GetById(...string) common.ResponseJson
		RemoveId(...string) common.ResponseJson
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
		return repository
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
				ctx      = this.ctx.GetParent()
				meta     *models.Meta
				items    []*models.TravelNotes
				page, _  = ctx.GetInt("page", 1)
				count, _ = ctx.GetInt("count", 20)
				limit    = models.NewListParam(page, count)
				ty       = ctx.GetString("type")
		)
		if ty == "" && len(typ) != 0 {
				ty = typ[0]
		}
		switch ty {
		case "my":
				id := ctx.GetSession(middlewares.AuthUserId)
				if id == nil {
						return common.NewFailedResp(common.RecordNotFound, common.RecordNotFoundError)
				}
				items, meta = this.service.Lists(id.(string), limit)
		case "address":
				items, meta = this.service.ListByAddress(ctx.GetString(":address"), limit)
		case "tags":
				address := ctx.GetString("tags")
				items, meta = this.service.ListByAddress(address, limit)
		case "user":
				if len(typ) <= 1 {
						break
				}
				userId := typ[1]
				items, meta = this.service.Lists(userId, limit)
		case "search":
				items, meta = this.service.Search(beego.M{}, limit)
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
		var data = this.service.GetById(id[0])
		if data == nil {
				return common.NewFailedResp(common.RecordNotFound, common.RecordNotFoundError)
		}
		data.DeletedAt = time.Now().Unix()
		err := data.Save()
		if err == nil {
				common.NewSuccessResp(beego.M{"timestamp": time.Now().Unix()}, "删除成功")
		}
		return common.NewFailedResp(common.RecordNotFound, "删除失败")
}

// 获取文章内容转换器
func (this *postRepositoryImpl) getPostTransform() func(m beego.M) beego.M {
		return func(m beego.M) beego.M {
				m = getMediaInfoTransform()(m)
				var (
						id, _      = m["id"]
						userId, ok = m["userId"]
						dto        = GetDtoRepository()
				)
				m["userInfo"] = nil
				// 获取用户是否已点赞
				if id != nil && id != "" {
						m["isUp"] = dto.IsThumbsUp(id.(string), this.GetUserId())
				}
				if !ok {
						return m
				}
				if id, ok := userId.(string); ok {
						m["userInfo"] = dto.GetUserById(id)
				}
				return m
		}
}
